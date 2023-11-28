package vld

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidMembername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidNameEntire = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("should have between %d and %d characters", minLength, maxLength)
	}
	return nil
}

func ValidateMembername(value string) error {
	if err := ValidateString(value, 2, 80); err != nil {
		return err
	}
	if !isValidMembername(value) {
		return fmt.Errorf("should only include lowercase letters, numbers, or underscores")
	}
	return nil
}

func ValidateNameEntire(value string) error {
	if err := ValidateString(value, 2, 80); err != nil {
		return err
	}
	if !isValidNameEntire(value) {
		return fmt.Errorf("should only consist of letters and spaces")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 32)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 80); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is an invalid email address")
	}
	return nil
}
