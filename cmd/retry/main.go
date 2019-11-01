package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
)

func main() {
	initalTimeout := 2 * time.Millisecond         // Inital timeout
	maxTimeout := 9 * time.Millisecond            // Max time out
	exponentFactor := 2.0                         // Multiplier
	maximumJitterInterval := 2 * time.Millisecond // Max jitter interval. It must be more than 1*time.Millisecond

	backoff := heimdall.NewExponentialBackoff(initalTimeout, maxTimeout, exponentFactor, maximumJitterInterval)

	// Create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	// Create a new hystrix-wrapped HTTP client with the fallbackFunc as fall-back function
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(10*time.Second),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(3),
	)

	fmt.Println("Making requests...")
	statusCodes := []int{http.StatusOK, http.StatusBadRequest, http.StatusInternalServerError}
	for i := 0; i <= 50; i++ {
		url := fmt.Sprintf("https://httpstat.us/%d", statusCodes[rand.Intn(len(statusCodes))])
		fmt.Printf("GET %s \n", url)

		if err := get(client, url); err != nil {
			fmt.Printf("failed: %s \n", err)
		}

		fmt.Println("success")
	}
}

func get(client heimdall.Client, url string) error {
	// Use the clients GET method to create and execute the request
	res, err := client.Get(url, nil)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusOK {
		return nil
	}

	return errors.New("request failed")
}
