package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/ynachi/url-shortner/server"
)

func main() {
	ctx := context.Background()
	gcpProjectID := os.Getenv("URLS_GCP_PROJECT_ID")
	port := os.Getenv("URLS_SERVER_PORT")
	serverPort, err := strconv.Atoi(port)
	if err != nil {
		server.Logger.Error("invalid port number", err, "port_number", port)
		os.Exit(2)
	}
	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDISHOST"), os.Getenv("REDISPORT"))
	// This is a conternerized workload so we set it to listen on all interfaces for now
	srv, err := server.NewServer(ctx, "0.0.0.0", serverPort, gcpProjectID, redisAddr)
	if err != nil {
		server.Logger.Error("server creation failed", err, "name", "server")
		os.Exit(1)
	}
	server.Logger.Info("server instantiated", "port", srv.Port)

	// registering the handlers
	srv.Mux.HandleFunc("/", Home)
	srv.Mux.HandleFunc("/url/create", CreateURL)
	srv.Mux.HandleFunc("/url/delete", DeleteURL)
	srv.Mux.HandleFunc("/url/update", UpdateURL)
	srv.Mux.HandleFunc("/url/get", GetURL)
	srv.Mux.HandleFunc("/urls/get", GetURLs)
	srv.Mux.HandleFunc("/url/redirect", Redirect)

	server.Logger.Info("starting server", "port", srv.Port)
	err = srv.Start()
	if err != nil {
		server.Logger.Error("server startup failed", err, "port", srv.Port)
		os.Exit(1)
	}
}
