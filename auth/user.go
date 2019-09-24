package auth

type userClause struct {
	*manager
}

func (m *manager) User() UserClause {
	return &userClause{m}
}

func (m *userClause) Create(u *User) (*User, error) {
	return u, m.store.CreateUser(u)
}
func (m *userClause) Get(name string) (*User, bool, error) {
	return m.store.GetUser(name)
}
func (m *userClause) List(filter ...string) ([]User, error) {
	return m.store.ListUser()
}
func (m *userClause) Update(u *User) error {
	return m.store.PutUser(u)
}
func (m *userClause) Delete(u *User) error {
	return m.store.DeleteUser(u.Name)
}
