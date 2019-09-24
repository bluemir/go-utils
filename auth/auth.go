package auth

import "net/http"

type StoreDriver interface {
	CreateUser(u *User) error
	GetUser(name string) (*User, bool, error)
	ListUser() ([]User, error)
	PutUser(u *User) error
	DeleteUser(name string) error

	CreateToken(t *Token) error
	GetToken(username, hashedKey string) (*Token, bool, error)
	ListToken(username string) ([]Token, error)
	PutToken(t *Token) error
	DeleteToken(revokeKey string) error

	//AddRule()
	//DeleteRule()
	//ListRule()

	Close() error
}

type Manager interface {
	User() UserClause
	Token() TokenClause

	Authn() AuthnClause
	Authz(token *Token) AuthzClause

	ABAC() ABACClause
	RBAC() RBACClause

	Rules() RulesClause

	Root(name string) (string, error) // generate-key
	RootWithKey(name string, key string) error

	Close() error
}
type UserClause interface {
	Create(u *User) (*User, error)
	Get(name string) (*User, bool, error)
	List(filter ...string) ([]User, error)
	Update(u *User) error
	Delete(u *User) error
}
type TokenClause interface {
	Issue(username, unhashedKey string) (*Token, error)
	List(username string) ([]Token, error)
	Revoke(revokeKey string) error
}
type AuthnClause interface {
	Default(name, unhashedKey string) (*Token, error)
	HTTPHeaderString(header string) (*Token, error)
	HTTP(header http.Header) (*Token, error)
}
type AuthzClause interface {
	Do(action Action, resource Resource) bool
	Create(resource Resource) bool
	Get(resource Resource) bool
	List(resource Resource) bool
	Update(resource Resource) bool
	Patch(resource Resource) bool
	Delete(resource Resource) bool
}

type ABACClause interface {
	User(user *User) AttrHandler
	Token(token *Token) AttrHandler
}

type RBACClause interface {
	Rules() RBACRulesClause
	User(user *User) RBACUserClause
}
type RBACRulesClause interface {
	Role(role Role) RBACRulesRoleClause
}
type RBACRulesRoleClause interface {
	//LoadYaml(buf []byte) error
	//LoadJson(buf []byte) error
	//DumpYaml() ([]byte, error)
	//DumpJson() ([]byte, error)
	AddPermission(action Action, resource Resource) error
	DeletePermission(action Action, resource Resource) error
}
type RBACUserClause interface {
	AddRole(role Role) error
	RemoveRole(role Role) error
	Roles() ([]Role, error)
}

type RulesClause interface {
	Add(rule string) error
	//AddCustom(id string, f func() bool) error
	List() ([]Rule, error)
	Delete(index int) error
}

type AttrHandler interface {
	Set(key, value string) error
	Get(key string) (string, error)
}
