package stringCheckers

import (
	"errors"
	"strings"
)

func CheckCompanyName(name string) error {
	if len(name) > 40 {
		return errors.New("za długa nazwa")
	}
	if len(name) < 3 {
		return errors.New("nazwa za krótka")
	}
	if strings.ContainsAny(name, "[!@#$%^&*()+={}[]|\\:;\"'<>,.?/~]") {
		return errors.New("zawiera niepoprawne znaki")
	}

	return nil
}
