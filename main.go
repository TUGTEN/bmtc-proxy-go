package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	// c := cache.New(5*time.Second, 30*time.Second)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		url := "https://bmtcmobileapi.karnataka.gov.in" + r.URL.Path

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		start := time.Now()

		fmt.Println(
			url,
		)
		defer func() {
			fmt.Println(
				r.URL.Path, time.Since(start).Milliseconds(), "ms",
			)
			defer cancel()
		}()

		req, err := http.NewRequestWithContext(ctx, r.Method, url, r.Body)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for name, headers := range r.Header {
			for _, h := range headers {
				req.Header.Add(name, h)
			}
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		defer resp.Body.Close()

		for name, headers := range resp.Header {
			for _, h := range headers {
				w.Header().Add(name, h)
			}
		}

		w.WriteHeader(resp.StatusCode)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		w.Write(bodyBytes)

	})

	corsMux := cors.AllowAll().Handler(mux)

	port := "5000"

	println("Server Listening on port: ", port)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, corsMux))

}
