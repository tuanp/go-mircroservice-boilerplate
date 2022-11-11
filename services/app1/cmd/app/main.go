package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	_, ctxCancel := context.WithCancel(context.Background())
	var teardownTimeout time.Duration = 15

	go waitForShutdownSignal(teardownTimeout, ctxCancel)
}

func waitForShutdownSignal(timeout time.Duration, callback func()) {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	<-osSignal
	// Wait for maximum 15s
	go func() {
		timer := time.NewTimer(timeout)
		<-timer.C
		//ll.Fatal("Force shutdown due to timeout!")
	}()
	callback()
}
