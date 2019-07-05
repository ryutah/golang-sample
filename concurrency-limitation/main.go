package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const concurrencyLimit = 10

func main() {
	var (
		bucket      = make(chan struct{}, concurrencyLimit)
		wg          = new(sync.WaitGroup)
		doneChan    = make(chan struct{})
		ctx, cancel = context.WithCancel(context.Background())
	)

	go func() {
		var cnt int
		for {
			select {
			case <-doneChan:
				cnt++
				if cnt%10 == 0 {
					fmt.Printf("Finish %d proc\n", cnt)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	start := time.Now()
	for i := 0; i < 100; i++ {
		bucket <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-bucket
				wg.Done()
			}()
			time.Sleep(1 * time.Second)
			doneChan <- struct{}{}
		}()
	}

	wg.Wait()
	cancel()
	fmt.Println("take time: ", time.Since(start)) // should be take about 10 seconds.
	fmt.Println("finish")
}
