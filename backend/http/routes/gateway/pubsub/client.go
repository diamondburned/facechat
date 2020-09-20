package pubsub

import (
	"net/http"
	"sync"

	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type Ready struct {
}

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
	return c.conn.Close()
}

func (c *Client) Start() error {
	for ev := range c.Events {
		e, err := NewEvent(ev)
		if err != nil {
			return errors.Wrap(err, "failed to make event")
		}

		if err := c.conn.WriteJSON(e); err != nil {
			return errors.Wrap(err, "failed to write JSON")
		}
	}

	return nil
}
