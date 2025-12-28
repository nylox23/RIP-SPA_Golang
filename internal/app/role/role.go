package role

type Role int

const (
	Guest Role = iota
	User
	Admin
)
