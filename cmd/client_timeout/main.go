package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
		},
	}

	url := "https://httpstat.us/200?sleep=6000"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("failed to create request: %s", err)
	}

	fmt.Printf("Making request to %s \n", url)
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Code: %d \n", res.StatusCode)
}
