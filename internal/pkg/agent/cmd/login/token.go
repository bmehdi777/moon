package login

import "time"

type KeycloakJWTS struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type TokenDisk struct {
	AccessToken        string `json:"at"`
	AccessTokenExpire  int64  `json:"at_exp"`
	RefreshToken       string `json:"rt"`
	RefreshTokenExpire int64  `json:"rt_exp"`
}

func (k *KeycloakJWTS) ToDisk() *TokenDisk {
	atExp := time.Now().Add(time.Second * time.Duration(k.ExpiresIn))
	rtExp := time.Now().Add(time.Second * time.Duration(k.RefreshExpiresIn))

	td := TokenDisk{
		AccessToken:        k.AccessToken,
		AccessTokenExpire:  atExp.Unix(),
		RefreshToken:       k.RefreshToken,
		RefreshTokenExpire: rtExp.Unix(),
	}

	return &td
}
