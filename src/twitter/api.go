package twitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jaumecapdevila/twitter-votes/src/persistence"
)

func read(db *persistence.MongoDB, votes chan<- string) {
	options, err := db.LoadOptions()
	if err != nil {
		log.Println("Failed to load options:", err)
		return
	}
	u, err := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
	if err != nil {
		log.Println("Creating filter request failed: ", err)
		return
	}
	query := make(url.Values)
	query.Set("track", strings.Join(options, ","))
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(query.Encode()))
	if err != nil {
		log.Println("Creating filter request failed: ", err)
	}
	resp, err := makeRequest(req, query)
	if err != nil {
		log.Println("Making request failed: ", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Printf("Invalid status code %d, 200 expected: ", resp.StatusCode)
		return
	}
	reader := resp.Body
	decoder := json.NewDecoder(reader)
	for {
		var t Tweet
		if err := decoder.Decode(&t); err != nil {
			log.Println("Decoding response failed with the following error: ", err)
			break
		}
		fmt.Println(t)
		for _, option := range options {
			if strings.Contains(
				strings.ToLower(t.Text),
				strings.ToLower(option),
			) {
				log.Println("Vote: ", option)
				votes <- option
			}
		}
	}
}

// StartStream opens stream with the twitter API
func StartStream(db *persistence.MongoDB, stopChan <-chan struct{}, votes chan<- string) <-chan struct{} {
	stoppedChan := make(chan struct{}, 1)
	go func() {
		defer func() {
			stoppedChan <- struct{}{}
		}()
		for {
			select {
			case <-stopChan:
				log.Println("Stopping Twitter stream")
				return
			default:
				log.Println("Querying Twitter...")
				read(db, votes)
				log.Println("(waiting)")
				time.Sleep(10 * time.Second) // wait before reconnecting
			}
		}
	}()
	return stoppedChan
}
