package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Stage 1 - provided, do not edit
func startServer(ctx context.Context, name string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(rand.Intn(500)) * time.Millisecond):
				out <- fmt.Sprintf("[%s] metric: %d", name, rand.Intn(100))
			}
		}
	}()
	return out
}

// FanIn merges multiple input channels into a single output channel.
// It closes the result channel only after all input channels are drained and closed.
// It also respects context cancellation.
func FanIn(ctx context.Context, channels ...<-chan string) <-chan string {
	result := make(chan string)
	var wg sync.WaitGroup

	// For each input channel, start a goroutine that forwards messages to result
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-c:
					if !ok {
						// Channel closed — goroutine is done
						return
					}
					select {
					case result <- val:
					case <-ctx.Done():
						return
					}
				}
			}
		}(ch)
	}

	// Close result channel only after all goroutines finish
	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch1 := startServer(ctx, "Alpha")
	ch2 := startServer(ctx, "Beta")
	ch3 := startServer(ctx, "Gamma")

	ch4 := FanIn(ctx, ch1, ch2, ch3)

	for val := range ch4 {
		fmt.Println(val)
	}
}
