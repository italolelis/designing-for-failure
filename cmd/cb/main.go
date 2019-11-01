package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/italolelis/designing-for-failure/internal/pkg/profile"
	"github.com/sony/gobreaker"
)

func main() {
	// Arguments to a defer statement is immediately evaluated and stored.
	// The deferred function receives the pre-evaluated values when its invoked.
	defer profile.Duration(time.Now(), "main")

	rand.Seed(time.Now().UnixNano())

	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name: "HTTP GET Example",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	})

	statusCodes := []int{http.StatusOK, http.StatusBadRequest, http.StatusInternalServerError}
	statusCodesBuckets := make(map[int]int)
	totalReq := 50

	fmt.Println("Making requests...")
	for i := 0; i < totalReq; i++ {
		url := fmt.Sprintf("https://httpstat.us/%d", statusCodes[rand.Intn(len(statusCodes))])

		resp, err := cb.Execute(func() (interface{}, error) {
			resp, err := http.Get(url)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode >= http.StatusBadRequest {
				return resp, errors.New("request failed")
			}

			return resp, nil
		})
		if err != nil && resp == nil {
			fmt.Println(err)
			statusCodesBuckets[0]++
			continue
		}
		res := resp.(*http.Response)

		statusCodesBuckets[res.StatusCode]++
		fmt.Printf("URL %s \n", url)
	}

	for c, i := range statusCodesBuckets {
		percentage := (float64(i) * float64(100)) / float64(totalReq)
		fmt.Printf("Code %d: %d (%g%%) \n", c, i, percentage)
	}
}
