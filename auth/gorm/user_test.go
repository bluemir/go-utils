package gorm

import (
	"testing"

	"github.com/bluemir/go-utils/auth"
)

func TestUser(t *testing.T) {
	store, first, err := New(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	if !first {
		t.Fatalf("auth store not clear")
	}

	u := &auth.User{Name: "bluemir"}

	if err := store.CreateUser(u); err != nil {
		t.Error(err)
	}

	user, _, err := store.GetUser("bluemir")
	if err != nil {
		t.Error(err)
	}
	if user.Name != "bluemir" {
		t.Error("name not matched")
	}
}

func TestUserAttr(t *testing.T) {
	store, first, err := New(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	if !first {
		t.Fatalf("auth store not clear")
	}

	u := &auth.User{Name: "bluemir", Attr: map[string]string{"test1": "ok"}}

	if err := store.CreateUser(u); err != nil {
		t.Error(err)
	}

	user, _, err := store.GetUser("bluemir")
	if err != nil {
		t.Error(err)
	}
	if user.Name != "bluemir" {
		t.Error("name not matched")
	}
	if v, ok := user.Attr["test1"]; !ok {
		t.Error("attr1 'test1' not found")
	} else if v != "ok" {
		t.Error("test1's value not matched")
	}
}
func TestUserUpdate(t *testing.T) {
	store, first, err := New(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	if !first {
		t.Fatalf("auth store not clear")
	}
	{
		u := &auth.User{Name: "bluemir", Attr: map[string]string{"test1": "ok"}}

		if err := store.CreateUser(u); err != nil {
			t.Error(err)
		}
	}
	{
		u, _, err := store.GetUser("bluemir")
		if err != nil {
			t.Error(err)
		}

		u.Attr["test2"] = "user"
		if err := store.PutUser(u); err != nil {
			t.Fatal(err)
		}
		user, _, err := store.GetUser("bluemir")
		if err != nil {
			t.Error(err)
		}

		if v, ok := user.Attr["test2"]; !ok {
			t.Error("attr1 'test1' not found")
		} else if v != "user" {
			t.Error("test1's value not matched")
		}
	}
	{
		u, _, err := store.GetUser("bluemir")
		if err != nil {
			t.Error(err)
		}

		delete(u.Attr, "test1")
		if err := store.PutUser(u); err != nil {
			t.Fatal(err)
		}
		user, _, err := store.GetUser("bluemir")
		if err != nil {
			t.Error(err)
		}

		if _, ok := user.Attr["test1"]; ok {
			t.Error("attr1 'test1' found must be empty")
		}
	}
}
