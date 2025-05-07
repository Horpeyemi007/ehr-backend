package auth

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
	aud     string
	issue   string
}

func NewJWTAuthenticator(privatePath, publicPath, aud, iss string) *JWTAuthenticator {
	privBytes, err := os.ReadFile(privatePath)
	if err != nil {
		log.Fatal("Error reading private key as byte: ", err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		log.Fatal("Error parsing private key file from pem: ", err)
	}
	pubBytes, err := os.ReadFile(publicPath)
	if err != nil {
		log.Fatal("Error reading public key as byte: ", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		log.Fatal("Error parsing public key file from pem: ", err)
	}

	return &JWTAuthenticator{privKey, pubKey, iss, aud}
}

func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(j.privKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return j.pubKey, nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(j.aud),
		jwt.WithIssuer(j.aud),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
	)
}
