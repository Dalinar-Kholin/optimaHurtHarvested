package stringCheckers

import (
	"errors"
	"regexp"
)

func CheckEan(ean string) error {
	if len(ean) < 5 || len(ean) > 15 {
		return errors.New("błędna długość ean")
	}
	allowedPattern := `^[0-9]+`
	if matched, err := regexp.MatchString(allowedPattern, ean); !matched || err != nil {
		// W razie błędu zwracamy false
		return errors.New("zawiera niedozwolone znaki")
	}

	return nil
}
