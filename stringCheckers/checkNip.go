package stringCheckers

import (
	"errors"
	"strconv"
	"strings"
)

func CheckNip(nip string) error {
	// Usuń wszelkie znaki inne niż cyfry
	nip = strings.ReplaceAll(nip, "-", "")

	// NIP powinien mieć dokładnie 10 cyfr
	if len(nip) != 10 {
		return errors.New("zła długość nip")
	}

	// Przekształć ciąg znaków na tablicę liczb całkowitych
	weights := []int{6, 5, 7, 2, 3, 4, 5, 6, 7}
	var sum int

	for i := 0; i < 9; i++ {
		digit, err := strconv.Atoi(string(nip[i]))
		if err != nil {
			return errors.New("niepoprawne znaki w nip")
		}
		sum += digit * weights[i]
	}

	controlDigit := sum % 11

	// Sprawdź, czy cyfra kontrolna zgadza się z ostatnią cyfrą NIP
	lastDigit, err := strconv.Atoi(string(nip[9]))
	if err != nil || controlDigit != lastDigit {
		return errors.New("niepoprawny nip")
	}

	return nil
}
