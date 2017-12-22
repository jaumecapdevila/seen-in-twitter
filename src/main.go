package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jaumecapdevila/twitter-votes/src/twittervotes"
)

func main() {
	var stopLock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		stopLock.Lock()
		stop = true
		stopLock.Unlock()
		log.Println("Stopping...")
		stopChan <- struct{}{}
		twittervotes.CloseConn()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
}
