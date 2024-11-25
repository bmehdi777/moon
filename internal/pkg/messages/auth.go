package messages

type AuthRequest struct {
	Version byte
	Email string
	AccessToken string
}

