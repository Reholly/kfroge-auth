package entity

const (
	UsernameClaim   = "username"
	ExpirationClaim = "exp"
	IdClaim         = "id"
	RoleClaim       = "role"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Claim struct {
	Title string
	Value any
}
