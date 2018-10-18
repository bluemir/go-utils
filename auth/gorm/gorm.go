package gorm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"

	"github.com/bluemir/go-utils/auth"
)

func init() {
	auth.RegisterStoreDrver("gorm", New)
}

type store struct {
	db *gorm.DB
}

func getOpt(opts map[string]interface{}, name string, def interface{}) interface{} {
	if val, ok := opts[name]; ok {
		return val
	}
	return def
}
func tryGetDBFormOpt(opts map[string]interface{}) (*gorm.DB, bool) {
	value, ok := opts["db"]
	if !ok {
		return nil, false
	}
	db, ok := value.(*gorm.DB)
	if !ok {
		return nil, false
	}
	return db, true
}
func New(opts map[string]interface{}) (auth.StoreDriver, error) {
	db, ok := tryGetDBFormOpt(opts)
	if !ok {
		filename := getOpt(opts, "filename", "test.db").(string)
		var err error
		db, err = gorm.Open("sqlite3", filename)
		if err != nil {
			return nil, errors.New("failed to connect database")
		}
	}

	db.AutoMigrate(&auth.User{})
	db.AutoMigrate(&auth.Token{})
	db.AutoMigrate(&auth.Rule{})

	return &store{db}, nil
}
func (s *store) Close() error {
	return s.db.Close()
}

func (s *store) CreateUser(u *auth.User) error {
	// TODO vaildate
	if s.db.NewRecord(u) {
		err := s.db.Create(u).Error
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("User already exist")
}
func (s *store) GetUser(username string) (*auth.User, bool, error) {
	user := &auth.User{}
	if err := s.db.Where("name = ?", username).First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return user, true, nil
}
func (s *store) ListUser() ([]auth.User, error) {
	users := []auth.User{}
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (s *store) PutUser(u *auth.User) error {
	if err := s.db.Table("users").Where("name = ?", u.Name).Update(u).Error; err != nil {
		return err
	}
	return nil
}
func (s *store) DeleteUser(username string) error {
	if err := s.db.Where("name = ?", username).Delete(&auth.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) CreateToken(t *auth.Token) error {
	if s.db.NewRecord(t) {
		err := s.db.Create(t).Error
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Token already exist")
}
func (s *store) GetToken(username, hashedKey string) (*auth.Token, bool, error) {
	token := &auth.Token{}
	err := s.db.Where("username = ? AND hashed_key = ?", username, hashedKey).First(token).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return token, true, nil
}
func (s *store) ListToken(username string) ([]auth.Token, error) {
	result := []auth.Token{}
	err := s.db.Where("username = ?", username).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, err
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

// create mean just put
func (s *store) HasRule(role auth.Role, action auth.Action) bool {
	cnt := 0

	err := s.db.Table("rules").Where("role = ? AND allow = ?", role, action).Count(&cnt).Error
	if err != nil {
		//TODO log
		return false
	}
	if cnt > 0 {
		return true
	}
	return false

}
func (s *store) PutRule(role auth.Role, action auth.Action) error {
	rule := &auth.Rule{Role: role, Allow: action}
	if err := s.db.Where(rule).Assign(rule).FirstOrCreate(rule).Error; err != nil {
		return err
	}
	return nil
}
func (s *store) DeleteRule(role auth.Role, action auth.Action) error {
	if err := s.db.Where("role = ? AND action = ?", role, action).Delete(&auth.Rule{}).Error; err != nil {
		return err
	}

	return nil
}
