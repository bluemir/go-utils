package auth

const (
	ActionCreate Action = "Create"
	ActionGet           = "GET"
	ActionList          = "LIST"
	ActionUpdate        = "UPDATE"
	ActionPatch         = "PATCH"
	ActionDelete        = "DELETE"
)

const (
	ResultAllow AuthzResult = true
	ResultDeny  AuthzResult = false
)

const (
	AttrKeyKind       = "kind"
	AttrKeyRolePrefix = "role/"
)
