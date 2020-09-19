package facechat

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/pkg/errors"
)

var node *snowflake.Node

func init() {
	snowflake.Epoch = 1288834974657
	var err error
	node, err = snowflake.NewNode(0)
	if err != nil {
		panic(errors.Wrap(err, "error creating snowflake node"))
	}
}

type ID uint64

func GenerateID() ID {
	return ID(node.Generate())
}

func (id ID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

type User struct {
	ID    ID     `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

type Relationship struct {
	TargetID ID               `json:"target_id" db:"target_id"`
	Type     RelationshipType `json:"type"      db:"type"`
}

type RelationshipType uint8

const (
	Stranger RelationshipType = iota
	Blocked
	Friend
	IncomingFriendRequest
	SentFriendRequest
)

const MinAccounts = 1

var ErrNoAccountsLinked = httperr.New(400, "no accounts linked")

type Account struct {
	Service string          `json:"service" db:"service"`
	Name    string          `json:"name"    db:"name"`
	URL     string          `json:"url"     db:"url"`
	Data    json.RawMessage `json:"data"    db:"data"`
	UserID  ID              `json:"userID"  db:"user_id"`
}

const SessionTimeout = 7 * 24 * time.Hour

type Session struct {
	UserID ID        `json:"-" db:"user_id"`
	Token  string    `json:"-" db:"token"`
	Expiry time.Time `json:"-" db:"expiry"`
}

// RoomsPerUser is the maximum number of rooms a user can make. This may change
// in the future.
const RoomsPerUser = 3

const MaxRoomNameLen = 64

type Room struct {
	ID    ID          `json:"id"    db:"id"`
	Type  RoomType    `json:"type"  db:"typee"`
	Name  string      `json:"name"  db:"name"`
	Topic string      `json:"topic" db:"topic"`
	Level SecretLevel `json:"level" db:"level"`
}

type RoomType int8

const (
	PublicLobby RoomType = iota
	PrivateRoom
)

type SecretLevel int8

const (
	// Anonymous means that no username, avatar nor any personal information is
	// exposed. This is the equivalent of anonymous image boards.
	Anonymous SecretLevel = iota
	// HalfOpen exposes the username and avatar.
	HalfOpen
	// FullyOpen exposes all information, including social media accounts.
	FullyOpen
)

// MaxMessageLen is the maximum number of bytes per message.
const MaxMessageLen = 2048

type Message struct {
	ID       ID          `json:"id"        db:"id"`
	Type     MessageType `json:"type"      db:"type"`
	RoomID   ID          `json:"room_id"   db:"room_id"`
	AuthorID ID          `json:"author_id" db:"author_id"`
	Markdown string      `json:"markdown"  db:"markdown"`
}

type MessageType int8

const (
	NormalMessage MessageType = iota
	JoinMessage
	LeaveMessage
	DeletedMessage // soft delete
)
