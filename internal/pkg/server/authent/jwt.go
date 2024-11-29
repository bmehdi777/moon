package authent

type AccessTokenDecoded struct {
	exp            int
	iat            int
	auth_time      int
	jti            string
	iss            string
	sub            string
	typ            string
	azp            string
	session_state  string
	scope          string
	sid            string
	email_verified bool
	email          string
}
