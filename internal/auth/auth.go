package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// HashPassword
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// CheckPasswordHash
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.Issuer != "chirpy-access" {
			return uuid.Nil, jwt.ErrTokenInvalidClaims
		}
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, err
		}
		return userID, nil
	}

	return uuid.Nil, jwt.ErrTokenInvalidClaims
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok {
		return "", errors.New("malformed authorization header")
	}

	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return "", errors.New("missing bearer token")
	}

	return tokenString, nil
}

func MakeRefreshToken() string {
	token := make([]byte, 32)
	_, _ = rand.Read(token)
	return hex.EncodeToString(token)
}


// GetAPIKey 

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("missing authorization header")
	}

	apiKey, ok := strings.CutPrefix(apiKey, "ApiKey ")
	if !ok {
		return "", errors.New("malformed authorization header")
	}

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return "", errors.New("missing api key")
	}

	return apiKey, nil
}
