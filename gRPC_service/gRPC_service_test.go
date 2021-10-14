package gRPC_service

import (
	"context"
	"testing"

	"github.com/nik-zaitsev/fibonacci_service/fibonacci_calculator"
	"github.com/nik-zaitsev/fibonacci_service/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRunGRPCServer(t *testing.T) {
	rpcServer := grpc.NewServer()
	go RunGRPCServer(rpcServer, nil)
	defer rpcServer.Stop()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	client := pb.NewFibonacciClient(conn)

	t.Run("correct parameters", func(t *testing.T) {
		res, err := client.GetFibonacciSlice(context.Background(), &pb.BorderValues{From: 1, To: 2})
		require.Equal(t, []uint64{1, 1}, res.FibonacciNums)
		require.Nil(t, err)
	})

	t.Run("bad parameters", func(t *testing.T) {
		res, err := client.GetFibonacciSlice(context.Background(), &pb.BorderValues{From: 0, To: 1})
		require.Nil(t, res)
		require.Equal(t, status.Error(codes.InvalidArgument, fibonacci_calculator.InvalidParametersValues.Error()), err)
	})
}
