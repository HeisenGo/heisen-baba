package wallet

import "strconv"

func isValidCardNumber(cardNumber string) bool {
	if len(cardNumber) != 16 {
		return false
	}

	sum := 0
	for i := 0; i < len(cardNumber); i++ {
		digit, err := strconv.Atoi(string(cardNumber[i]))
		if err != nil {
			return false
		}

		if i%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
