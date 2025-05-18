package stringCheckers

import "errors"

func CheckToken(token string) error {

	if len(token) > 100 || len(token) < 32 {
		return errors.New("bad token") // zwracamy tylko tokeny poniżej 100 i powyżej 32 znaków więc to oznacza że ktoś majstruje przy api
	}

	return nil
}
