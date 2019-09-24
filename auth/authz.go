package auth

type authzClause struct {
	*manager
	token *Token
}

func (m *manager) Authz(token *Token) AuthzClause {
	return &authzClause{m, token}
}

func (c *authzClause) Do(action Action, resource Resource) bool {
	return c.manager.rules.check(&RuleContext{
		Token:    c.token,
		Resource: resource,
		Action:   action,
	})
}
func (c *authzClause) Create(resource Resource) bool {
	return c.Do(ActionCreate, resource)
}
func (c *authzClause) Get(resource Resource) bool {
	return c.Do(ActionGet, resource)
}
func (c *authzClause) List(resource Resource) bool {
	return c.Do(ActionList, resource)
}
func (c *authzClause) Update(resource Resource) bool {
	return c.Do(ActionUpdate, resource)
}
func (c *authzClause) Patch(resource Resource) bool {
	return c.Do(ActionPatch, resource)
}
func (c *authzClause) Delete(resource Resource) bool {
	return c.Do(ActionDelete, resource)
}
