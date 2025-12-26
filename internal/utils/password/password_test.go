package password

import "testing"

func TestHashVerify(t *testing.T) {
	h, err := Hash("password123")
	if err != nil {
		t.Fatalf("hash err: %v", err)
	}
	if !Verify(h, "password123") {
		t.Fatalf("verify should pass")
	}
	if Verify(h, "wrong") {
		t.Fatalf("verify should fail")
	}
}
