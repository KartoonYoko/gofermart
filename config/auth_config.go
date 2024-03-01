package config

type AuthConfig struct {
	Sault string // соль для создания хешей паролей
}

func NewAuthConfig(sault string) *AuthConfig {
	return &AuthConfig{
		Sault: sault,
	}
}