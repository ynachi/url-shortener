package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-redis/redis/v8"
	"golang.org/x/exp/slog"
)

const portMax = 65535
const userPortMin = 49152

// CacheDuration is the expiration time set on each cached object. We set it at 0 as we will be using LRU or LFU at the server side.
// Could be a parameter passed to the server as environment variable.
const CacheDuration = 0 * time.Hour

// Set logger and sentinel errors
var (
	Logger                   = slog.New(slog.NewJSONHandler(os.Stdout))
	ErrRedisClientCreate     = errors.New("unable to create redis client")
	ErrFirestoreClientCreate = errors.New("unable to create firestore client")
)

// Server is a struct representing an instance of the url shortening web application
type Server struct {
	IPAddr          string
	Port            string
	Mux             *http.ServeMux
	ctx             context.Context
	projectID       string
	redisAddr       string
	redisClient     *redis.Client
	firestoreClient *firestore.Client
}

// init logger
func init() {
	slog.SetDefault(Logger)
}

// NewServer MakeServer creates a new instance of an url shortening server. The port should be a valid port and not a port reserved for clients.
// For reference, ports reserved to clients are 49152 - 65535 and valid port ranges are 0 - 65535. ipaddr and ports should be valid.
func NewServer(ctx context.Context, ipaddr string, port string, gcpProjectID string, redisAddr string) (*Server, error) {
	serverPort, err := strconv.Atoi(port)
	if err != nil {
		Logger.Error("invalid port number", err, "port_number", port)
		return nil, err
	}
	ip := net.ParseIP(ipaddr)
	switch {
	case ip == nil:
		return nil, fmt.Errorf("ip address %s is invalid", ipaddr)
	case serverPort < 0 || serverPort > portMax:
		return nil, fmt.Errorf("port number %s is invalid", port)
	case serverPort >= userPortMin:
		return nil, fmt.Errorf("user reserved port %s is not allowed", port)
	}
	mux := http.NewServeMux()
	firestoreClient, err1 := newFirestoreClient(ctx, gcpProjectID)
	redisClient, err2 := newRedisClient(redisAddr)
	srv := &Server{
		IPAddr:          ipaddr,
		Port:            port,
		Mux:             mux,
		ctx:             ctx,
		projectID:       gcpProjectID,
		redisAddr:       redisAddr,
		redisClient:     redisClient,
		firestoreClient: firestoreClient,
	}
	// We want to return a non nil server instance even if redis and firestore clients fail
	if err1 != nil {
		return srv, err1
	}
	if err2 != nil {
		return srv, err2
	}
	return srv, nil
}

// Start starts the server
func (srv *Server) Start() error {
	listenAddr := srv.IPAddr + ":" + srv.Port
	if srv.firestoreClient != nil {
		defer func(firestoreClient *firestore.Client) {
			err := firestoreClient.Close()
			if err != nil {
				Logger.Warn("failed to close firestore client", err)
			}
		}(srv.firestoreClient)
	}
	if srv.redisClient != nil {
		defer func(redisClient *redis.Client) {
			err := redisClient.Close()
			if err != nil {
				Logger.Warn("failed to close redis client", err)
			}
		}(srv.redisClient)
	}
	return http.ListenAndServe(listenAddr, srv.Mux)
}
