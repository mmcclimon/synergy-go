package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
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
	usernames map[string]string
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
		usernames: make(map[string]string),
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

	var wg sync.WaitGroup

	wg.Add(3)
	go client.loadUsers(&wg)
	go client.loadChannels(&wg)
	go client.loadDMs(&wg)
	wg.Wait()

	log.Println("loaded all the things!")
	client.isReady = true
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

	// is it bad form to reuse this? maybe!
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

		log.Printf("got message from user %s on channel %s: %s", client.usernameFor(msg.User), msg.Channel, msg.Text)
	}
}

func (client *Client) usernameFor(id string) string {
	username, ok := client.usernames[id]
	if !ok {
		log.Printf("no user found for id %s??", id)
		return id
	}

	return username
}

func apiURL(method string) string {
	return fmt.Sprintf("https://slack.com/api/%s", method)
}

func (client *Client) apiAuthHeader() string {
	return fmt.Sprintf("Bearer %s", client.apiKey)
}

// probably this should also return an err, but will ignore for now
func (client *Client) apiCall(endpoint string, postData arbitraryJSON) arbitraryJSON {
	data, err := json.Marshal(postData)
	if err != nil {
		log.Fatal("couldn't encode json: ", err)
	}

	req, _ := http.NewRequest("POST", apiURL(endpoint), bytes.NewBuffer(data))
	req.Header.Set("Authorization", client.apiAuthHeader())
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("error talking to Slack:", err)
	}

	var ret arbitraryJSON
	err = json.NewDecoder(res.Body).Decode(&ret)

	if err != nil {
		log.Fatal("invalid json from slack", err)
	}

	return ret
}

func (client *Client) loadUsers(wg *sync.WaitGroup) {
	defer wg.Done()

	postData := make(arbitraryJSON)
	postData["presence"] = false

	data := client.apiCall("users.list", postData)

	type user struct {
		ID   string
		Name string
	}

	var users []user

	err := mapstructure.Decode(data["members"], &users)
	if err != nil {
		log.Println("error loading users", err)
		return
	}

	for _, user := range users {
		client.usernames[user.ID] = user.Name
	}

	// fix up our own data
	client.usernames[client.ourID] = client.ourName

	log.Println("loaded users")
}

func (client *Client) loadChannels(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(2 * time.Second)
	log.Println("loaded channels")
}

func (client *Client) loadDMs(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(1 * time.Second)
	log.Println("loaded DMs")
}
