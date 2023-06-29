package timer_demo

import (
	"context"
	"fmt"
	"time"
)

type timerDemo struct {
	stopHeartBeat context.CancelFunc
}

func (td *timerDemo) timer() {
	var ctx context.Context

	ctx, td.stopHeartBeat = context.WithCancel(context.Background())
	go td.purge(ctx)

}

func (td *timerDemo) purge(ctx context.Context) {
	heartbeat := time.NewTicker(2 * time.Second)
	defer func() {
		heartbeat.Stop()
	}()
	for  {
		select {
		case <- heartbeat.C:
		case <- ctx.Done():
			return
		}
		fmt.Println("定时任务")
	}
}
