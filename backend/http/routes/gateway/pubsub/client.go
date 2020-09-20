package pubsub

import (
	"log"
	"net/http"
	"sync"

	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	UserID facechat.ID
	Events chan interface{}

	sendMu sync.Mutex

	// nillabe
	conn *websocket.Conn
}

func NewClient(userID facechat.ID) *Client {
	return &Client{
		UserID: userID,
		Events: make(chan interface{}),
	}
}

func (c *Client) Connect(w http.ResponseWriter, r *http.Request) error {
	u, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.Wrap(err, "failed to upgrade")
	}

	c.conn = u
	return nil
}

func (c *Client) Close() error {
	close(c.Events)
	return c.conn.Close()
}

func (c *Client) Background() <-chan struct{} {
	var stop = make(chan struct{})

	go func() {
		for ev := range c.Events {
			e, err := NewEvent(ev)
			if err != nil {
				stop <- struct{}{}
				return
			}

			if err := c.conn.WriteJSON(e); err != nil {
				stop <- struct{}{}
				return
			}
		}
	}()

	go func() {
		for {
			_, b, err := c.conn.ReadMessage()
			if err != nil {
				stop <- struct{}{}
				return
			}

			log.Println("Received", string(b))
		}
	}()

	return stop
}
