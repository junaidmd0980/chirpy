package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	token_secret := "super-secret-key"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, token_secret, expiresIn)
	if err != nil {
		t.Fatalf("error making JWT: %v", err)
	}

	parsedID, err := ValidateJWT(token, token_secret)
	if err != nil {
		t.Fatalf("error validating JWT: %v", err)
	}

	if parsedID != userID {
		t.Errorf("expected user ID %v got %v", userID, parsedID)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	token_secret := "super-secret-key"
	expiresIn := -time.Second

	token, err := MakeJWT(userID, token_secret, expiresIn)
	if err != nil {
		t.Fatalf("error making JWT: %v", err)
	}

	_, err = ValidateJWT(token, token_secret)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestWrongSecretJWt(t *testing.T) {
	userID := uuid.New()
	token_secret := "super-secret-key"
	wrong_secret := "wrong-secret-key"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, token_secret, expiresIn)
	if err != nil {
		t.Fatalf("error making JWT: %v", err)
	}

	_, err = ValidateJWT(token, wrong_secret)
	if err == nil {
		t.Error("expected error for wrong secret, got nil")
	}
}
