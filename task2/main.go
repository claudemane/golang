package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("💰 Processing payment...")
	time.Sleep(2 * time.Second)

	w.WriteHeader(200)
	w.Write([]byte(`{"status":"paid","amount":1000}`))
}

func main() {
	store := NewMemoryStore()

	handler := IdempotencyMiddleware(store, http.HandlerFunc(paymentHandler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go server.ListenAndServe()

	time.Sleep(time.Second)

	var wg sync.WaitGroup
	key := "abc-123"

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
			req.Header.Set("Idempotency-Key", key)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("Request", i, "status:", resp.StatusCode)
		}(i)
	}

	wg.Wait()
}