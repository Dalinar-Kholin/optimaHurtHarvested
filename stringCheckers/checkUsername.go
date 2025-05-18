package stringCheckers

import (
	"errors"
	"regexp"
)

func CheckUsername(username string) error {
	if len(username) < 6 {
		return errors.New("za krókta nazwa")
	}
	if len(username) > 20 {
		return errors.New("za długa nazwa")
	}

	allowedPattern := `^[a-zA-Z0-9_.-@]+$`
	if matched, err := regexp.MatchString(allowedPattern, username); !matched || err != nil {
		// W razie błędu zwracamy false
		return errors.New("zawiera niedozwolone znaki, jedyne dozwolone znaki to a-z A-Z 0-9 oraz _, ., -")
	}

	return nil
}
