package auth

/*
User:Role   = n:n
User:Token  = 1:n
Role:Action = n:n
*/

type User struct {
	Name string
	Attr map[string]string
}

type Token struct {
	*User
	HashedKey string
	RevokeKey string
	Attr      map[string]string
}

type Action = string

type Resource interface {
	Set(key, value string) error
	Get(key string) (string, error)
	List() map[string]string
}

type AuthzResult = bool

type Role = string

// simple Resource implements
type KV map[string]string

func (kv KV) Get(key string) (string, error) {
	return kv[key], nil
}
func (kv KV) Set(key, value string) error {
	kv[key] = value
	return nil
}
func (kv KV) List() map[string]string {
	return kv
}
