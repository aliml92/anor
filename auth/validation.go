package auth

import (
	"errors"
	"fmt"
	"regexp"
)

func validatePassword(password string) error {
	// Check for length between 8 and 20 characters
	if len(password) < 8 || len(password) > 20 {
		return fmt.Errorf("password length should be between 8 and 20 characters")
	}

	// Check for at least one lowercase letter
	if ok, _ := regexp.MatchString("[a-z]", password); !ok {
		return fmt.Errorf("password should contain at least one lowercase letter")
	}

	// Check for at least one uppercase letter
	if ok, _ := regexp.MatchString("[A-Z]", password); !ok {
		return fmt.Errorf("password should contain at least one uppercase letter")
	}

	// Check for at least one digit
	if ok, _ := regexp.MatchString("\\d", password); !ok {
		return fmt.Errorf("password should contain at least one digit")
	}

	// Check for at least one special character from @$!%*?&
	if ok, _ := regexp.MatchString("[@$!%*?&]", password); !ok {
		return errors.New("password should contain at least one special character from @$!%*?&")
	}

	return nil
}
