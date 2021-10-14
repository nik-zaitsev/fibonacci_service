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
	if resSlice, err := fibonacci_calculator.Fibonacci(argFrom, argTo); err != nil {
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

func RunHTTPServer(httpServer *http.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("starting HTTP server on %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Printf("error while running HTTP server, %s", err.Error())
	}
}
