package main

import (
	"testing"

	"github.com/bluemir/go-utils/auth"
	_ "github.com/bluemir/go-utils/auth/gorm"
)

func TestBasicABAC(t *testing.T) {
	authManager, first, err := auth.New(&auth.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !first {
		t.Fatalf("auth store not clear")
	}
	{
		authManager.Rules().Add(`User.Attr.role=="admin"`)

		user, err := authManager.User().Create(&auth.User{
			Name: "bluemir",
			Attr: map[string]string{
				"role": "admin",
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		token, err := authManager.Token().Issue(user.Name, "1234")
		if err != nil {
			t.Fatal(err)
		}
		if token == nil {
			t.Fatal("token is empty")
		}
	}
	{

		token, err := authManager.Authn().Default("bluemir", "1234")
		if err != nil {
			t.Fatal(err)
		}

		if !authManager.Authz(token).Get(auth.Kind("article")) {
			t.Error("user not allowed Get article")
		}
	}
}
