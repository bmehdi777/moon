package authent

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"moon/internal/pkg/server/config"
)

func VerifyJwt(tokenStr string) (*jwt.Token, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("token is empty")
	}
	spkiPem := "-----BEGIN PUBLIC KEY-----\n" + config.GlobalConfig.RealmConfig.PublicKey + "\n-----END PUBLIC KEY-----"

	spkiBlock, _ := pem.Decode([]byte(spkiPem))
	if spkiBlock == nil {
		return nil, fmt.Errorf("invalid JWT public key PEM")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(spkiBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse JWT public key: %w", err)
	}
	spkiKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("JWT public key is not RSA")
	}

	expectedAlgorithm := config.GlobalConfig.Auth.Algorithm
	if expectedAlgorithm == "" {
		expectedAlgorithm = jwt.SigningMethodRS256.Alg()
	}
	if !strings.HasPrefix(expectedAlgorithm, "RS") {
		return nil, fmt.Errorf("unsupported JWT signing algorithm: %s", expectedAlgorithm)
	}
	token, err := jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
		return spkiKey, nil
	}, jwt.WithValidMethods([]string{expectedAlgorithm}))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}
	if _, err := token.Claims.GetSubject(); err != nil {
		return nil, fmt.Errorf("token subject is invalid: %w", err)
	}
	return token, nil
}
