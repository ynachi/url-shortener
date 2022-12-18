package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Home define a home handler function which respond to users hiting the home page
func (srv *Server) Home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't, use
	// the http.NotFound() function to send a 404 response to the client.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Welcome to Yao's Cloud Native URL shortening App"))
}

// CreateURL creates a shortened URL from a long URL
func (srv *Server) CreateURL(w http.ResponseWriter, r *http.Request) {
	// Only POST method is accepted to create resources
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	longURL := r.URL.Query().Get("url")
	u, err := url.ParseRequestURI(longURL)
	if err != nil {
		http.Error(w, "Bad URL format", http.StatusBadRequest)
		return
	}
	// Do not accept to parse localhost
	if u.Host == "127.0.0.1" || u.Host == "localhost" {
		http.Error(w, "Localhost shortening is not allowed", http.StatusBadRequest)
		return
	}
	// remove trailing / as we want goo.com and goo.com/ to encode the same ID
	longURL = strings.TrimSuffix(longURL, "/")
	var requestData = requestURL{URL: longURL}
	encodedID, err := requestData.encodeLongURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Logger.Info("persisting encoded URL", "long_url", longURL, "short_url_id", encodedID)
	err = persistURL(srv.ctx, longURL, encodedID, srv.firestoreClient)
	if err != nil {
		http.Error(w, "Failed to save encoded", http.StatusInternalServerError)
		return
	}
	const msg = `"{message": Short url %s created and saved.}`
	fmt.Fprintf(w, msg, encodedID)
}

// DeleteURL deletes a saved URL
func (srv *Server) DeleteURL(w http.ResponseWriter, r *http.Request) {

}

// UpdateURL updates the expiration date of a given shortened URL
func (srv *Server) UpdateURL(w http.ResponseWriter, r *http.Request) {

}

// GetURL ViewURL displays information about a shortened URL, like the long URL it points to
func (srv *Server) GetURL(w http.ResponseWriter, r *http.Request) {
	// Only GET method is accepted to create resources
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	shortID := r.URL.Query().Get("shortid")
	longURL, err := decodeURL(srv.ctx, shortID, srv.redisClient, srv.firestoreClient)
	if err != nil {
		switch {
		case errors.Is(err, ErrCacheSave):
			// log the error but respond to the request as the answer is still valid
			Logger.Error("failed to save cold item to cache", err, "redis_host", srv.redisAddr)
			const msg = `"{message": Long url for short url %s is %s.}`
			fmt.Fprintf(w, msg, shortID, longURL)
		case errors.Is(err, ErrStorageMiss):
			Logger.Error("item not found", err, "short_id", shortID)
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		default:
			Logger.Error("internal error", err, "short_id", shortID)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
	} else {
		const msg = `"{message": Long url for short url %s is %s.}`
		fmt.Fprintf(w, msg, shortID, longURL)
	}
}

// GetURLs ViewURLs displays all shortened URLs matching some criteria
func (srv *Server) GetURLs(w http.ResponseWriter, r *http.Request) {

}

// Redirect redirect redirect a long URL to a shortened URL
func (srv *Server) Redirect(w http.ResponseWriter, r *http.Request) {

}
