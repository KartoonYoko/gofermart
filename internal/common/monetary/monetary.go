package monetary

func GetCurencyFromFloat64(f float64) int {
	return int(f * 100)
}

func GetFloat64FromCurrency(currency int) float64 {
	intNum := currency/100
	remainder := currency%100
	return float64(intNum) + float64(remainder) / 100
}
