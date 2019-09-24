package auth

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/bluemir/go-utils/auth/utils"
)

type DriverInit func(opt map[string]interface{}) (StoreDriver, bool, error)

var drivers = map[string]DriverInit{}

func RegisterStoreDrver(name string, i DriverInit) {
	drivers[name] = i
}

func New(opt *Options) (Manager, bool, error) {
	if opt.RandomKeyLength == 0 {
		opt.RandomKeyLength = 32
	}

	if opt.StoreDriver == "" {
		opt.StoreDriver = "gorm"
	}

	driver, ok := drivers[opt.StoreDriver]
	if !ok {
		return nil, true, errors.Errorf("StoreDriver not found: '%s'", opt.StoreDriver)
	}
	sd, first, err := driver(opt.DriverOpts)
	if err != nil {
		return nil, first, err
	}

	// TODO predefined Role, Rule
	return &manager{
		store: sd,
		opt:   opt,
		rules: &Rules{
			items: []Rule{},
		},
	}, first, nil
}

type Options struct {
	StoreDriver     string
	DriverOpts      map[string]interface{}
	RandomKeyLength int
}

type manager struct {
	store StoreDriver
	opt   *Options
	rules *Rules
}

func (m *manager) Close() error {
	return m.store.Close()
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
			Attr: map[string]string{
				"rbac/role-root": "true",
			},
		}
		err := m.store.CreateUser(root)
		if err != nil {
			return nil
		}
	} else {

		//	m.store.SetUserAttr(root.ID)

		// m.store.PutUser(root)
	}

	tokens, err := m.store.ListToken(name)
	for _, token := range tokens {
		if err := m.store.DeleteToken(token.RevokeKey); err != nil {
			return err
		}
	}

	token := &Token{
		User:      root,
		HashedKey: hash(key, salt(name)),
		RevokeKey: hash(name+time.Now().String()+key[:4], "__revoke__"),
	}

	err = m.store.CreateToken(token)
	if err != nil {
		return err
	}
	return nil
}
