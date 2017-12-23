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
	"github.com/spf13/viper"
)

var mongoDB *persistence.MongoDB

func init() {
	loadConfig()
	setupDatabase()
}

func setupDatabase() {
	if mongoDB, err := persistence.New(viper.Get("database.source")); err != nil {
		log.Fatal("Failed to establishing a connection with the database")
	}
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to load the configuration file")
	}
}

func main() {
	defer mongoDB.CloseConnection()
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
