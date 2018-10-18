package codes

type ErrorCode int32

const (
	None ErrorCode = iota
	EmptyAccount
	WrongEncoding
	EmptyHeader
	BadToken
	NotImplement
	Store
	Unauthorized
	Unknown
)
