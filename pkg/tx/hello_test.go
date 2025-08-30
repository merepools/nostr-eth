package tx

import (
	"testing"
)

func TestHelloWorld(t *testing.T) {
	expected := "Hello, World!"
	result := HelloWorld()

	if result != expected {
		t.Errorf("HelloWorld() = %q, want %q", result, expected)
	}
}

func TestHelloWorldWithName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Alice", "Alice", "Hello, Alice!"},
		{"Bob", "Bob", "Hello, Bob!"},
		{"Empty", "", "Hello, !"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HelloWorldWithName(tt.input)
			if result != tt.expected {
				t.Errorf("HelloWorldWithName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
