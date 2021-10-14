package HTTP_service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

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
