package main

import (
	"testing"

	_ "github.com/bluemir/go-utils/auth/gorm"
)

func TestRBACRoleTest(t *testing.T) {
	m := setupManager(t)
	token, err := m.Authn().Default("bluemir", "1234")
	if err != nil {
		t.Fatal(err)
	}

	m.RBAC().User(token.User).AddRole("admin")
	m.RBAC().User(token.User).AddRole("test")

	roles, err := m.RBAC().User(token.User).Roles()
	if err != nil {
		t.Fatal(err)
	}

	for _, role := range roles {
		if role == "admin" || role == "test" {
			continue
		}
		t.Errorf("role is wrong")
	}

	m.RBAC().User(token.User).RemoveRole("test")

	roles, err = m.RBAC().User(token.User).Roles()
	if err != nil {
		t.Fatal(err)
	}
	for _, role := range roles {
		if role == "admin" {
			continue
		}
		t.Errorf("role is wrong: was %s", role)
	}
}
