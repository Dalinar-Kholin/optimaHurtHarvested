package stringCheckers

import "errors"

// hasło będzie i tak argonowane(moje pieniądze chlip chlip) więc nie ma niebezpieczeństawa złych znaków
func CheckPassword(password string) error {
	if len(password) < 4 {
		return errors.New("password to short")
	}
	if len(password) > 20 {
		return errors.New("password to long")
	}

	return nil
}
