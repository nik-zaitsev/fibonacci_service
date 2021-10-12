package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nik-zaitsev/fibonacci_service/pb"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/require"
)

func TestFibonacci(t *testing.T) {
	testCases := []struct {
		inputFrom   uint64
		inputTo     uint64
		expectedRes []uint64
		expectedErr error
	}{
		{inputFrom: 1, inputTo: 1, expectedRes: []uint64{1}, expectedErr: nil},
		{inputFrom: 1, inputTo: 2, expectedRes: []uint64{1, 1}, expectedErr: nil},
		{inputFrom: 1, inputTo: 3, expectedRes: []uint64{1, 1, 2}, expectedErr: nil},
		{inputFrom: 2, inputTo: 3, expectedRes: []uint64{1, 2}, expectedErr: nil},
		{inputFrom: 3, inputTo: 3, expectedRes: []uint64{2}, expectedErr: nil},
		{inputFrom: 2, inputTo: 4, expectedRes: []uint64{1, 2, 3}, expectedErr: nil},
		{inputFrom: 0, inputTo: 2, expectedRes: nil, expectedErr: InvalidParametersValues},
		{inputFrom: 2, inputTo: 0, expectedRes: nil, expectedErr: InvalidParametersValues},
		{inputFrom: 0, inputTo: 0, expectedRes: nil, expectedErr: InvalidParametersValues},
		{inputFrom: 4, inputTo: 3, expectedRes: nil, expectedErr: InvalidParametersValues},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("FROM=%d, TO=%d", tc.inputFrom, tc.inputTo), func(t *testing.T) {
			res, err := Fibonacci(tc.inputFrom, tc.inputTo)
			require.Equal(t, tc.expectedRes, res)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}

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
		require.Equal(t, status.Error(codes.InvalidArgument, InvalidParametersValues.Error()), err)
	})
}

func TestRunHTTPServer(t *testing.T) {
	httpHandler := &HttpHandler{}
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler,
	}
	go RunHTTPServer(httpServer, nil)
	defer func() {
		err := httpServer.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	t.Run("correct parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "localhost:8080?from=1&to=2", nil)
		w := httptest.NewRecorder()
		httpHandler.ServeHTTP(w, req)
		res := w.Result()
		defer func() {
			err := res.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()
		data, errRead := ioutil.ReadAll(res.Body)
		require.Nil(t, errRead)
		j := make(map[string][]uint64, 1)
		errJ := json.Unmarshal(data, &j)
		if errJ != nil {
			t.Fatal(errJ)
		}
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Equal(t, map[string][]uint64{"fibonacciSlice": {1, 1}}, j)
	})

	t.Run("bad parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "localhost:8080?to=2", nil)
		w := httptest.NewRecorder()
		httpHandler.ServeHTTP(w, req)
		res := w.Result()
		defer func() {
			err := res.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

}
