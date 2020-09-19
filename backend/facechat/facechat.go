package facechat

import "encoding/json"

type ID uint64

type User struct {
	ID    ID     `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`

	Accounts []Account `json:"accounts" db:"accounts"`
}

type Account struct {
	ID      string `json:"id" db:"id"`
	Service string `json:"service" db:"service"`
	Cookies json.RawMessage `json:"cookies" db:"cookies"`

	UserID ID `json:"userID" db:"user_id"`
}
