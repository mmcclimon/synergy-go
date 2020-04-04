package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	urllib "net/url"
	"os"
)

type Client struct {
	apiKey    string
	ourName   string
	ourID     string
	connected bool
	isReady   bool
	// websocket, plus internal data
}

func NewClient() (*Client, error) {
	apiKey := os.Getenv("SLACK_API_TOKEN")

	if len(apiKey) == 0 {
		return nil, errors.New("Error: must set SLACK_API_TOKEN in env!")
	}

	client := Client{
		apiKey:    apiKey,
		connected: false,
		isReady:   false,
	}

	return &client, nil
}

func (client *Client) Connect() {
	url, _ := urllib.Parse("https://slack.com/api/rtm.connect")
	q := url.Query()
	q.Set("token", client.apiKey)
	url.RawQuery = q.Encode()

	res, err := http.Get(url.String())

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	var data map[string]interface{}
	// body, _ := ioutil.ReadAll(res.Body)
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
}
