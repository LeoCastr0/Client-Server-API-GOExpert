package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	ContexHandler(ctx)
}

func ContexHandler(ctx context.Context) {
	GetDollar(ctx)
	select {
	case <-time.After(300 * time.Millisecond):
		log.Println("Request successfully processed")
	case <-ctx.Done():
		log.Println("Request Cancelled")
	}
}

func GetDollar(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	println(string(body))
}
