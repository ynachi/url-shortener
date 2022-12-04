package server

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
)

// Server is a struct representing an instance of the url shortening web application
type Server struct {
	IPAddr string
	Port   int
	Mux    *http.ServeMux
}

const portMax = 65535
const userPortMin = 49152

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
func (srv Server) Start() error {
	listenAddr := srv.IPAddr + ":" + strconv.Itoa(srv.Port)
	return http.ListenAndServe(listenAddr, srv.Mux)
}
