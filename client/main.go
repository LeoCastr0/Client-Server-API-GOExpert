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
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	GetDollar(ctx)
}

func GetDollar(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Println(error.Error(err))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(error.Error(err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(error.Error(err))
	}
	println(string(body))
}
