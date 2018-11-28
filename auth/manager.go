package auth

import (
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/bluemir/go-utils/auth/codes"
	"github.com/bluemir/go-utils/auth/utils"
)

type DriverInit func(opt map[string]interface{}) (StoreDriver, bool, error)

var drivers = map[string]DriverInit{}

func RegisterStoreDrver(name string, i DriverInit) {
	drivers[name] = i
}

func New(opt *Options) (Manager, bool, error) {
	driver, ok := drivers[opt.StoreDriver]
	if !ok {
		return nil, true, errors.Errorf("StoreDriver not found: '%s'", opt.StoreDriver)
	}
	sd, first, err := driver(opt.DriverOpts)
	if err != nil {
		return nil, first, err
	}

	if opt.RandomKeyLength == 0 {
		opt.RandomKeyLength = 32
	}

	// TODO predefined Role, Rule
	return &manager{
		store: sd,
		opt:   opt,
	}, first, nil
}

type Options struct {
	DefaultRole     string
	RootRole        string
	StoreDriver     string
	DriverOpts      map[string]interface{}
	RandomKeyLength int
}

type manager struct {
	store StoreDriver
	opt   *Options
}

func (m *manager) Close() error {
	return m.store.Close()
}

func (m *manager) HttpAuth(header string) (*Token, error) {
	if header == "" {
		return nil, Errorf(codes.EmptyHeader, "Empty header")
	}
	arr := strings.SplitN(header, " ", 2)
	switch arr[0] {
	case "Basic", "basic", "Token", "token":
		str, err := base64.StdEncoding.DecodeString(arr[1])
		if err != nil {
			return nil, Error(codes.WrongEncoding, err)
		}

		authStr := strings.SplitN(string(str), ":", 2)
		if len(authStr) != 2 {
			return nil, Errorf(codes.BadToken, "Token invaildated")
		}
		return m.Default(authStr[0], authStr[1])

	case "Bearer", "bearer":
		// TODO
		return nil, Errorf(codes.NotImplement, "Not Implements")
	default:
		return nil, Errorf(codes.NotImplement, "Not Implements")
	}
	return nil, nil
}
func (m *manager) Default(name, pw string) (*Token, error) {
	if name == "" && pw == "" {
		return nil, Errorf(codes.EmptyAccount, "EmptyAccount")
	}
	token, ok, err := m.store.GetToken(name, hash(pw, "__salt__"+name+"__salt__"))
	if err != nil {
		return nil, Error(codes.Store, err)
	}
	if !ok {
		return nil, Errorf(codes.Unauthorized, "token not found")
	}
	return token, nil
}

func (m *manager) CreateUser(u *User) error {
	if u.Role == "" {
		u.Role = Role(m.opt.DefaultRole)
	}

	return m.store.CreateUser(u)
}
func (m *manager) GetUser(name string) (*User, bool, error) {
	return m.store.GetUser(name)
}
func (m *manager) ListUser(filter ...string) ([]User, error) {
	return m.store.ListUser()

}
func (m *manager) UpdateUser(u *User) error {
	// TODO vaildate
	return m.store.PutUser(u)
}
func (m *manager) DeleteUser(u *User) error {
	return m.store.DeleteUser(u.Name)
}

func (m *manager) IssueToken(username, unhashedKey string) (*Token, error) {
	//check user exist
	user, ok, err := m.store.GetUser(username)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("not found")
	}
	if user == nil {
		return nil, errors.New("User not found")
	}
	if len(unhashedKey) < 4 {
		return nil, errors.New("too short key")
	}

	token := &Token{
		Username:  username,
		HashedKey: hash(unhashedKey, salt(username)),
		RevokeKey: hash(username+time.Now().String()+unhashedKey[:4], "__revoke__"),
	}

	err = m.store.CreateToken(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}
func (m *manager) ListToken(username string) ([]Token, error) {
	return m.store.ListToken(username)
}
func (m *manager) UpdateToken(t *Token) error {
	// TODO vaildate
	return m.UpdateToken(t)
}
func (m *manager) RevokeToken(revokeKey string) error {
	return m.store.DeleteToken(revokeKey)
}
func (m *manager) Is(val interface{}) IsClause {
	switch v := val.(type) {
	case *User:
		return &isUserClause{m, v}
	case *Token:
		return &isTokenClause{m, v}
	case Role:
		return &isRoleClause{m, v}
	default:
		return &isDefaultClause{}
	}
}
func (m *manager) PutRule(role string, actions ...string) error {
	for _, action := range actions {
		m.store.PutRule(Role(role), Action(action))
	}
	return nil
}
func (m *manager) DeleteRule(role string, actions ...string) error {
	for _, action := range actions {
		m.store.DeleteRule(Role(role), Action(action))
	}
	return nil
}

// Ensure root account
// root account MUST have only one key. if forget it? recall this function
// recalling this function will be revoke all root's token except new one.
func (m *manager) Root(name string) (string, error) {
	unhashedKey := utils.RandomString(m.opt.RandomKeyLength)

	err := m.RootWithKey(name, unhashedKey)
	if err != nil {
		return "", nil
	}
	str := fmt.Sprintf("%s:%s", name, unhashedKey)
	return base64.StdEncoding.EncodeToString([]byte(str)), nil
}
func (m *manager) RootWithKey(name string, key string) error {
	root, ok, err := m.store.GetUser(name)
	if err != nil {
		return err
	}
	if !ok {
		root = &User{
			Name: name,
			Role: Role(m.opt.RootRole),
		}
		err := m.store.CreateUser(root)
		if err != nil {
			return nil
		}
	} else {
		root.Role = Role(m.opt.RootRole)

		m.store.PutUser(root)
	}

	tokens, err := m.store.ListToken(name)
	for _, token := range tokens {
		if err := m.store.DeleteToken(token.RevokeKey); err != nil {
			return err
		}
	}

	token := &Token{
		Username:  name,
		HashedKey: hash(key, salt(name)),
		RevokeKey: hash(name+time.Now().String()+key[:4], "__revoke__"),
	}

	err = m.store.CreateToken(token)
	if err != nil {
		return err
	}
	return nil
}

func hashRawHex(str string) string {
	hashed := crypto.SHA512.New()
	io.WriteString(hashed, str)
	return hex.EncodeToString(hashed.Sum(nil))
}
func hash(unhashedKey, saltSeed string) string {
	salt := hashRawHex(saltSeed)
	return hashRawHex(salt[:64] + unhashedKey + salt[64:])
}
func salt(username string) string {
	return "__salt__" + username + "__salt__"
}

const (
	HeaderAuthorization   = "Authorization"
	HeaderWWWAuthenticate = "WWW-Authenticate"
)

func HttpRealm(relam string) string {
	return `Basic realm="` + relam + `"`
}
