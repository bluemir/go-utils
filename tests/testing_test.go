package main

import (
	"testing"

	"github.com/bluemir/go-utils/auth"
	_ "github.com/bluemir/go-utils/auth/gorm"
)

func setupManager(t *testing.T) auth.Manager {
	authManager, first, err := auth.New(&auth.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !first {
		t.Fatalf("auth store not clear")
	}

	user, err := authManager.User().Create(&auth.User{
		Name: "bluemir",
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
	return authManager
}
