package token

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var PrivateKey *rsa.PrivateKey
var PublicKey *rsa.PublicKey

func Init() error {
	parsedPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("RSA_PRIVATE_KEY")))
	if err != nil {
		return err
	}
	PrivateKey = parsedPrivateKey

	parsedPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("RSA_PUBLIC_KEY")))
	if err != nil {
		return err
	}
	PublicKey = parsedPublicKey

	return nil
}

func Sign(claims *jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(PrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Verify(tokenString string, claims *jwt.MapClaims) (*jwt.MapClaims, bool) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return PublicKey, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		return claims, token.Valid
	}
	
	return nil, token.Valid
}
