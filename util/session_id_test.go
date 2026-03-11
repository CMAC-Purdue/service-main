package util

import "testing"

func TestGenerateSessionIDSize(t *testing.T) {
	id, err := GenerateSessionIDSize(32)
	if err != nil {
		t.Fatalf("GenerateSessionIDSize returned error: %v", err)
	}

	if id == "" {
		t.Fatal("expected non-empty id")
	}

	if len(id) != 43 {
		t.Fatalf("expected id length 43 for 32 random bytes, got %d", len(id))
	}
}

func TestGenerateSessionIDSizeInvalid(t *testing.T) {
	if _, err := GenerateSessionIDSize(0); err == nil {
		t.Fatal("expected error for invalid size")
	}
}
