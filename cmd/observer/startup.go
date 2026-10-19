package main

import (
	"context"
	"os"
)

func openHardware(signalChan <-chan os.Signal, open func(context.Context) (context.Context, context.CancelFunc, error)) (context.Context, context.CancelFunc, bool, error) {
	openCtx, cancelOpen := context.WithCancel(context.Background())
	openDone := make(chan struct{})
	interruptedChan := make(chan bool, 1)
	go func() {
		select {
		case <-signalChan:
			cancelOpen()
			interruptedChan <- true
		case <-openDone:
			interruptedChan <- false
		}
	}()
	hardwareCtx, abort, err := open(openCtx)
	close(openDone)
	return hardwareCtx, abort, <-interruptedChan, err
}
