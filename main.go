package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/nik-zaitsev/fibonacci_service/HTTP_service"
	"github.com/nik-zaitsev/fibonacci_service/gRPC_service"
	"google.golang.org/grpc"
)

func ParsePort(service string) uint64 {
	var portString string
	fmt.Printf("enter %s port...\n", service)
	for {
		if _, err := fmt.Scan(&portString); err != nil {
			fmt.Println(err.Error())
		} else {
			if portUint, errParse := strconv.ParseUint(portString, 10, 64); errParse != nil {
				fmt.Println(errParse.Error())
			} else {
				return portUint
			}
		}
	}
}

func main() {
	wg := &sync.WaitGroup{}

	rpcServer := grpc.NewServer()
	wg.Add(1)
	grpcPort := ParsePort("gRPC")
	go gRPC_service.RunGRPCServer(rpcServer, wg, grpcPort)

	httpPort := ParsePort("HTTP")
	httpHandler := &HTTP_service.HttpHandler{}
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: httpHandler,
	}
	wg.Add(1)
	go HTTP_service.RunHTTPServer(httpServer, wg)

	wg.Wait()
}
