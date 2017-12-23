package queue

import (
	"fmt"
	"log"

	nsq "github.com/bitly/go-nsq"
)

// NSQQueue to publish votes to the nsq queue
type NSQQueue struct {
	ProducerConfig string
}

// New returns a new NSQQueue
func New(config string) *NSQQueue {
	return &NSQQueue{
		ProducerConfig: config,
	}
}

// PublishVotes sends the votes to the nsq queue
func (n *NSQQueue) PublishVotes(votes <-chan string) <-chan struct{} {
	stopChan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer(n.ProducerConfig, nsq.NewConfig())
	go func() {
		for vote := range votes {
			err := pub.Publish("votes", []byte(vote))
			if err != nil {
				fmt.Println("Publishing vote failed with error: ", err)
			}
		}
		log.Println("Publisher: Stopping")
		pub.Stop()
		log.Println("Publisher: Stopped")
		stopChan <- struct{}{}
	}()
	return stopChan
}
