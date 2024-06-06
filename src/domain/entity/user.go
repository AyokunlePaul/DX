package entity

type UserType int

const (
	User UserType = iota
	ClientManager
	UserAdmin
	SuperAdmin
)
