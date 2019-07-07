package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		bucket      = make(chan struct{})
	)
	defer cancel()
	go tokenGen(ctx, bucket, 50)

	start := time.Now()
	for i := 1; i <= 100; i++ {
		<-bucket
		if i%10 == 0 {
			fmt.Println(".")
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println("finish", time.Since(start))
}

func tokenGen(ctx context.Context, bucket chan<- struct{}, ratePerSec uint) {
	for {
		select {
		case <-time.After(time.Second / time.Duration(ratePerSec)):
			bucket <- struct{}{}
		case <-ctx.Done():
			return
		}
	}
}
