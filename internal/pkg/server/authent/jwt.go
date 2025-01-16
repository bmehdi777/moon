package authent

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"moon/internal/pkg/server/config"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyJwt(tokenStr string) (*jwt.Token, error) {
	spkiPem := "-----BEGIN PUBLIC KEY-----\n" + config.GlobalConfig.RealmConfig.PublicKey + "\n-----END PUBLIC KEY-----"

	spkiBlock, _ := pem.Decode([]byte(spkiPem))
	var spkiKey *rsa.PublicKey
	pubInterface, _ := x509.ParsePKIXPublicKey(spkiBlock.Bytes)
	spkiKey = pubInterface.(*rsa.PublicKey)

	token, err := jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tok.Header["alg"])
		}

		return spkiKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
