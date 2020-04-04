package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Client is a slack client, with a bunch of internal fields and probably some
// methods.
type Client struct {
	apiKey    string
	ourName   string
	ourID     string
	connected bool
	isReady   bool
	// websocket, plus internal data
}

// NewClient gives you an instance of a slack client (not connected)
func NewClient() *Client {
	apiKey := os.Getenv("SLACK_API_TOKEN")

	if len(apiKey) == 0 {
		log.Fatal("Error: must set SLACK_API_TOKEN in env!")
	}

	client := Client{
		apiKey:    apiKey,
		connected: false,
		isReady:   false,
	}

	return &client
}

// Connect sets up a connection to slack
func (client *Client) Connect() {
	u, _ := url.Parse("https://slack.com/api/rtm.connect")
	q := u.Query()
	q.Set("token", client.apiKey)
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	type slackConnection struct {
		Ok   bool                   `json:"ok"`
		URL  string                 `json:"url"`
		Team map[string]interface{} `json:"team"`
		Self struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"self"`
	}

	data := slackConnection{}
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		log.Fatalf("Could not decode JSON: %s", err)
	}

	fmt.Println(data)
}
