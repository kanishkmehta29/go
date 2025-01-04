package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server interface{
	Address() string
	IsAlive() bool
	Serve(w http.ResponseWriter,req *http.Request)

}

type SimpleServer struct{
	addr string
	proxy *httputil.ReverseProxy
}

func (s *SimpleServer) Address() string{
	return s.addr
}

func (s *SimpleServer) IsAlive() bool{
	return true
}

func (s *SimpleServer) Serve(w http.ResponseWriter,req *http.Request){
	s.proxy.ServeHTTP(w,req)
}

func newSimpleServer(addr string) *SimpleServer{
	SimpleServerUrl,err := url.Parse(addr)
	if err != nil{
	 log.Fatalln("Error parsing url")
	}
 
	return &SimpleServer{
	 addr: addr,
	 proxy: httputil.NewSingleHostReverseProxy(SimpleServerUrl),
	}
 }

type LoadBalancer struct{
	port string
	roundRobinCount int
	servers []Server
}

func (lb *LoadBalancer) getNextAvailableServer() Server{
    server := lb.servers[lb.roundRobinCount%(len(lb.servers))]
	for !server.IsAlive(){
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%(len(lb.servers))]
	}
	lb.roundRobinCount++
	return server
}

func (lb *LoadBalancer) serveProxy(w http.ResponseWriter,req *http.Request){
    targetserver := lb.getNextAvailableServer()
	targetserver.Serve(w,req)
}

func NewLoadBalancer(port string,servers []Server) *LoadBalancer{
	return &LoadBalancer{
		port: port,
		roundRobinCount: 0,
        servers: servers,
	}
}

func main(){
	servers := []Server{
		newSimpleServer("https://www.wikipedia.com"),
		newSimpleServer("https://apple.com"),
		newSimpleServer("https://www.screener.in"),
	}
	lb := NewLoadBalancer("8080",servers)
	http.HandleFunc("/",func(rw http.ResponseWriter, req *http.Request) {
		lb.serveProxy(rw, req)
	})
	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080",nil)

}