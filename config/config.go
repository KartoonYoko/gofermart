package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress           string // адрес и порт запуска сервиса
	DatabaseURL          string // адрес подключения к базе данных
	AccrualSystemAddress string // адрес системы расчёта начислений
}

func New() *Config {
	runAddressFlag := flag.String("a", ":8000", "Server address")
	databaseURLConnectionFlag := flag.String("d", "", "Database string connection")
	accrualSystemAddressFlag := flag.String("r", "", "Accrual system address")
	flag.Parse()

	c := &Config{
		RunAddress:           *runAddressFlag,
		DatabaseURL:          *databaseURLConnectionFlag,
		AccrualSystemAddress: *accrualSystemAddressFlag,
	}

	if envServerAddr := os.Getenv("RUN_ADDRESS"); envServerAddr != "" {
		c.RunAddress = envServerAddr
	}

	if emvDbsc := os.Getenv("DATABASE_URI"); emvDbsc != "" {
		c.DatabaseURL = emvDbsc
	}

	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		c.AccrualSystemAddress = envAccrualSystemAddress
	}

	return c
}
