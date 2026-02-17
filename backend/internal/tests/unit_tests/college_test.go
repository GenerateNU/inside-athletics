package unitTests

import (
	h "inside-athletics/internal/handlers/college"
	"testing"
)

func TestStringPtrOrNil(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string returns nil",
			input: "",
		},
		{
			name:  "non-empty string returns pointer to string",
			input: "https://www.northeastern.edu",
		},
		{
			name:  "string with spaces returns pointer",
			input: "  ",
		},
		{
			name:  "single character returns pointer",
			input: "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := h.StringPtrOrNil(tt.input)

			if tt.input == "" {
				if result != nil {
					t.Fatalf("expected nil for empty string, got %v", result)
				}
			} else {
				if result == nil {
					t.Fatalf("expected pointer to %s, got nil", tt.input)
					return
				}
				if *result != tt.input {
					t.Fatalf("expected %s, got %s", tt.input, *result)
				}
			}
		})
	}
}
