package messages

type AuthRequest struct {
	Version byte
	AccessTokenJWT string
}

