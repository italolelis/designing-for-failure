package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gojektech/heimdall/hystrix"
	"github.com/italolelis/designing-for-failure/internal/pkg/profile"
)

func main() {
	// Arguments to a defer statement is immediately evaluated and stored.
	// The deferred function receives the pre-evaluated values when its invoked.
	defer profile.Duration(time.Now(), "main")

	// Create a new hystrix-wrapped HTTP client with the fallbackFunc as fall-back function
	client := hystrix.NewClient(
		hystrix.WithHTTPTimeout(10*time.Second),
		hystrix.WithCommandName("MyCommand"),
		hystrix.WithHystrixTimeout(10*time.Second),
		hystrix.WithMaxConcurrentRequests(100),
		hystrix.WithErrorPercentThreshold(40),
		hystrix.WithSleepWindow(10),
		hystrix.WithRequestVolumeThreshold(10),
	)

	statusCodes := []int{http.StatusOK, http.StatusBadRequest, http.StatusInternalServerError}
	statusCodesBuckets := make(map[int]int)
	totalReq := 50

	for i := 0; i <= totalReq; i++ {
		url := fmt.Sprintf("https://httpstat.us/%d", statusCodes[rand.Intn(len(statusCodes))])

		res, err := client.Get(url, nil)
		if err != nil {
			statusCodesBuckets[0] = statusCodesBuckets[0] + 1
			continue
		}

		statusCodesBuckets[res.StatusCode] = statusCodesBuckets[res.StatusCode] + 1
	}

	for c, i := range statusCodesBuckets {
		percentage := ((float64(totalReq) * float64(i)) / float64(100))
		fmt.Printf("Code %d: %d (%g%%) \n", c, i, percentage)
	}
}
