package pubsub

import (
	"sync"

	"github.com/cskr/pubsub"
	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/facechat"
)

const Capacity = 10

type Collection struct {
	*pubsub.PubSub

	db *db.DB

	cliMutex sync.RWMutex
	clients  map[facechat.ID]*Client
}

func NewCollection(db *db.DB) *Collection {
	return &Collection{
		PubSub:  pubsub.New(Capacity),
		clients: map[facechat.ID]*Client{},
	}
}

// Register registers a client.
func (c *Collection) Register(client *Client) {
	c.cliMutex.Lock()
	c.clients[client.UserID] = client
	c.cliMutex.Unlock()
}

// Unregister unsubscribes the given client, then removes it from the
// collection.
func (c *Collection) Unregister(userID facechat.ID) {
	c.cliMutex.Lock()

	client, ok := c.clients[userID]
	if ok {
		delete(c.clients, userID)
	}

	c.cliMutex.Unlock()

	c.Unsub(client.Events)
}

// BroadcastMessage broadcasts the given message to room ID topic.
func (c *Collection) BroadcastMessage(msg facechat.Message) {
	c.Pub(msg, msg.RoomID.String())
}

// SubscribeRoom subscribes the user with the ID to a room.
func (c *Collection) SubscribeRoom(userID, roomID facechat.ID) {
	client := c.getClient(userID)
	if client == nil {
		return
	}

	c.AddSub(client.Events, roomID.String())
}

// UnsubscribeRoom unsubscribes the user with the ID to a room.
func (c *Collection) UnsubscribeRoom(userID, roomID facechat.ID) {
	client := c.getClient(userID)
	if client == nil {
		return
	}

	c.Unsub(client.Events, roomID.String())
}

func (c *Collection) getClient(userID facechat.ID) *Client {
	c.cliMutex.RLock()
	cl, _ := c.clients[userID]
	c.cliMutex.RUnlock()
	return cl
}
