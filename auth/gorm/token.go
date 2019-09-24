package gorm

import (
	"github.com/bluemir/go-utils/auth"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func (s *store) CreateToken(t *auth.Token) error {
	token := fromAuthToken(t)

	if err := s.db.Create(token).Error; err != nil {
		return errors.Wrapf(err, "Token already exist")
	}
	return nil

}
func (s *store) GetToken(username, hashedKey string) (*auth.Token, bool, error) {
	token := &Token{
		Username:  username,
		HashedKey: hashedKey,
	}

	if err := s.db.Preload("Attrs").Where(token).Take(token).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	t, err := s.toAuthToken(token)
	if err != nil {
		return nil, false, err
	}
	return t, true, nil
}
func (s *store) ListToken(username string) ([]auth.Token, error) {
	result := []Token{}
	if err := s.db.Preload("Attrs").Where(&Token{Username: username}).Find(&result).Error; err != nil {
		return nil, err
	}

	tokens := []auth.Token{}
	for _, token := range result {
		t, err := s.toAuthToken(&token)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, *t)
	}

	return tokens, nil
}
func (s *store) PutToken(t *auth.Token) error {
	return nil
}
func (s *store) DeleteToken(revokeKey string) error {
	err := s.db.Where("revoke_key = ?", revokeKey).Delete(&auth.Token{}).Error
	if err != nil {
		return err
	}
	return nil
}
