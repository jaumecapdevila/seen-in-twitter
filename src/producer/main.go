package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	nsqeueu "github.com/jaumecapdevila/seen-in-twitter/src/producer/queue"
	"github.com/jaumecapdevila/seen-in-twitter/src/producer/twitter"
	"github.com/spf13/viper"
)

var nsq *nsqeueu.NSQQueue

func init() {
	loadConfig()
	setupQueue()
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
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	tweets := make(chan twitter.Tweet)
	publisherStoppedChan := nsq.PublishVotes(tweets)
	twitterStoppedChan := twitter.StartStream(mongoDB, stopChan, tweets)
	go func() {
		<-signalChan
		stopChan <- struct{}{}
	}()
	<-twitterStoppedChan
	close(tweets)
	<-publisherStoppedChan
}
