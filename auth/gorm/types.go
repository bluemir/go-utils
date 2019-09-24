package gorm

import (
	"github.com/bluemir/go-utils/auth"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Attrs []UserAttr
}
type UserAttr struct {
	gorm.Model
	UserID uint
	Key    string
	Value  string
}
type Token struct {
	gorm.Model
	Username  string
	HashedKey string
	RevokeKey string
	Attrs     []TokenAttr
}
type TokenAttr struct {
	gorm.Model
	TokenID uint
	Key     string
	Value   string
}

func fromAuthUser(u *auth.User) *User {
	user := &User{
		Name: u.Name,
	}

	for k, v := range u.Attr {
		user.Attrs = append(user.Attrs, UserAttr{
			Key: k, Value: v,
		})
	}

	return user
}
func (u *User) toAuthUser() *auth.User {
	user := &auth.User{
		Name: u.Name,
		Attr: map[string]string{},
	}
	for _, attr := range u.Attrs {
		user.Attr[attr.Key] = attr.Value
	}
	return user
}
func fromAuthToken(t *auth.Token) *Token {
	token := &Token{
		Username:  t.User.Name,
		HashedKey: t.HashedKey,
		RevokeKey: t.RevokeKey,
	}
	for k, v := range t.Attr {
		token.Attrs = append(token.Attrs, TokenAttr{
			Key: k, Value: v,
		})
	}
	return token
}
func (s *store) toAuthToken(t *Token) (*auth.Token, error) {
	user, _, err := s.GetUser(t.Username)
	if err != nil {
		return nil, err
	}
	token := &auth.Token{
		User:      user,
		HashedKey: t.HashedKey,
		RevokeKey: t.RevokeKey,
		Attr:      map[string]string{},
	}
	for _, attr := range t.Attrs {
		token.Attr[attr.Key] = attr.Value
	}
	return token, nil
}
