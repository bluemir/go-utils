package main

import (
	"testing"

	"github.com/bluemir/go-utils/auth"
	_ "github.com/bluemir/go-utils/auth/gorm"
)

func TestBasicRBAC(t *testing.T) {
	authManager, first, err := auth.New(&auth.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !first {
		t.Fatalf("auth store not clear")
	}
	{
		authManager.RBAC().Rules().Role("reader").AddPermission(auth.ActionGet, auth.Kind("article"))

		user, err := authManager.User().Create(&auth.User{
			Name: "bluemir",
		})
		if err != nil {
			t.Fatal(err)
		}

		authManager.RBAC().User(user).AddRole("reader")

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
			t.Error("user not allowed read")
		}
	}
}
