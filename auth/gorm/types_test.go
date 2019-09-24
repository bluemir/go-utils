package gorm

import (
	"testing"

	"github.com/bluemir/go-utils/auth"
)

func TestFromAuthUser(t *testing.T) {
	{
		user := fromAuthUser(&auth.User{Name: "user"})

		if user.Name != "user" {
			t.Error("name not matched")
		}
	}
	{
		user := fromAuthUser(&auth.User{
			Name: "user",
			Attr: map[string]string{
				"test": "ok",
			},
		})

		if user.Name != "user" {
			t.Error("name not matched")
		}

		if user.Attrs[0].Key != "test" {
			t.Error("test attr not found")
		}

		if user.Attrs[0].Value != "ok" {
			t.Error("test attr not matched")
		}
	}
}
func TestToAuthUser(t *testing.T) {
	{
		u := User{
			Name: "user",
			Attrs: []UserAttr{
				{Key: "test", Value: "ok"},
			},
		}
		user := u.toAuthUser()

		if user.Name != "user" {
			t.Error("name not matched")
		}

		if v, ok := user.Attr["test"]; !ok {
			t.Error("attr not found")
		} else if v != "ok" {
			t.Error("attr not matched")
		}
	}
}
