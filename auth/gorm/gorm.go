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
func New(opts map[string]interface{}) (auth.StoreDriver, bool, error) {
	db, ok := tryGetDBFormOpt(opts)
	if !ok {
		filename := getOpt(opts, "filename", ":memory:").(string)
		var err error
		db, err = gorm.Open("sqlite3", filename)
		if err != nil {
			return nil, true, errors.New("failed to connect database")
		}
		db.DB().SetMaxOpenConns(1)
	}

	first := true &&
		!db.HasTable(&User{}) && !db.HasTable(&UserAttr{}) &&
		!db.HasTable(&Token{}) && !db.HasTable(&TokenAttr{})

	db.AutoMigrate(
		&User{},
		&UserAttr{},
		&Token{},
		&TokenAttr{},
	)
	db.Model(&User{}).Association("Attrs")
	db.Model(&Token{}).Association("Attrs")

	return &store{db}, first, nil
}
func (s *store) Close() error {
	return s.db.Close()
}
