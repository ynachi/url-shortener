package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ynachi/url-shortner/server"
)

func main() {
	ctx := context.Background()
	gcpProjectID := os.Getenv("URLS_GCP_PROJECT_ID")
	serverPort := os.Getenv("URLS_SERVER_PORT")
	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDISHOST"), os.Getenv("REDISPORT"))
	srv, err := server.NewServer(ctx, "0.0.0.0", serverPort, gcpProjectID, redisAddr)
	if err != nil {
		switch {
		case errors.Is(err, server.ErrRedisClientCreate):
			server.Logger.Error("redis client error, caching disabled", err, "redis_host", redisAddr)
		case errors.Is(err, server.ErrFirestoreClientCreate):
			server.Logger.Error("firestore client error", err, "project_id", gcpProjectID)
		default:
			server.Logger.Error("server creation failed", err)
			os.Exit(2)
		}
	}
	server.Logger.Info("server instantiated", "port", srv.Port)

	// registering the handlers
	srv.Mux.HandleFunc("/", srv.Redirect)
	srv.Mux.HandleFunc("/urls/", srv.URLs)
	
	server.Logger.Info("starting server", "port", srv.Port)
	err = srv.Start()
	if err != nil {
		server.Logger.Error("server startup failed", err, "port", srv.Port)
		os.Exit(1)
	}
}
