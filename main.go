package main

import (
	"github.com/ynachi/url-shortner/server"
)

func main() {
	srv, err := server.MakeServer("0.0.0.0", 8080)
	if err != nil {
		server.Logger.Error("server creation failed", err, "name", "server")
		return
	}
	server.Logger.Info("server instantiated", "port", srv.Port)

	// registering the handlers
	srv.Mux.HandleFunc("/", server.Home)
	srv.Mux.HandleFunc("/url/create", server.CreateURL)
	srv.Mux.HandleFunc("/url/delete", server.DeleteURL)
	srv.Mux.HandleFunc("/url/update", server.UpdateURL)
	srv.Mux.HandleFunc("/url/get", server.GetURL)
	srv.Mux.HandleFunc("/urls/get", server.GetURLs)
	srv.Mux.HandleFunc("/url/redirect", server.Redirect)

	server.Logger.Info("starting server", "port", srv.Port)
	err = srv.Start(server.Ctx)
	if err != nil {
		server.Logger.Error("server startup failed", err, "port", srv.Port)
		return
	}
}
