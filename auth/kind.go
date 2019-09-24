package auth

func Kind(kind string) Resource {
	return KV{AttrKeyKind: kind}
}
