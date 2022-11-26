package main

import (
	"os"

	"github.com/ynachi/url-shortner/server"
	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout))
	slog.SetDefault(logger)
	srv, err := server.MakeServer("0.0.0.0", 8080)
	if err != nil {
		logger.Error("server creation failed",
			err,
			"name", "server")
		return
	}
	logger.Info("server instantiated",
		"port", srv.Port)
	srv.Mux.HandleFunc("/", server.Home)
	logger.Info("starting server",
		"port", srv.Port)
	err = srv.Start()
	if err != nil {
		logger.Error("server startup failed",
			err,
			"port", srv.Port)
		return
	}
}
