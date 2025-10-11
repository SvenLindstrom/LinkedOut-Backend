package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Tokens struct {
	Access  string
	Refresh string
}

const (
	Refresh string = "refresh"
	Access  string = "access"
)

var keys = map[string][]byte{
	Refresh: []byte(os.Getenv("TOKEN_REFRESH_KEY")),
	Access:  []byte(os.Getenv("TOKEN_ACCESS_KEY")),
}

func CreatTokenPair(id string) (Tokens, error) {

	newAccessToken, errAcc := newAccessToken(id)
	newRefreshToken, errRef := newRefreshToken(id)

	if errAcc != nil || errRef != nil {
		return Tokens{}, errors.New("failed creat")
	}
	return Tokens{newAccessToken, newRefreshToken}, nil
}

func newClaims(audience string, subject string) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "linkedOut",
		Subject:   subject,
		ID:        uuid.NewString(),
		Audience:  []string{audience},
	}
}

func newRefreshToken(id string) (string, error) {
	claim := newClaims(Refresh, id)
	token, err := Sign(claim)
	return token, err
}

func newAccessToken(id string) (string, error) {
	claim := newClaims(Access, id)
	token, err := Sign(claim)
	return token, err
}

func Verify(tokeString string, audience string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokeString,
		&jwt.RegisteredClaims{},
		keyFunc(audience),
		jwt.WithIssuedAt(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithAudience(audience),
		jwt.WithIssuer("linkedOut"),
	)

	if err != nil {
		return &jwt.RegisteredClaims{}, err

	}

	claims := token.Claims.(*jwt.RegisteredClaims)
	return claims, err
}

func keyFunc(audience string) func(*jwt.Token) (any, error) {
	key := keys[audience]
	return func(token *jwt.Token) (any, error) {
		return key, nil
	}
}

func Sign(claims jwt.Claims) (string, error) {
	aud, err := claims.GetAudience()
	key := keys[aud[0]]
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed_token, err := token.SignedString(key)
	return signed_token, err
}
