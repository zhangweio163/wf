package main

import (
	"Reginster/http"
	"Reginster/server"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

var (
	port     = flag.Int("port", 8500, "port")
	HttpPort = flag.Int("http", 8300, "http")
)

func RunServer(work *sync.WaitGroup) {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	register := server.NewRegServer()
	register.InitServer()
	server.RegisterRegisterServer(s, register)
	log.Printf("RegServer started on port %d. Can be used Ctrl+C to stop server", *port)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	work.Done()
}
func RunHttpServer(work *sync.WaitGroup) {
	flag.Parse()
	http.HttpServer(*HttpPort)
	work.Done()
}

func main() {
	var work sync.WaitGroup
	work.Add(2)
	go RunServer(&work)
	go RunHttpServer(&work)
	work.Wait()
}
