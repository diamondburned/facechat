package facechat

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

type Account interface {
	ServiceName() string
	ID() string
	Name() string
	Email() string
}
