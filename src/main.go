package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jaumecapdevila/twitter-votes/src/persistence"
	"github.com/jaumecapdevila/twitter-votes/src/queue"
	"github.com/jaumecapdevila/twitter-votes/src/twitter"
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
		twitter.CloseConn()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := persistence.DialDB(); err != nil {
		log.Fatalln("failed to dial MongoDB:", err)
	}
	defer persistence.CloseDB()
	votes := make(chan string) // chan for votes
	publisherStoppedChan := queue.PublishVotes(votes)
	twitterStoppedChan := twitter.StartTwitterStream(stopChan, votes)
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			twitter.CloseConn()
			stopLock.Lock()
			if stop {
				stopLock.Unlock()
				return
			}
			stopLock.Unlock()
		}
	}()
	<-twitterStoppedChan
	close(votes)
	<-publisherStoppedChan
}
