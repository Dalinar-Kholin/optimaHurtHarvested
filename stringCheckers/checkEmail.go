package stringCheckers

import (
	"errors"
	"strings"
)

func CheckEmail(email string) error {
	if len(email) > 40 {
		return errors.New("email to long")
	}
	if len(email) < 10 {
		return errors.New("email to short")
	}

	if strings.Count(email, "@") != 1 {
		return errors.New("bac @ count")
	}
	if !strings.Contains(email, "@") || strings.ContainsAny(email, "(),:;<>\"[\\] ") {
		return errors.New("bad email syntax")
	}

	return nil
}
