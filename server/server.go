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
	Logger          = slog.New(slog.NewJSONHandler(os.Stdout))
	Ctx             = context.Background()
	GCPProjectID    = os.Getenv("URL_SHORTNER_PROJECT_ID")
	redisAddr       = fmt.Sprintf("%s:%s", os.Getenv("REDISHOST"), os.Getenv("REDISPORT"))
	redisClient     *redis.Client
	firestoreClient *firestore.Client
)

// Server is a struct representing an instance of the url shortening web application
type Server struct {
	IPAddr string
	Port   int
	Mux    *http.ServeMux
	Logger *slog.Logger
	Ctx    context.Context
}

// init logger and external services clients
func init() {
	// multiple assigment is not recognize in init().
	// for instance, firestoreClient, err := NewFirestoreClient(ctx, GCPProjectID) will claim firestoreClient is a
	// new variable despite being declared outside init()
	var err error
	slog.SetDefault(Logger)
	Ctx, cancelCtx := context.WithCancel(Ctx)
	defer cancelCtx()
	firestoreClient, err = NewFirestoreClient(Ctx, GCPProjectID)
	if err != nil {
		cancelCtx()
	}
	redisClient, err = NewRedisClient(redisAddr)
	if err != nil {
		cancelCtx()
	}
}

// MakeServer creates a new instance of a url shortening server. The port should be a valid port and not a port reserved for clients.
// For reference, ports reserved to clients are 49152 - 65535 and valid port ranges are 0 - 65535.
func MakeServer(ipaddr string, port int) (Server, error) {
	mux := http.NewServeMux()
	ip := net.ParseIP(ipaddr)
	if ip == nil {
		return Server{ipaddr, port, nil}, fmt.Errorf("ip address %s is invalid", ipaddr)
	}
	if port < 0 || port > portMax {
		return Server{ipaddr, port, nil}, fmt.Errorf("port number %d is invalid", port)
	}
	if port >= userPortMin {
		return Server{ipaddr, port, nil}, fmt.Errorf("user reserved port %d is not allowed", port)
	}
	return Server{ipaddr, port, mux}, nil
}

// Start starts the server
func (srv Server) Start(ctx context.Context) error {
	listenAddr := srv.IPAddr + ":" + strconv.Itoa(srv.Port)
	// fail if the init() method returned an error
	if err := ctx.Err(); err != nil {
		Logger.Error("server initialization failed", err, "redis_server", redisAddr, "firestore_project_id", GCPProjectID)
		return err
	}
	return http.ListenAndServe(listenAddr, srv.Mux)
}
