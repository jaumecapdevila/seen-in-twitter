package twitter

// Tweet struct represents a tweet from twitter
type Tweet struct {
	Topic string
	Text  string `json:"text"`
}
