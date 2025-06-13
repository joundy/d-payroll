package entity

type Login struct {
	Username string
	Password string
}

type AuthToken struct {
	Token string
}

type AuthTokenPayload struct {
	Id   uint
	Role UserRole
}
