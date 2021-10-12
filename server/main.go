package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/nik-zaitsev/fibonacci_service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var InvalidParametersValues = errors.New("invalid parameters values")

func Fibonacci(from uint64, to uint64) ([]uint64, error) {
	if from < 1 || to < 1 || to < from {
		return nil, InvalidParametersValues
	}

	res := make([]uint64, to-from+1)
	if from == 1 {
		res[0] = 1
	}

	var n2, n1 uint64 = 0, 1
	for i := uint64(1); i < to; i++ {
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

func (s *Service) GetFibonacciSlice(_ context.Context, values *pb.BorderValues) (*pb.FibonacciSlice, error) {
	log.Printf("new request received: From = %d, To = %d", values.From, values.To)
	if resSlice, err := Fibonacci(values.From, values.To); err != nil {
		log.Printf("bad arguments, skipping..")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else {
		log.Printf("sending back answer: %v", resSlice)
		return &pb.FibonacciSlice{FibonacciNums: resSlice}, nil
	}
}

type HttpHandler struct {
}

func (h *HttpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	argFrom, errArgFrom := strconv.ParseUint(args.Get("from"), 10, 64)
	if errArgFrom != nil {
		log.Printf("fail parsing FROM parameter to uint64: %v", errArgFrom)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	argTo, errArgTo := strconv.ParseUint(args.Get("to"), 10, 64)
	if errArgTo != nil {
		log.Printf("fail parsing TO parameter to uint64: %v", errArgTo)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("new request received: From = %d, To = %d", argFrom, argTo)
	if resSlice, err := Fibonacci(argFrom, argTo); err != nil {
		log.Printf("bad arguments, %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	} else {
		resp := make(map[string][]uint64, 1)
		resp["fibonacciSlice"] = resSlice
		if jsonResp, err := json.Marshal(resp); err != nil {
			if _, err := rw.Write([]byte(err.Error())); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			rw.Header().Set("Content-Type", "application/json; charset=utf-8")
			rw.WriteHeader(http.StatusOK)
			if _, err := rw.Write(jsonResp); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Printf("sending back answer: %v", resp)
		}
	}
}

func RunGRPCServer(rpcServer *grpc.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	lsn, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterFibonacciServer(rpcServer, new(Service))
	log.Printf("starting gRPC server on %s", lsn.Addr().String())
	if err := rpcServer.Serve(lsn); err != nil {
		log.Printf("error while running gRPC server, %s", err.Error())
	}
}

func RunHTTPServer(httpServer *http.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("starting HTTP server on %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Printf("error while running HTTP server, %s", err.Error())
	}
}

func main() {
	wg := &sync.WaitGroup{}

	rpcServer := grpc.NewServer()
	wg.Add(1)
	go RunGRPCServer(rpcServer, wg)

	httpHandler := &HttpHandler{}
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler,
	}
	wg.Add(1)
	go RunHTTPServer(httpServer, wg)

	wg.Wait()
}
