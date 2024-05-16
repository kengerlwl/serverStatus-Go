package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	// 新建一个上下文，ctx，就是context实例。 cancel,是一个函数， 就是ctx取消的函数
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel() // 执行这个，ctx.Done()会批量取消阻塞

	go handle(ctx, 500*time.Millisecond)
	// select {
	// case <-ctx.Done():
	// 	fmt.Println("main", ctx.Err())
	// }
	time.Sleep(1 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
}

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}
