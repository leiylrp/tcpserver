package ipool

import (
	"time"
)

type IPool interface {
	// Submit submits a task to this pool.
	//
	// Note that you are allowed to call Pool.Submit() from the current Pool.Submit(),
	// but what calls for special attention is that you will get blocked with the latest
	// Pool.Submit() call once the current Pool runs out of its capacity, and to avoid this,
	// you should instantiate a Pool with ants.WithNonblocking(true).
	Submit(task func()) error

	// retrieveWorker returns an available worker to run the tasks.
	//retrieveWorker() (w *goWorker)

	// revertWorker puts a worker back into free pool, recycling the goroutines.
	//revertWorker(worker *goWorker) bool

	// Running returns the number of workers currently running.
	Running() int

	// Free returns the number of available goroutines to work, -1 indicates this pool is unlimited.
	Free() int

	// Waiting returns the number of tasks which are waiting be executed.
	Waiting() int

	// Cap returns the capacity of this pool.
	Cap() int

	// Tune changes the capacity of this pool, note that it is noneffective to the infinite or pre-allocation pool.
	Tune(size int)

	// IsClosed indicates whether the pool is closed.
	IsClosed() bool

	// Release closes this pool and releases the worker queue.
	Release()

	// ReleaseTimeout is like Release but with a timeout, it waits all workers to exit before timing out.
	ReleaseTimeout(timeout time.Duration) error

	// Reboot reboots a closed pool.
	Reboot()
}
