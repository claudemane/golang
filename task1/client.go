package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func ExecutePayment(ctx context.Context, url string) error {
	client := &http.Client{}

	var resp *http.Response
	var err error

	maxRetries := 5

	for attempt := 0; attempt < maxRetries; attempt++ {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		resp, err = client.Do(req)

		// 1. Сначала проверяем на УСПЕХ (перенесите этот блок ВЫШЕ)
		if err == nil && resp.StatusCode == 200 {
			fmt.Println("✅ Success on attempt:", attempt+1)
			return nil
		}

		// 2. И только потом проверяем, нужно ли делать ретрай
		if !IsRetryable(resp, err) {
			return err
		}

		delay := CalculateBackoff(attempt)
		fmt.Printf("❌ Attempt %d failed. Retrying in %v\n", attempt+1, delay)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("failed after retries")
}
