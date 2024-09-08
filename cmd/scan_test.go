package cmd

import (
	"errors"
	"testing"
)

func TestParseSizeWithDecimalUnits(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{"1KB", 1000},
		{"2MB", 2000000},
		{"3GB", 3000000000},
	}

	for _, tt := range tests {
		result, err := parseSize(tt.input)
		if err != nil {
			t.Errorf("unexpected error for input %s: %v", tt.input, err)
		}
		if result != tt.expected {
			t.Errorf("expected %d, got %d for input %s", tt.expected, result, tt.input)
		}
	}
}

func TestParseSizeWithBinaryUnits(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{"1KiB", 1024},
		{"2MiB", 2097152},
		{"3GiB", 3221225472},
	}

	for _, tt := range tests {
		result, err := parseSize(tt.input)
		if err != nil {
			t.Errorf("unexpected error for input %s: %v", tt.input, err)
		}
		if result != tt.expected {
			t.Errorf("expected %d, got %d for input %s", tt.expected, result, tt.input)
		}
	}
}

func TestParseSizeWithNoUnit(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{"1", 1},
		{"100", 100},
		{"999", 999},
	}

	for _, tt := range tests {
		result, err := parseSize(tt.input)
		if err != nil {
			t.Errorf("unexpected error for input %s: %v", tt.input, err)
		}
		if result != tt.expected {
			t.Errorf("expected %d, got %d for input %s", tt.expected, result, tt.input)
		}
	}
}

func TestParseSizeWithMixedCaseUnits(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{"1kb", 1000},
		{"2Mb", 2000000},
		{"3gB", 3000000000},
	}

	for _, tt := range tests {
		result, err := parseSize(tt.input)
		if err != nil {
			t.Errorf("unexpected error for input %s: %v", tt.input, err)
		}
		if result != tt.expected {
			t.Errorf("expected %d, got %d for input %s", tt.expected, result, tt.input)
		}
	}
}

func TestParseSizeWithSpaces(t *testing.T) {
	tests := []struct {
		input    string
		expected uint64
	}{
		{" 123  ", 123},
		{" 1KB ", 1000},
		{"1 MB", 1000000},
		{" 3 GB", 3000000000},
	}

	for _, tt := range tests {
		result, err := parseSize(tt.input)
		if err != nil {
			t.Errorf("unexpected error for input %s: %v", tt.input, err)
		}
		if result != tt.expected {
			t.Errorf("expected %d, got %d for input %s", tt.expected, result, tt.input)
		}
	}
}

func TestParseSizeWithMissingNumericValue(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"KB"},
		{"MB"},
		{"GB"},
		{"KIB"},
		{"MIB"},
		{"GIB"},
	}

	for _, tt := range tests {
		_, err := parseSize(tt.input)
		if !errors.Is(err, ErrMissingNumericValue) {
			t.Errorf("expected ErrMissingNumericValue for input %s, but got %v", tt.input, err)
		}
	}
}

func TestParseSizeWithInvalidNumericPart(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"2  5KB"},
		{"2,5KB"},
		{"2.5KB"},
		{"3_000MB"},
		{"3,000GB"},
		{"3,0,00GB"},
	}

	for _, tt := range tests {
		_, err := parseSize(tt.input)
		if !errors.Is(err, ErrInvalidNumericValue) {
			t.Errorf("expected ErrInvalidNumericValue for input %s, but got %v", tt.input, err)
		}
	}
}

func TestParseSizeWithSpecialCharacters(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"1@KB"},
	}

	for _, tt := range tests {
		_, err := parseSize(tt.input)
		if !errors.Is(err, ErrUnknownSizeUnit) {
			t.Errorf("expected ErrUnknownSizeUnit for input %s, but got %v", tt.input, err)
		}
	}
}

func TestParseSizeWithMixedCharacters(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"1K2B"},
	}

	for _, tt := range tests {
		_, err := parseSize(tt.input)
		if !errors.Is(err, ErrUnknownSizeUnit) {
			t.Errorf("expected ErrUnknownSizeUnit for input %s, but got %v", tt.input, err)
		}
	}
}

func TestParseSizeWithInvalidFormat(t *testing.T) {
	invalidInputs := []string{"1XB", "10ZB", "1abc"}

	for _, input := range invalidInputs {
		_, err := parseSize(input)
		if !errors.Is(err, ErrUnknownSizeUnit) {
			t.Errorf("expected ErrUnknownSizeUnit for input %s, but got %v", input, err)
		}
	}
}

func TestParseSizeWithEmptyStringInput(t *testing.T) {
	_, err := parseSize("")
	if !errors.Is(err, ErrMissingArgumentValue) {
		t.Errorf("expected ErrMissingArgumentValue for empty string input, but got %v", err)
	}
}
