package main

import (
	"testing"

	"github.com/bluemir/go-utils/auth"
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
func TestRBACRuleTest(t *testing.T) {
	m := setupManager(t)
	token, err := m.Authn().Default("bluemir", "1234")
	if err != nil {
		t.Fatal(err)
	}
	m.RBAC().User(token.User).AddRole("test")

	m.RBAC().Rules().Role("test").AddPermission(auth.ActionGet, auth.Kind("article"))

	if !m.Authz(token).Get(auth.Kind("article")) {
		t.Fatal("must allow GET article")
	}

	m.RBAC().Rules().Role("test").DeletePermission(auth.ActionGet, auth.Kind("article"))

	if m.Authz(token).Get(auth.Kind("article")) {
		t.Fatal("must not allow GET article")
	}

}
