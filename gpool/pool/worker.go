package pool

import (
	"fmt"
	"runtime"
	"time"
)

// goWorker is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type goWorker struct {
	Parent 		*Pool

	task 		chan func()

	recycleTime	time.Time
}

// run 执行goWorker关联的任务
func (w *goWorker) run() {
	w.Parent.addRunning(1)
	fmt.Println("running task numbers:", w.Parent.Running())
	go func() {
		defer func() {
			w.Parent.addRunning(-1)
			fmt.Println("---------Put Pool")
			w.Parent.workerCache.Put(w)
			if p := recover(); p != nil {
				if ph := w.Parent.options.PanicHandler; ph != nil {
					ph(p)
				} else {
					w.Parent.options.Logger.Printf("worker exits from a panic: %v\n", p)
					var buf [4096]byte
					n := runtime.Stack(buf[:], false)
					w.Parent.options.Logger.Printf("worker exits from panic: %s\n", string(buf[:n]))
				}
			}
			// Call Signal() here in case there are goroutines waiting for available workers.
			w.Parent.cond.Signal()
		}()

		for f := range w.task {
			if f == nil {
				return
			}
			// 执行任务
			fmt.Println("exec task before...")
			f()
			fmt.Println("exec task after...")
			if ok := w.Parent.revertWorker(w); !ok {
				return
			}
		}
	}()
	fmt.Println("run task numbers: ", w.Parent.Running())
}
