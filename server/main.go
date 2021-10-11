package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/nik-zaitsev/fibonacci_service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

var InvalidParametersValues = errors.New("invalid parameters values")

func Fibonacci(from uint64, to uint64) ([]uint64, error) {
	if from < 1 || to < 1 || to <= from {
		return nil, InvalidParametersValues
	}

	res := make([]uint64, to-from+1)
	if from == 1 {
		res[0] = 1
	}

	var n2, n1 uint64 = 0, 1
	for i := uint64(1); i < to; i++ {
		fmt.Println(i, n1, n2)
		n2, n1 = n1, n1+n2
		if i >= from-1 {
			res[i-from+1] = n1
		}
	}

	return res, nil
}

type Service struct {
	pb.UnimplementedFibonacciServer
}

func (s *Service) GetFibonacciSlice(ctx context.Context, values *pb.BorderValues) (*pb.FibonacciSlice, error) {
	log.Printf("new request received: From = %d, To = %d", values.From, values.To)
	if resSlice, err := Fibonacci(values.From, values.To); err != nil {
		log.Printf("bad arguments, skipping..")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else {
		return &pb.FibonacciSlice{FibonacciNums: resSlice}, nil
	}
}

func main() {
	lsn, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	pb.RegisterFibonacciServer(server, new(Service))

	log.Printf("starting server on %s", lsn.Addr().String())
	if err := server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}
