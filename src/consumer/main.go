package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	nsq "github.com/bitly/go-nsq"
	"github.com/jaumecapdevila/seen-in-twitter/src/consumer/persistence"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var db *gorm.DB

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
	if db, err = gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local"); err != nil {
		log.Fatalf("Establish connection to database failed with error: %s", err.Error())
	}
}

func main() {
	defer db.Close()
	log.Println("Connecting to nsq...")
	q, err := nsq.NewConsumer("recetas", "recetas", nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
		return
	}
	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		tweet := string(m.Body)
		persistTweet(tweet)
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

func persistTweet(tweetBody string) {
	db.Create(&persistence.Tweet{
		Text: tweetBody,
	})
}
