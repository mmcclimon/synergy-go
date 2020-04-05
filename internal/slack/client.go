package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

type arbitraryJSON = map[string]interface{}

// Client is a slack client, with a bunch of internal fields and probably some
// methods.
type Client struct {
	apiKey    string
	ourName   string
	ourID     string
	connected bool
	isReady   bool
	ws        *websocket.Conn
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

func (client *Client) connect() {
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
		Ok   bool
		URL  string
		Team arbitraryJSON
		Self struct {
			ID   string
			Name string
		}
	}

	data := slackConnection{}
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		log.Fatalf("Could not decode JSON: %s", err)
	}

	client.ourName = data.Self.Name
	client.ourID = data.Self.ID

	// connect to the rtm
	conn, _, err := websocket.DefaultDialer.Dial(data.URL, nil)

	if err != nil {
		log.Fatalf("could not connect to slack: %s", err)
	}

	client.ws = conn
	log.Println("connected to Slack!")
}

// Run is the central listen loop,
func (client *Client) Run() {
	if client.ws == nil {
		client.connect()
	}

	defer client.ws.Close()

	type slackMessage struct {
		Type    string
		TS      string
		Subtype string
		Channel string
		Text    string
		User    string
	}

	var msg slackMessage

	for {
		err := client.ws.ReadJSON(&msg)

		if err != nil {
			log.Println("error on read: ", err)
			return
		}

		// only handle message types
		if msg.Type != "message" {
			continue
		}

		log.Printf("got message from user %s on channel %s: %s", msg.User, msg.Channel, msg.Text)
	}
}
