package facechat

import "encoding/json"

type ID uint64

type User struct {
	ID    ID     `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	Accounts []Account `json:"accounts"`
}

type UserAuth struct {
	User
	PassHash string `json:"-"`
}

type Account struct {
	ID      string `json:"id" db:"id"`
	Service string `json:"service"`
	Cookies json.RawMessage

	UserID ID
}
