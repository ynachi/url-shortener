package main

import (
	"os"

	"github.com/ynachi/url-shortner/server"
	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout))
	slog.SetDefault(logger)
	srv, err := server.MakeServer("0.0.0.0", 8080)
	if err != nil {
		logger.Error("server creation failed", err, "name", "server")
		return
	}
	logger.Info("server instantiated", "port", srv.Port)

	// registrering the handlers
	srv.Mux.HandleFunc("/", server.Home)
	srv.Mux.HandleFunc("/url/create", server.CreateURL)
	srv.Mux.HandleFunc("/url/delete", server.DeleteURL)
	srv.Mux.HandleFunc("/url/update", server.UpdateURL)
	srv.Mux.HandleFunc("/url/view", server.ViewURL)
	srv.Mux.HandleFunc("/urls/view", server.ViewURLs)
	srv.Mux.HandleFunc("/url/redirect", server.Redirect)

	logger.Info("starting server", "port", srv.Port)
	err = srv.Start()
	if err != nil {
		logger.Error("server startup failed", err, "port", srv.Port)
		return
	}
}
