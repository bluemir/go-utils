package auth

type abacClause struct {
	*manager
}

func (m *manager) ABAC() ABACClause {
	return &abacClause{m}
}
func (c *abacClause) User(user *User) AttrHandler {
	return &abacUserClause{c.manager, user}
}
func (c *abacClause) Token(token *Token) AttrHandler {
	return &abacTokenClause{c.manager, token}
}

type abacUserClause struct {
	*manager
	user *User
}

func (c *abacUserClause) Get(key string) (string, error) {
	return c.user.Attr[key], nil
}
func (c *abacUserClause) Set(key, value string) error {
	// TODO
	c.user.Attr[key] = value

	return c.manager.User().Update(c.user)
}

type abacTokenClause struct {
	*manager
	token *Token
}

func (c *abacTokenClause) Get(key string) (string, error) {
	return c.token.Attr[key], nil
}
func (c *abacTokenClause) Set(key, value string) error {
	// TODO
	return nil
}
