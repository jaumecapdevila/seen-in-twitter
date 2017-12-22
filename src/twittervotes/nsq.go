package twittervotes

import (
	"log"

	nsq "github.com/bitly/go-nsq"
)

// PublishVotes sends the votes to the nsq queue
func PublishVotes(votes <-chan string) <-chan struct{} {
	stopChan := make(chan struct{}, 1)
	// TODO: Handle the error properly
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
