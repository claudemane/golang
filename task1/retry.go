package main

import (
	"math"
	"math/rand"
	"net/http"
	"time"
)

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true // network error
	}

	if resp == nil {
		return true
	}

	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	case 401, 404:
		return false
	default:
		return false
	}
}

func CalculateBackoff(attempt int) time.Duration {
	baseDelay := 500 * time.Millisecond
	maxDelay := 5 * time.Second

	backoff := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
	if backoff > maxDelay {
		backoff = maxDelay
	}

	jitter := time.Duration(rand.Int63n(int64(backoff)))
	return jitter
}