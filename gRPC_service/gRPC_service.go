package gRPC_service

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/nik-zaitsev/fibonacci_service/gRPC_service/pb"

	"github.com/nik-zaitsev/fibonacci_service/fibonacci_calculator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedFibonacciServer
}

func (s *Service) GetFibonacciSlice(ctx context.Context, values *pb.BorderValues) (*pb.FibonacciSlice, error) {
	log.Printf("new request received: From = %d, To = %d", values.From, values.To)
	if resSlice, err := fibonacci_calculator.Fibonacci(values.From, values.To, ctx.Done()); err != nil {
		if err == fibonacci_calculator.OperationRejected {
			log.Printf("operation rejected from context")
			return nil, status.Error(codes.Aborted, ctx.Err().Error())
		} else {
			log.Printf("bad arguments, skipping..")
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	} else {
		log.Printf("sending back answer: %v", resSlice)
		return &pb.FibonacciSlice{FibonacciNums: resSlice}, nil
	}
}

func RunGRPCServer(rpcServer *grpc.Server, wg *sync.WaitGroup, port uint64) {
	defer wg.Done()
	lsn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterFibonacciServer(rpcServer, new(Service))
	log.Printf("starting gRPC server on %s", lsn.Addr().String())
	if err := rpcServer.Serve(lsn); err != nil {
		log.Printf("error while running gRPC server, %s", err.Error())
	}
}
