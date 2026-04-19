package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"base":"USD","target":"EUR","rate":0.92}`))
	}))
	defer server.Close()

	svc := NewExchangeService(server.URL)
	rate, err := svc.GetRate("USD", "EUR")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rate != 0.92 {
		t.Errorf("expected rate 0.92, got %f", rate)
	}
}

func TestGetRate_APIBusinessError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid currency pair"}`))
	}))
	defer server.Close()

	svc := NewExchangeService(server.URL)
	_, err := svc.GetRate("USD", "INVALID")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "api error: invalid currency pair" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetRate_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Internal Server Error`))
	}))
	defer server.Close()

	svc := NewExchangeService(server.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}

func TestGetRate_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	svc := NewExchangeService(server.URL)
	svc.Client.Timeout = 100 * time.Millisecond

	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestGetRate_ServerPanic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer server.Close()

	svc := NewExchangeService(server.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetRate_EmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(``))
	}))
	defer server.Close()

	svc := NewExchangeService(server.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected error for empty body, got nil")
	}
}
