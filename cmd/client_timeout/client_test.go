package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestFlagParser(t *testing.T) {
	tests := []struct {
		c             *http.Client
		url           string
		duration      time.Duration
		shouldTimeout bool
	}{
		{newSimpleTimeout(), "https://httpstat.us/200?sleep=6000", 5 * time.Second, true},
		{newCompleteTimeout(), "https://httpstat.us/200?sleep=6000", 5 * time.Second, true},
	}

	ctx := context.Background()
	for _, tc := range tests {
		ctx, cancel := context.WithTimeout(ctx, tc.duration)
		defer cancel()

		if err := execute(ctx, tc.c, tc.url); err != nil {
			if !tc.shouldTimeout {
				t.Errorf("the call was expected to timeout")
			}
		}
	}
}

func newSimpleTimeout() *http.Client {
	return &http.Client{
		Timeout: 3 * time.Second,
	}
}

func newCompleteTimeout() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

func execute(ctx context.Context, client *http.Client, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
