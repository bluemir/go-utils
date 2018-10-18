package auth

type StoreDriver interface {
	CreateUser(u *User) error
	GetUser(username string) (*User, bool, error)
	ListUser() ([]User, error)
	PutUser(u *User) error
	DeleteUser(username string) error

	CreateToken(t *Token) error
	GetToken(username, hashedKey string) (*Token, bool, error)
	ListToken(username string) ([]Token, error)
	PutToken(t *Token) error
	DeleteToken(revokeKey string) error

	// create mean just put
	HasRule(role Role, action Action) bool // error?
	PutRule(role Role, action Action) error
	DeleteRule(role Role, action Action) error

	Close() error
}

// auth.Is(token).Allow("action")
type Manager interface {
	HttpAuth(header string) (*Token, error)
	Default(name, key string) (*Token, error)

	CreateUser(u *User) error
	GetUser(name string) (*User, bool, error)
	ListUser(filter ...string) ([]User, error)
	UpdateUser(u *User) error
	DeleteUser(u *User) error

	IssueToken(username, unhashedKey string) (*Token, error)
	ListToken(username string) ([]Token, error)
	UpdateToken(t *Token) error
	RevokeToken(RevokeKey string) error

	Is(interface{}) IsClause

	PutRule(role string, action ...string) error
	DeleteRule(role string, action ...string) error

	Root(name string) (string, error) // generate-key

	Close() error
}

type IsClause interface {
	Allow(action Action) bool
}
