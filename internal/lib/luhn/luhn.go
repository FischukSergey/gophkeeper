package luhn

func Valid(cardNumber string) bool {
	sum := 0
	isSecond := false
	
	// Идем с конца строки
	for i := len(cardNumber) - 1; i >= 0; i-- {
		// Пропускаем пробелы
		if cardNumber[i] == ' ' {
			continue
		}
		
		// Проверяем, что символ является цифрой
		if cardNumber[i] < '0' || cardNumber[i] > '9' {
			return false
		}

		digit := int(cardNumber[i] - '0')

		if isSecond {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isSecond = !isSecond
	}

	return sum%10 == 0 && len(cardNumber) > 1
}
