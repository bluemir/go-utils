package auth

type isTokenClause struct {
	*manager
	t *Token
}

func (c *isTokenClause) Allow(action Action) bool {
	if c.t == nil {
		return false
	}

	user, ok, err := c.store.GetUser(c.t.Username)
	if err != nil {
		return false
	}
	if !ok {
		// user not exist
		return false
	}
	return c.checkPerm(action) && c.Is(user).Allow(action)
}
func (c *isTokenClause) NotAllow(action Action) bool {
	return !c.Allow(action)
}
func (c *isTokenClause) checkPerm(action Action) bool {
	// permission is empty it means same as user
	if c.t.Allows == nil || len(c.t.Allows) == 0 {
		return true
	}
	for _, perm := range c.t.Allows {
		if action == perm {
			return true
		}
	}
	return false
}

type isUserClause struct {
	*manager
	u *User
}

func (c *isUserClause) Allow(action Action) bool {
	// find role and check action
	return c.Is(c.u.Role).Allow(action)
	//return c.manager.store.HasRule(c.u.Role, action)
}
func (c *isUserClause) NotAllow(action Action) bool {
	return !c.Allow(action)
}

type isRoleClause struct {
	*manager
	r Role
}

func (c *isRoleClause) Allow(action Action) bool {
	if Role(c.manager.opt.RootRole) == c.r {
		return true
	}
	return c.manager.store.HasRule(c.r, action)
}
func (c *isRoleClause) NotAllow(action Action) bool {
	return !c.Allow(action)
}

type isDefaultClause struct{}

func (c *isDefaultClause) Allow(action Action) bool {
	return false
}
func (c *isDefaultClause) NotAllow(action Action) bool {
	return !c.Allow(action)
}
