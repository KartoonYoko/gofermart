package monetary

func GetCurencyFromFloat64(f float64) int {
	return int(f * 100)
}

func GetFloat64FromCurrency(currency int) float64 {
	return float64(currency/100 + currency%100)
}
