package domain

import (
	"testing"
)

func TestAirline_ToLogger(t *testing.T) {
	t.Run("returns expected log fields", func(t *testing.T) {
		airline := &Airline{
			ID:          "airline-uuid-123",
			AirlineName: "Test Airlines",
			AirlineCode: "TST",
			Status:      "active",
		}

		result := airline.ToLogger()

		if len(result) != 4 {
			t.Fatalf("expected 4 log fields, got %d", len(result))
		}

		expected := []string{
			"id:airline-uuid-123",
			"name:Test Airlines",
			"code:TST",
			"status:active",
		}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("expected result[%d] = %q, got %q", i, exp, result[i])
			}
		}
	})

	t.Run("handles empty fields", func(t *testing.T) {
		airline := &Airline{}

		result := airline.ToLogger()

		if len(result) != 4 {
			t.Fatalf("expected 4 log fields, got %d", len(result))
		}

		expected := []string{
			"id:",
			"name:",
			"code:",
			"status:",
		}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("expected result[%d] = %q, got %q", i, exp, result[i])
			}
		}
	})
}

func TestAirline_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "active status returns true",
			status:   "active",
			expected: true,
		},
		{
			name:     "inactive status returns false",
			status:   "inactive",
			expected: false,
		},
		{
			name:     "empty status returns false",
			status:   "",
			expected: false,
		},
		{
			name:     "ACTIVE (uppercase) returns false",
			status:   "ACTIVE",
			expected: false,
		},
		{
			name:     "unknown status returns false",
			status:   "pending",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			airline := &Airline{Status: tt.status}
			result := airline.IsActive()
			if result != tt.expected {
				t.Errorf("expected IsActive() = %v, got %v", tt.expected, result)
			}
		})
	}
}
