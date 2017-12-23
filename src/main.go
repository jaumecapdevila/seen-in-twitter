package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jaumecapdevila/twitter-votes/src/persistence"
	nsqeueu "github.com/jaumecapdevila/twitter-votes/src/queue"
	"github.com/jaumecapdevila/twitter-votes/src/twitter"
	"github.com/spf13/viper"
)

var mongoDB *persistence.MongoDB
var nsq *nsqeueu.NSQQueue

func init() {
	loadConfig()
	setupDatabase()
	setupQueue()
}

func setupDatabase() {
	var err error
	if mongoDB, err = persistence.New(viper.GetString("database.source")); err != nil {
		log.Fatalf("Establishing a connection to the database failed with the following error: %s", err.Error())
	}
}

func setupQueue() {
	nsq = nsqeueu.New(fmt.Sprintf("%s:%s", viper.GetString("queue.producer"), viper.GetString("queue.port")))
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
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	votes := make(chan string)
	publisherStoppedChan := nsq.PublishVotes(votes)
	twitterStoppedChan := twitter.StartStream(mongoDB, stopChan, votes)
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			twitter.CloseConn()
		}
	}()
	<-twitterStoppedChan
	close(votes)
	votes <- "bitcoin"
	<-publisherStoppedChan
}
