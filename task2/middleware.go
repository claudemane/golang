package main

import (
	"net/http"
	"net/http/httptest"
)

func IdempotencyMiddleware(store *MemoryStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Missing Idempotency-Key", http.StatusBadRequest)
			return
		}

		if cached, exists := store.Get(key); exists {
			if cached.Completed {
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
				return
			}
			http.Error(w, "Request in progress", http.StatusConflict)
			return
		}

		if !store.StartProcessing(key) {
			http.Error(w, "Duplicate request", http.StatusConflict)
			return
		}

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		store.Finish(key, rec.Code, rec.Body.Bytes())

		w.WriteHeader(rec.Code)
		w.Write(rec.Body.Bytes())
	})
}