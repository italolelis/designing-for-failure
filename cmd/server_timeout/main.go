package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	h := timeoutHandler{}

	srv := &http.Server{
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
		Handler:           h,
	}

	log.Println(srv.ListenAndServe())
}

type timeoutHandler struct{}

func (h timeoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	timer := time.AfterFunc(5*time.Second, func() {
		r.Body.Close()
	})

	bodyBytes := make([]byte, 0)
	for {
		//We reset the timer, for the variable time
		timer.Reset(1 * time.Second)

		_, err := io.CopyN(bytes.NewBuffer(bodyBytes), r.Body, 256)
		if err == io.EOF {
			// This is not an error in the common sense
			// io.EOF tells us, that we did read the complete body
			break
		} else if err != nil {
			//You should do error handling here
			break
		}
	}
}
