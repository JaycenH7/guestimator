package client

import (
	"log"

	"github.com/mrap/guestimator/server/match/event"

	"golang.org/x/net/websocket"
)

type Client struct {
	PlayerID       string
	ws             *websocket.Conn
	ReceivedEvents chan event.Event
}

const RECV_BUFFER_SIZE = 1000

func NewClient(playerID string) *Client {
	cli := Client{
		PlayerID:       playerID,
		ReceivedEvents: make(chan event.Event, RECV_BUFFER_SIZE),
	}
	return &cli
}

func (c *Client) Connect(urlStr string) error {
	ws, err := websocket.Dial(urlStr, "", "http://localhost/")
	if err != nil {
		log.Println("Could not create client connection with url", urlStr, err)
		return err
	}

	c.ws = ws

	go c.receiveLoop()
	return err
}

func (c *Client) receiveLoop() {
	for {
		var ev event.Event
		err := websocket.JSON.Receive(c.ws, &ev)
		// If there's an error here, we're assuming the ws is disconnected.
		if err != nil {
			log.Println("Error receiving on client ws", err)
			return
		}
		c.ReceivedEvents <- ev
	}
}
