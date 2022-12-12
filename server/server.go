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

// CacheDuration Should be LRU like in real world. For now, we use time based expiration
const CacheDuration = 4 * time.Hour

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

// MakeServer creates a new instance of a url shortening server. The port should be a valid port and not a port reserved for clients.
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
	firestoreClient, err1 := NewFirestoreClient(ctx, gcpProjectID)
	redisClient, err2 := NewRedisClient(redisAddr)
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
func (srv Server) Start() error {
	listenAddr := srv.IPAddr + ":" + srv.Port
	if srv.firestoreClient != nil {
		defer srv.firestoreClient.Close()
	}
	if srv.redisClient != nil {
		defer srv.redisClient.Close()
	}
	return http.ListenAndServe(listenAddr, srv.Mux)
}
