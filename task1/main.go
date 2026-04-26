package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++

		if attempts <= 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Println("Server: returning 503")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
		fmt.Println("Server: returning 200")
	}))

	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := ExecutePayment(ctx, server.URL)
	if err != nil {
		fmt.Println("Final error:", err)
	}
}