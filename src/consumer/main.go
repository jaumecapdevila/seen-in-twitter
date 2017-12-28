package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	nsq "github.com/bitly/go-nsq"
	"github.com/jaumecapdevila/seen-in-twitter/src/consumer/persistence"
	"github.com/spf13/viper"
)

var mongoDB *persistence.MongoDB

func init() {
	loadConfig()
	setupDatabase()
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to load the configuration file")
	}
}

func setupDatabase() {
	var err error
	if mongoDB, err = persistence.New(viper.GetString("database.host")); err != nil {
		log.Fatalf("Establishing a connection to the database failed with the following error: %s", err.Error())
	}
}

func main() {
	log.Println("Connecting to nsq...")
	q, err := nsq.NewConsumer("recetas", "recetas", nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
		return
	}
	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		tweet := string(m.Body)
		fmt.Println(tweet)
		return nil
	}))
	if err := q.ConnectToNSQLookupd(fmt.Sprintf("%s:%s", viper.GetString("nsqlookupd.host"), viper.GetString("nsqlookupd.port"))); err != nil {
		log.Fatal(err)
		return
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			q.Stop()
		case <-q.StopChan:
			return
		}
	}
}
