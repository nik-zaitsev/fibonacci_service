package main

import (
	"net/http"
	"sync"

	"github.com/nik-zaitsev/fibonacci_service/HTTP_service"
	"github.com/nik-zaitsev/fibonacci_service/gRPC_service"
	"google.golang.org/grpc"
)

func main() {
	wg := &sync.WaitGroup{}

	rpcServer := grpc.NewServer()
	wg.Add(1)
	go gRPC_service.RunGRPCServer(rpcServer, wg)

	httpHandler := &HTTP_service.HttpHandler{}
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler,
	}
	wg.Add(1)
	go HTTP_service.RunHTTPServer(httpServer, wg)

	wg.Wait()
}
