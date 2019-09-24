package auth

import (
	"time"

	"github.com/pkg/errors"
)

type tokenClause struct {
	*manager
}

func (m *manager) Token() TokenClause {
	return &tokenClause{m}
}

func (m *tokenClause) Issue(username, unhashedKey string) (*Token, error) {
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
		User:      user,
		HashedKey: hash(unhashedKey, salt(username)),
		RevokeKey: hash(username+time.Now().String()+unhashedKey[:4], "__revoke__"),
	}

	err = m.store.CreateToken(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}
func (m *tokenClause) List(username string) ([]Token, error) {
	return m.store.ListToken(username)
}
func (m *tokenClause) Revoke(revokeKey string) error {
	return m.store.DeleteToken(revokeKey)
}
