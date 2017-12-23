package queue

import (
	"fmt"
	"log"

	nsq "github.com/bitly/go-nsq"
	"github.com/jaumecapdevila/seen-in-twitter/src/producer/twitter"
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
func (n *NSQQueue) PublishVotes(tweets <-chan twitter.Tweet) <-chan struct{} {
	stopChan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer(n.ProducerConfig, nsq.NewConfig())
	go func() {
		for tweet := range tweets {
			err := pub.Publish(tweet.Topic, []byte(tweet.Text))
			if err != nil {
				fmt.Println("Publishing tweet failed with error: ", err)
			}
			fmt.Println("Tweet published to queue: ", tweet)
		}
		log.Println("Publisher: Stopping")
		pub.Stop()
		log.Println("Publisher: Stopped")
		stopChan <- struct{}{}
	}()
	return stopChan
}
