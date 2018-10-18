package auth

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

/*
User:Role   = n:1
User:Token  = 1:n
Role:Action = n:n
*/

type User struct {
	gorm.Model
	Name   string `gorm:"unique"`
	Role   Role   `sql:"type:string"`
	Labels Labels `sql:"type:json"`
}

type Token struct {
	gorm.Model
	Username  string
	HashedKey string `json:"-"`
	RevokeKey string
	Labels    Labels `sql:"type:json"`
	// if nil it mean same as user
	Allows Allows `sql:"type:json"`
}
type Allows []Action // for gorm ...
func (allows *Allows) Scan(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		return errors.New("must []byte")
	}
	return json.Unmarshal(str, allows)
}
func (allows Allows) Value() (driver.Value, error) {
	return json.Marshal(allows)
}

type Role string
type Action string

type Rule struct {
	Role  Role   `sql:"type:string"`
	Allow Action `sql:"type:string"`
}

// TODO labels gorm bindig
type Labels map[string]string

func (labels *Labels) Scan(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		return errors.New("must []byte")
	}
	err := json.Unmarshal(str, labels)
	if err != nil {
		return err
	}
	return nil
}
func (labels Labels) Value() (driver.Value, error) {
	return json.Marshal(labels)
}
