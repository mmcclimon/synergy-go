package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type arbitraryJSON = map[string]interface{}

// Client is a slack client, with a bunch of internal fields and probably some
// methods.
type Client struct {
	apiKey    string
	OurName   string
	ourID     string
	connected bool
	ws        *websocket.Conn

	// these are all id: pretty-name
	usernames          map[string]string
	channels           map[string]string
	groupConversations map[string]string
	dmChannels         map[string]string // userid => dmChannelId
}

// Message represents a raw json message from slack
type Message struct {
	Type    string
	TS      string
	Subtype string
	Channel string
	Text    string
	User    string
}

type slackUser struct {
	ID   string
	Name string
}

type slackChannel struct {
	ID        string
	Name      string // not in DMs
	User      string // for DMs
	IsChannel bool   `mapstructure:"is_channel"`
	IsIM      bool   `mapstructure:"is_im"`
	IsGroup   bool   `mapstructure:"is_group"`
}

// NewClient gives you an instance of a slack client (not connected)
func NewClient(apiKey string) *Client {
	client := Client{
		apiKey:             apiKey,
		connected:          false,
		usernames:          make(map[string]string),
		channels:           make(map[string]string),
		dmChannels:         make(map[string]string),
		groupConversations: make(map[string]string),
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

	client.OurName = data.Self.Name
	client.ourID = data.Self.ID

	// connect to the rtm
	conn, _, err := websocket.DefaultDialer.Dial(data.URL, nil)

	if err != nil {
		log.Fatalf("could not connect to slack: %s", err)
	}

	client.ws = conn
	client.connected = true
	log.Println("connected to Slack")

	var wg sync.WaitGroup

	wg.Add(2)
	go client.loadUsers(&wg)
	go client.loadConversations(&wg)
	wg.Wait()
}

// Run is the central listen loop,
func (client *Client) Run(rawEvents chan<- Message) {
	if client.ws == nil {
		client.connect()
	}

	defer client.ws.Close()

	// is it bad form to reuse this? maybe!
	var msg Message

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

		// log.Printf("got message from user %s on channel %s: %s", client.usernameFor(msg.User), msg.Channel, msg.Text)
		rawEvents <- msg
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
func (client *Client) apiCall(req *http.Request) arbitraryJSON {
	req.Header.Set("Authorization", client.apiAuthHeader())

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

func (client *Client) apiCallForm(endpoint string, postData map[string]string) arbitraryJSON {
	vals := url.Values{}
	for key, val := range postData {
		vals.Add(key, val)
	}

	req, _ := http.NewRequest("POST", apiURL(endpoint), strings.NewReader(vals.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return client.apiCall(req)
}

func (client *Client) loadUsers(wg *sync.WaitGroup) {
	defer wg.Done()

	data := client.apiCallForm("users.list", nil)

	var users []slackUser

	err := mapstructure.Decode(data["members"], &users)
	if err != nil {
		log.Println("error loading users", err)
		return
	}

	for _, user := range users {
		client.usernames[user.ID] = user.Name
	}

	// fix up our own data
	client.usernames[client.ourID] = client.OurName

	log.Println("loaded users")
}

func (client *Client) loadConversations(wg *sync.WaitGroup) {
	defer wg.Done()

	var postData = map[string]string{
		"excludeArchived": "true",
		"types":           "public_channel,mpim,im",
	}

	data := client.apiCallForm("conversations.list", postData)

	var channels []slackChannel
	err := mapstructure.Decode(data["channels"], &channels)

	if err != nil {
		log.Println("error loading channels", err)
		return
	}

	for _, channel := range channels {
		id := channel.ID
		switch {
		case channel.IsChannel:
			client.channels[id] = channel.Name

		case channel.IsGroup:
			client.groupConversations[id] = channel.Name

		case channel.IsIM:
			client.dmChannels[channel.User] = id

		default:
			panic("unknown channel type!")
		}
	}

	log.Println("loaded conversations")
}

// DMChannelForAddress gives you the dm channel for a U12345 string
func (client *Client) DMChannelForAddress(userAddr string) (string, bool) {
	channel, ok := client.dmChannels[userAddr]

	if !ok {
		// TODO: missing a bunch of logic from real synergy, who opens it on
		// demand
		log.Printf("could not find dm channel for %s", userAddr)
	}

	return channel, ok
}

// UsernameFor returns the username for an address
func (client *Client) UsernameFor(addr string) string {
	return client.usernames[addr]
}
