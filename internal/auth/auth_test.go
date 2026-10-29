package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	t.Run("Valid Header", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("Authorization", "Bearer my-secret-token-123")

		token, err := GetBearerToken(headers)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "my-secret-token-123"
		if token != expected {
			t.Errorf("expected token %q, got %q", expected, token)
		}
	})

	t.Run("Missing Header", func(t *testing.T) {
		headers := http.Header{}

		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("expected error for missing header, got nil")
		}
	})

	t.Run("Malformed Header - No Bearer Prefix", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("Authorization", "my-secret-token-123")

		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("expected error for malformed header, got nil")
		}
	})

	t.Run("Malformed Header - Wrong Prefix", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("Authorization", "Basic my-secret-token-123")

		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("expected error for wrong prefix, got nil")
		}
	})
}
