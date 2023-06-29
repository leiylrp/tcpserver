package timer_demo

import (
	"os"
	"os/signal"
	"testing"
)

func TestTimerDemo(t *testing.T) {
	exitChan := make(chan os.Signal, 1)
	td := new(timerDemo)
	td.timer()

	signal.Notify(exitChan, os.Interrupt)
	<- exitChan
}
