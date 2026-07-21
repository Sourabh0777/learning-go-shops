package config
import "time"

type serverConfig struct{
	Port string
	GinMode string
}

type DatabaseConfig struct {
	Host string
	Port string
	User string
	Password string
	Name string
	SSLMode string
}

type JWTConfig struct {
	Secret string
	ExpiresIn time.
}