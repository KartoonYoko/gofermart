package config

type JWTConfig struct {
	SecretJWTKey string // секрет для подписания JWT токенов
}

func NewJWTConfig() (*JWTConfig, error) {
	// можно реализовать логику получения данных из хранилища секретов
	conf := &JWTConfig{
		SecretJWTKey: "supersecretkey",
	}

	return conf, nil
}
