package main

import (
	"testing"

	"github.com/bluemir/go-utils/auth"
	_ "github.com/bluemir/go-utils/auth/gorm"
)

func TestCreateUser(t *testing.T) {
	authManager, first, err := auth.New(&auth.Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !first {
		t.Fatalf("auth store not clear")
	}

	if _, err := authManager.User().Create(&auth.User{
		Name: "bluemir",
	}); err != nil {
		t.Fatal(err)
	}

	user, ok, err := authManager.User().Get("bluemir")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("user not found")
	}

	if user.Name != "bluemir" {
		t.Error("name not match")
	}
}

func TestCreateMultipleUser(t *testing.T) {
	authManager, first, err := auth.New(&auth.Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !first {
		t.Fatalf("auth store not clear")
	}

	if _, err := authManager.User().Create(&auth.User{
		Name: "bluemir",
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := authManager.User().Create(&auth.User{
		Name: "redmir",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := authManager.User().Create(&auth.User{
		Name: "blackmir",
	}); err != nil {
		t.Fatal(err)
	}

	user, ok, err := authManager.User().Get("bluemir")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("user not found")
	}

	if user.Name != "bluemir" {
		t.Error("name not match")
	}
}
func TestCreateToken(t *testing.T) {
	authManager, first, err := auth.New(&auth.Options{})
	if err != nil {
		t.Fatal(err)
	}

	if !first {
		t.Fatalf("auth store not clear")
	}

	if _, err := authManager.User().Create(&auth.User{
		Name: "bluemir",
	}); err != nil {
		t.Fatal(err)
	}

	token, err := authManager.Token().Issue("bluemir", "1234")
	if token.User.Name != "bluemir" {
		t.Error("token not matched")
	}
}
