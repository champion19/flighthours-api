package idencoder

import (
	"testing"

	"github.com/google/uuid"
)

func TestIDEncoder_EncodeDecodeUUID(t *testing.T) {
	encoder, err := NewHashidsEncoder(Config{
		Secret:    "test-secret-12345",
		MinLength: 10,
	}, nil)
	if err != nil {
		t.Fatalf("Error creating encoder: %v", err)
	}

	testUUID := uuid.New()
	testUUIDStr := testUUID.String()

	encoded, err := encoder.Encode(testUUIDStr)
	if err != nil {
		t.Fatalf("Error encoding UUID: %v", err)
	}

	if len(encoded) < 10 {
		t.Errorf("Encoded string is too short: got %d, want at least 10", len(encoded))
	}

	decoded, err := encoder.Decode(encoded)
	if err != nil {
		t.Fatalf("Error decoding UUID: %v", err)
	}

	if decoded != testUUIDStr {
		t.Errorf("Decoded UUID doesn't match original: got %v, want %v", decoded, testUUIDStr)
	}

	t.Logf("Original UUID: %v", testUUID)
	t.Logf("UUID String: %s", testUUIDStr)
	t.Logf("Encoded: %s", encoded)
	t.Logf("Decoded UUID: %s", decoded)
}
