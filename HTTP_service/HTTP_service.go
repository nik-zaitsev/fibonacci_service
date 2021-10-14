package HTTP_service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/nik-zaitsev/fibonacci_service/fibonacci_calculator"
)

type HttpHandler struct {
}

type JsonAnswer struct {
	FibonacciSlice []uint64 `json:"fibonacciSlice,omitempty"`
	Error          string   `json:"error,omitempty"`
}

func FillHttpResponseBody(rw http.ResponseWriter, resp JsonAnswer) {
	if jsonResp, err := json.Marshal(resp); err != nil {
		if _, err := rw.Write([]byte(err.Error())); err != nil {
			log.Printf("could not write json marshalling error to buffer")
		}
	} else {
		if _, err := rw.Write(jsonResp); err != nil {
			log.Printf("could not write json message to buffer")
		}
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	log.Println("sending back answer...")
}

func (h *HttpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	args := r.URL.Query()
	argFrom, errArgFrom := strconv.ParseUint(args.Get("from"), 10, 64)
	if errArgFrom != nil {
		log.Printf("fail parsing FROM parameter to uint64: %v", errArgFrom)
		FillHttpResponseBody(rw, JsonAnswer{Error: errArgFrom.Error()})
		return
	}
	argTo, errArgTo := strconv.ParseUint(args.Get("to"), 10, 64)
	if errArgTo != nil {
		log.Printf("fail parsing TO parameter to uint64: %v", errArgTo)
		FillHttpResponseBody(rw, JsonAnswer{Error: errArgTo.Error()})
		return
	}
	log.Printf("new request received: From = %d, To = %d", argFrom, argTo)
	if resSlice, err := fibonacci_calculator.Fibonacci(argFrom, argTo, make(<-chan struct{})); err != nil {
		log.Printf("bad arguments, %v", err)
		FillHttpResponseBody(rw, JsonAnswer{Error: err.Error()})
		return
	} else {
		FillHttpResponseBody(rw, JsonAnswer{FibonacciSlice: resSlice})
		return
	}
}

func RunHTTPServer(httpServer *http.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("starting HTTP server on %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Printf("error while running HTTP server, %s", err.Error())
	}
}
