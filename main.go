package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// ----------------------------
// SERVER INTERFACE
// ----------------------------
type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

// ----------------------------
// SIMPLE SERVER (BACKEND)
// ----------------------------
type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	handleErr(err)

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func (s *simpleServer) Address() string {
	return s.addr
}

func (s *simpleServer) IsAlive() bool {
	// For simplicity return true
	return true
}

func (s *simpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}

// ----------------------------
// LOAD BALANCER
// ----------------------------
type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

// ----------------------------
// ROUND ROBIN SERVER PICKER
// ----------------------------
func (lb *LoadBalancer) getNextAvailable() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	lb.roundRobinCount++
	return server
}

// ----------------------------
// PROXY REQUEST HANDLER
// ----------------------------
func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	server := lb.getNextAvailable()
	fmt.Printf("Forwarding request to: %s\n", server.Address())
	server.Serve(rw, r)
}

// ----------------------------
// MAIN FUNCTION
// ----------------------------
func main() {

	// BACKEND SERVERS (YOU CAN REPLACE WITH YOUR APIs)
	servers := []Server{
		newSimpleServer("https://www.facebook.com"),
		newSimpleServer("https://www.google.com"),
		newSimpleServer("https://www.youtube.com"),
	}

	lb := NewLoadBalancer("8000", servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serveProxy(rw, req)
	}

	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Load Balancer running at http://localhost:%s\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
