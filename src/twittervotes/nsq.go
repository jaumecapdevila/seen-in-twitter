package twittervotes

import (
	"log"

	nsq "github.com/bitly/go-nsq"
)

func publishVotes(votes <-chan string) <-chan struct{} {
	stopChan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer("nsqd:4150", nsq.NewConfig())
	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote))
		}
		log.Println("Publisher: Stopping")
		pub.Stop()
		log.Println("Publisher: Stopped")
		stopChan <- struct{}{}
	}()
	return stopChan
}
