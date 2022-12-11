package server

import (
	"context"
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

// Set package level variables. These are mostly shared service like contexts, logging and storage partameters
var (
	Logger = slog.New(slog.NewJSONHandler(os.Stdout))
)

// Server is a struct representing an instance of the url shortening web application
type Server struct {
	IPAddr          string
	Port            int
	Mux             *http.ServeMux
	Ctx             context.Context
	GCPProjectID    string
	RedisAddr       string
	RedisClient     *redis.Client
	FirestoreClient *firestore.Client
}

// init logger
func init() {
	slog.SetDefault(Logger)
}

// MakeServer creates a new instance of a url shortening server. The port should be a valid port and not a port reserved for clients.
// For reference, ports reserved to clients are 49152 - 65535 and valid port ranges are 0 - 65535. ipaddr and ports should be valid.
func NewServer(ctx context.Context, ipaddr string, port int, gcpProjectID string, redisAddr string) (*Server, error) {
	ip := net.ParseIP(ipaddr)
	switch {
	case ip == nil:
		return nil, fmt.Errorf("ip address %s is invalid", ipaddr)
	case port < 0 || port > portMax:
		return nil, fmt.Errorf("port number %d is invalid", port)
	case port >= userPortMin:
		return nil, fmt.Errorf("user reserved port %d is not allowed", port)
	}
	mux := http.NewServeMux()
	firestoreClient, err := NewFirestoreClient(ctx, gcpProjectID)
	if err != nil {
		return nil, err
	}
	redisClient, err := NewRedisClient(redisAddr)
	if err != nil {
		return nil, err
	}
	srv := &Server{
		IPAddr:          ipaddr,
		Port:            port,
		Mux:             mux,
		Ctx:             ctx,
		GCPProjectID:    gcpProjectID,
		RedisAddr:       redisAddr,
		RedisClient:     redisClient,
		FirestoreClient: firestoreClient,
	}
	return srv, nil
}

// Start starts the server
func (srv Server) Start() error {
	listenAddr := srv.IPAddr + ":" + strconv.Itoa(srv.Port)
	defer srv.FirestoreClient.Close()
	defer srv.RedisClient.Close()
	return http.ListenAndServe(listenAddr, srv.Mux)
}
