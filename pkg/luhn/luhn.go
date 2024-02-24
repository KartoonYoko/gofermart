// Реализация упрощенного алгоритма Луна описанного тут
// https://ru.wikipedia.org/wiki/%D0%90%D0%BB%D0%B3%D0%BE%D1%80%D0%B8%D1%82%D0%BC_%D0%9B%D1%83%D0%BD%D0%B0

package luhn

import (
	"errors"
	"strconv"
)

// ErrNotValid сигнализирует, что переданные данные не проходят проверку алгоритма Луна
var ErrNotValid = errors.New("luhn: data not valid")

func CheckStr(str string) (int64, error) {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return -1, ErrNotValid
	}

	if !valid(num) {
		return -1, ErrNotValid
	}

	return num, nil
}

func valid(number int64) bool {
	return (number%10+check(number/10))%10 == 0
}

func check(number int64) int64 {
	var luhn int64

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
