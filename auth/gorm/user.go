package gorm

import (
	"github.com/bluemir/go-utils/auth"
	"github.com/jinzhu/gorm"
)

func (s *store) CreateUser(u *auth.User) error {
	// TODO vaildate
	user := fromAuthUser(u)

	err := s.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil

}
func (s *store) GetUser(username string) (*auth.User, bool, error) {
	user := &User{}
	if err := s.db.Preload("Attrs").Where(&User{Name: username}).Take(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	u := user.toAuthUser()
	return u, true, nil
}
func (s *store) ListUser() ([]auth.User, error) {
	users := []User{}
	if err := s.db.Preload("Attrs").Find(&users).Error; err != nil {
		return nil, err
	}
	result := []auth.User{}
	for _, user := range users {
		result = append(result, *user.toAuthUser())
	}

	return result, nil
}
func (s *store) PutUser(u *auth.User) error {
	// TODO primary key change
	old := &User{Name: u.Name}
	if err := s.db.Take(old, old).Error; err != nil {
		return err
	}
	user := fromAuthUser(u)

	if err := s.db.Model(old).Association("Attrs").Replace(&user.Attrs).Error; err != nil {
		return err
	}

	return nil
}
func (s *store) DeleteUser(username string) error {
	u := &User{Name: username}
	if err := s.db.Where(u).Delete(u).Error; err != nil {
		return err
	}
	return nil
}
