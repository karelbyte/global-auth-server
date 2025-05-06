package libs

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"maps"
)

var privateKey *rsa.PrivateKey

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	if privateKey != nil {
		return privateKey, nil
	}
	keyData, err := os.ReadFile("certificates/private.pem")
	if err != nil {
		return nil, fmt.Errorf("could not read private key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}
	privateKey = key
	return privateKey, nil
}

func GenerateJWT(payload map[string]any, duration time.Duration) (string, int64, error) {
	key, err := LoadPrivateKey()
	if err != nil {
		return "", 0, err
	}
	exp := time.Now().Add(duration).Unix()
	claims := jwt.MapClaims{}
	maps.Copy(claims, payload)
	claims["exp"] = exp

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(key)
	if err != nil {
		return "", 0, err
	}
	return signed, exp, nil
}
