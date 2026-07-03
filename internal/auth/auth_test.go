package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "correct horse battery staple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == password {
		t.Fatal("hash should not match the raw password")
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}
	if !match {
		t.Fatal("expected password to match hash")
	}

	match, err = CheckPasswordHash("wrong password", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}
	if match {
		t.Fatal("expected wrong password not to match hash")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "super-secret"

	tokenString, err := MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	gotUserID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	if gotUserID != userID {
		t.Fatalf("expected user ID %s, got %s", userID, gotUserID)
	}
}

func TestValidateJWTRejectsExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "super-secret"

	tokenString, err := MakeJWT(userID, tokenSecret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	if _, err := ValidateJWT(tokenString, tokenSecret); err == nil {
		t.Fatal("expected expired token to be rejected")
	}
}

func TestValidateJWTRejectsWrongSecret(t *testing.T) {
	userID := uuid.New()

	tokenString, err := MakeJWT(userID, "super-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	if _, err := ValidateJWT(tokenString, "wrong-secret"); err == nil {
		t.Fatal("expected token signed with wrong secret to be rejected")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer  abc123  ")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned error: %v", err)
	}

	if token != "abc123" {
		t.Fatalf("expected token abc123, got %q", token)
	}
}

func TestGetBearerTokenRejectsMissingHeader(t *testing.T) {
	if _, err := GetBearerToken(http.Header{}); err == nil {
		t.Fatal("expected missing authorization header to error")
	}
}

func TestGetBearerTokenRejectsMalformedHeader(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Token abc123")

	if _, err := GetBearerToken(headers); err == nil {
		t.Fatal("expected malformed authorization header to error")
	}
}
