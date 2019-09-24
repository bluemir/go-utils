package auth

// manager.User().Create(user *User) (*User, error)
// manager.User().Get() (*User, error)
// manager.User().List(byAttr ...ListOpt) ([]User, error)
// manager.User().Update(user *User) error
// manager.User().Delete(username string) error

// manager.Token().Issue(username, unhasedKey string) (*Token, error)
// manager.Token().List(username string) ([]Token, error)
// manager.Token().Revoke(revorkKey string) error

// manager.Authn().Default(id, password string) (*Token, error)
// manager.Authn().HTTP(header http.Header) (*Token, error)

// manager.Authz(token).Do(action Action, resource Resource) bool
// manager.Authz(token).Create(resource Resource) bool
// manager.Authz(token).Get(resource Resource) bool
// manager.Authz(token).List(resource Resource) bool
// manager.Authz(token).Update(resource Resource) bool
// manager.Authz(token).Patch(resource Resource) bool
// manager.Authz(token).Delete(resource Resource) bool

// manager.ABAC().User(token).GetAttr(key string)
// manager.ABAC().User(token).SetAttr(key, value string)
// manager.ABAC().Token(token).GetAttr(key string)
// manager.ABAC().Token(token).SetAttr(key, value string)

// manager.RBAC().Rules().Role("admin").AddPermission("GET", Kind("article"))
// manager.RBAC().Rules().Role("admin").AddPermission(action, resource)
// manager.RBAC().Rules().Role("admin").DeletePermission(action, resource)
// manager.RBAC().Rules().Load();
// manager.RBAC().User(token).AddRole(role Role)
// manager.RBAC().User(token).RemoveRole(role Role)
// manager.RBAC().User(token).Roles() ([]Role, error)
// manager.RBAC().Token(token).AddRole(role Role)
// manager.RBAC().Token(token).RemoveRole(role Role)

// manager.Rule().Add(rule string) error
// manager.Rule().AddCustom(id string, f func(user *User, token *Token, Resource resource, action Action) bool) error (for custom ABAC)
// manager.Rule().List() ([]Rule, error)
// manager.Rule().Delete(id uint) error

// user.Attr[Email]

// token.User.Attr[Email]

// err.Code()
// err.Cause

// Kind(str) Resource // for simple use of RBAC
