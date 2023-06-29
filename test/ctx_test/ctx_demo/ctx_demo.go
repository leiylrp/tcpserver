package ctx_demo

import (
	"context"
	"fmt"
	"time"
)

func ctxDemo() {
	ctx := context.Background()

	// 父无法访问子内容 兄弟context也无法相互访问
	parent := context.WithValue(ctx, "my key", "my value")
	sub := context.WithValue(ctx, "my  new key", "my new value")

	fmt.Println(parent.Value("my key"))
	fmt.Println(sub.Value("my key"))
}

// timeoutDemo 控制，父context的超时会决定子context
func timeoutDemo() {
	bg := context.Background()

	parent, cancelp := context.WithTimeout(bg, 1*time.Second)
	sub, cancels := context.WithTimeout(parent, 3*time.Second)

	go func() {
		<- sub.Done()
		fmt.Println("timeout")
	}()
	time.Sleep(5*time.Second)
	cancelp()
	cancels()
}

func timeoutControl() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	bsChan := make(chan struct{})

	go func() {
		slowBusiness()
		bsChan <- struct{}{}
	}()

	select {
	case <- ctx.Done():
		fmt.Println("ctx control...")
	case <- bsChan:
		fmt.Println("business control...")
	}
}

func slowBusiness() {
	time.Sleep(2*time.Second)
}

func timeAfterFunc() {
	bsChan := make(chan struct{})
	go func() {
		slowBusiness()
		bsChan <- struct{}{}
	}()

	time.AfterFunc(time.Second, func() {
		fmt.Println("time after func timeout...")
	})
	<- bsChan
	fmt.Println("business end...")

}