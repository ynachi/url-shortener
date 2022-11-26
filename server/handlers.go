package server

import "net/http"

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Yao's App"))
}
