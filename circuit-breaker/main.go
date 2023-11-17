package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex
	return func(ctx context.Context) (string, error) {
		m.RLock() // Establish a "read lock"
		d := consecutiveFailures - int(failureThreshold)
		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}
		m.RUnlock()                   // Release read lock
		response, err := circuit(ctx) // Issue request proper
		m.Lock()                      // Lock around shared resources
		defer m.Unlock()
		lastAttempt = time.Now() // Record time of attempt
		if err != nil {          // Circuit returned an error,
			consecutiveFailures++ // so we count the failure
			fmt.Printf("Request failed : %d times!", consecutiveFailures)
			return response, err // and return
		}
		fmt.Println("Request Succeeded!")
		consecutiveFailures = 0 // Reset failures counter
		return response, nil
	}
}

// ExternalAPIService simulates an external API with potential failures.
type ExternalAPIService struct {
	counter int
}

func (e *ExternalAPIService) Call() (string, error) {
	// Simulate external API call
	e.counter++
	if e.counter%10 != 0 {
		return "", errors.New("External API error")
	}
	return "External API response", nil
}
func main() {
	externalAPI := &ExternalAPIService{}
	circuit := func(ctx context.Context) (string, error) {
		return externalAPI.Call()
	}

	breaker := Breaker(circuit, 3) // Break after 3 consecutive failures

	for i := 0; i < 10; i++ {
		response, err := breaker(context.Background())
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Success: %v\n", response)
		}

		time.Sleep(time.Second)
	}
}
