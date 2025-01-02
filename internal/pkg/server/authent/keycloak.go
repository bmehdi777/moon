package authent

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyJwt(tokenStr string) (*jwt.Token, error) {
	// local cert for testing purpose
	spkiPem := `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0GGJjxtCXGQgKxwcFZwd2AkNdaPSMN76A5bJyk6Dve8gMi8sbypzKngzhkziqofVe9g5H9kWRyZNIVzKiK4OnFhTRvRtAXoeWj98EINRMmvmWGv5BKwGmfr7g/mVvr+viyROUrrPUWx6TslyVD7VxLFrSchLiAdV6pZdMrKD1tlSXNQ78N3Q2Nw/SmuYd07wBIbtDCTwG9XaCJFaw0jgbKs6wdpTSqkfTNnYE2ekOlI8nAtTwAthjJeIfuPuScG4wVvbTTMx+Hd3z4kU2ripynSOVOWioyWUw6uerJqt1sgclNdQkFwdXgCzcOmJYIt8cOvCm8jEkNPmL3jJMN/eVQIDAQAB
-----END PUBLIC KEY-----
	`

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
