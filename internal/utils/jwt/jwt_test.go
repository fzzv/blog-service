package jwtutil

import (
	"testing"
	"time"
)

func TestJWTSignParse(t *testing.T) {
	m := Manager{
		Secret: []byte("test-secret-very-long"),
		Issuer: "blog-service",
		TTL:    time.Minute,
	}

	token, err := m.Sign(123, "user")
	if err != nil {
		t.Fatalf("sign err: %v", err)
	}

	claims, err := m.Parse(token)
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}

	if claims.UserID != 123 {
		t.Fatalf("uid want 123 got %d", claims.UserID)
	}
	if claims.Role != "user" {
		t.Fatalf("role want user got %s", claims.Role)
	}
}
