package internal

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var (
	ErrEmptyString       = errors.New("input cannot be empty")
	ErrExceedsMaxLength  = errors.New("input exceeds maximum allowed length")
	ErrInvalidCharacters = errors.New("input contains invalid characters")
	ErrConsecutiveSpaces = errors.New("input contains consecutive spaces")
)

// ValidateName checks if the given name (e.g., terminal name, city name) is valid.
// It returns an error if the name is invalid, or nil if it's valid.
func ValidateName(name string, maxLength int) error {
	// Check if the name is empty
	if strings.TrimSpace(name) == "" {
		return ErrEmptyString
	}

	// Check if the name exceeds the maximum length
	if len(name) > maxLength {
		return ErrExceedsMaxLength
	}

	// Check for invalid characters
	invalidChars := []rune{';', '\'', '"', '-', '=', '*', '/', '\\', '<', '>', '(', ')', '{', '}', '[', ']'}
	for _, char := range invalidChars {
		if strings.ContainsRune(name, char) {
			return ErrInvalidCharacters
		}
	}

	// Check for consecutive spaces
	if strings.Contains(name, "  ") {
		return ErrConsecutiveSpaces
	}

	// Check if the name contains only printable characters
	for _, char := range name {
		if !unicode.IsPrint(char) {
			return ErrInvalidCharacters
		}
	}

	// Additional check: name should start and end with a letter or number
	if !regexp.MustCompile(`^[a-zA-Z0-9].*[a-zA-Z0-9]$`).MatchString(name) {
		return ErrInvalidCharacters
	}

	return nil
}
