package server

import (
	"fmt"
	"net/http"
	"net/url"
)

// Home define a home handler function which writes a byte slice containing
func Home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't, use
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the handler
	// would keep executing and also write the "Hello from SnippetBox" message.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Welcome to Yao's Cloud Native URL shortening App"))
}

// CreateURL creates a shortened URL from a long URL
func CreateURL(w http.ResponseWriter, r *http.Request) {
	// Only POST method is accepted to create ressources
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
	var requestData = requestURL{longURL}
	encodedID, err := requestData.encodeLongURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "URL that is being procecessed: %s", encodedID)
}

// DeleteURL deletes a saved URL
func DeleteURL(w http.ResponseWriter, r *http.Request) {

}

// UpdateURL updates the exoiration date of a given shortened URL
func UpdateURL(w http.ResponseWriter, r *http.Request) {

}

// ViewURL displays information about a shortened URL, like the long URL it points to
func ViewURL(w http.ResponseWriter, r *http.Request) {

}

// ViewURLs displays all shortened URLs matching some criterias
func ViewURLs(w http.ResponseWriter, r *http.Request) {

}

// redirect redirect a long URL to a shortened URL
func Redirect(w http.ResponseWriter, r *http.Request) {

}
