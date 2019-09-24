package auth

import (
	"strings"

	"github.com/pkg/errors"
)

// manager.RBAC().Rule().Role("admin").Can(Get).Resource(Type("article"))
// manager.RBAC().User(token).SetRole(role)
// manager.RBAC().Can(token).Create(resource)
// manager.RBAC().Can(token).Get(resource)
// manager.RBAC().Can(token).List(resource)
// manager.RBAC().Can(token).Update(resource)
// manager.RBAC().Can(token).Patch(resource)
// manager.RBAC().Can(token).Delete(resource)

type rbacClause struct {
	*manager
}

func (m *manager) RBAC() RBACClause {
	return &rbacClause{m}
}

func (m *rbacClause) Rules() RBACRulesClause {
	return &rbacRulesClause{m.manager}
}
func (m *rbacClause) User(user *User) RBACUserClause {
	return &rbacUserClause{m.manager, user}
}

type rbacRulesClause struct {
	*manager
	// Role(role Role) RBACCanClause
}

func (m *rbacRulesClause) Role(role Role) RBACRulesRoleClause {
	return &rbacRulesRoleClause{m.manager, role}
}

type rbacRulesRoleClause struct {
	*manager
	role Role
}

func (c *rbacRulesRoleClause) AddPermission(action Action, resource Resource) error {
	rules := []string{}
	for k, v := range resource.List() {
		rules = append(rules, `Resource.Get("`+k+`") == "`+v+`"`)
	}
	rules = append(rules, `Token.User.Attr["`+AttrKeyRolePrefix+c.role+`"]=="true"`)
	rules = append(rules, `Action == "`+action+`"`)

	return c.manager.rules.Add(strings.Join(rules, " && "))
}

func (c *rbacRulesRoleClause) DeletePermission(action Action, resource Resource) error {
	return errors.Errorf("not implements")
}

type rbacUserClause struct {
	*manager
	user *User
	// SetRole(role Role) error
}

func (c *rbacUserClause) AddRole(role Role) error {
	if c.user.Attr == nil {
		c.user.Attr = map[string]string{}
	}
	c.user.Attr[AttrKeyRolePrefix+role] = "true"

	if err := c.manager.store.PutUser(c.user); err != nil {
		return err
	}

	return nil
}
func (c *rbacUserClause) RemoveRole(role Role) error {

	if c.user.Attr == nil {
		c.user.Attr = map[string]string{}
	}

	delete(c.user.Attr, AttrKeyRolePrefix+role)

	if err := c.manager.store.PutUser(c.user); err != nil {
		return err
	}

	return nil
}
func (c *rbacUserClause) Roles() ([]Role, error) {
	result := []Role{}
	for k, _ := range c.user.Attr {
		if strings.HasPrefix(k, AttrKeyRolePrefix) {
			result = append(result, strings.TrimPrefix(k, AttrKeyRolePrefix))
		}
	}
	return result, nil
}
