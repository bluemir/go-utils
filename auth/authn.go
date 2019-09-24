package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/bluemir/go-utils/auth/codes"
)

type authnClause struct {
	*manager
}

func (m *manager) Authn() AuthnClause {
	return &authnClause{m}
}

func (m *authnClause) Default(name, unhashedKey string) (*Token, error) {
	if name == "" && unhashedKey == "" {
		return nil, Errorf(codes.EmptyAccount, "EmptyAccount")
	}
	token, ok, err := m.store.GetToken(name, hash(unhashedKey, salt(name)))
	if err != nil {
		return nil, Error(codes.Store, err)
	}
	if !ok {
		return nil, Errorf(codes.Unauthorized, "token not found")
	}
	return token, nil

}
func (m *authnClause) HTTP(header http.Header) (*Token, error) {
	return m.HTTPHeaderString(header.Get(HeaderAuthorization))
}
func (m *authnClause) HTTPHeaderString(header string) (*Token, error) {
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
}
