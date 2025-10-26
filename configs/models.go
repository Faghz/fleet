package configs

import "time"

// Configuration models for the accounting application
//
// Note: For API responses requiring pagination, use the standardized Metadata struct
// from pkg/api/http/response package. All paginated endpoints should return data
// in the format: { "items": [...], "metadata": {...} }
//
// See METADATA_IMPLEMENTATION.md for detailed usage guidelines.

type Config struct {
	App      AppConfig      `mapstructure:"APP"`
	Function FunctionConfig `mapstructure:"FUNCTION"`
	Http     HttpConfig     `mapstructure:"HTTP"`
	Database Database       `mapstructure:"DATABASE"`
	MQTT     MQTTConfig     `mapstructure:"MQTT"`
}

type Database struct {
	PostgreSQL struct {
		Host            string        `mapstructure:"HOST" validate:"required"`
		Port            string        `mapstructure:"PORT" validate:"required"`
		User            string        `mapstructure:"USER" validate:"required"`
		Password        string        `mapstructure:"PASSWORD" validate:"required"`
		DbName          string        `mapstructure:"DB_NAME" validate:"required"`
		MaxConns        int32         `mapstructure:"MAX_CONNS" validate:"required"`
		MinConns        int32         `mapstructure:"MIN_CONNS" validate:"required"`
		MaxConnLifetime time.Duration `mapstructure:"MAX_CONN_LIFETIME" validate:"required"`
		MaxConnIdleTime time.Duration `mapstructure:"MAX_CONN_IDLE_TIME" validate:"required"`
	} `mapstructure:"POSTGRESQL"`
	Redis struct {
		Host     string `mapstructure:"HOST" validate:"required"`
		Port     string `mapstructure:"PORT" validate:"required"`
		User     string `mapstructure:"USER" validate:"omitempty,required"`
		Password string `mapstructure:"PASSWORD" validate:"omitempty,required"`
		DB       int    `mapstructure:"DB" validate:"omitempty,required"`
	} `mapstructure:"REDIS"`
}

type HttpConfig struct {
	Port             string `mapstructure:"PORT"`
	AllowedOrigins   string `mapstructure:"ALLOWED_ORIGINS" validate:"required"`
	AllowCredentials bool   `mapstructure:"ALLOW_CREDENTIALS"`
}

type AppConfig struct {
	Name           string        `mapstructure:"NAME" validate:"required"`
	ContextTimeout time.Duration `mapstructure:"CONTEXT_TIMEOUT" validate:"required"`
	Env            string        `mapstructure:"ENV" validate:"required"`
	LogLevel       string        `mapstructure:"LOG_LEVEL" validate:"omitempty,required"`
}

type FunctionUserSecretKey struct {
	Email           string `mapstructure:"EMAIL" validate:"required"`
	EmailSalt       string `mapstructure:"EMAIL_SALT" validate:"required"`
	EmailSaltLength int    `mapstructure:"EMAIL_SALT_LENGTH" validate:"required"`
}

type FunctionUser struct {
	SecretKey FunctionUserSecretKey `mapstructure:"SECRET_KEY"`
}

type FunctionAuthSecretKey struct {
	PasswordSalt string `mapstructure:"PASSWORD_SALT" validate:"required"`
	SessionID    string `mapstructure:"SESSION_ID" validate:"required"`
}

type FunctionAuthToken struct {
	SecretKey string        `mapstructure:"SECRET_KEY" validate:"required"`
	Expire    time.Duration `mapstructure:"EXPIRE" validate:"required"`
}

type FunctionAuth struct {
	SecretKey FunctionAuthSecretKey `mapstructure:"SECRET_KEY"`
	Token     FunctionAuthToken     `mapstructure:"TOKEN"`
	Session   FunctionSession       `mapstructure:"SESSION"`
}

type FunctionConfig struct {
	User FunctionUser `mapstructure:"USER"`
	Auth FunctionAuth `mapstructure:"AUTH"`
}

type FunctionSession struct {
	CacheExpireTime time.Duration `mapstructure:"CACHE_EXPIRE_TIME" validate:"required"`
}

type AuthConfig struct {
	SecretKey     string        `mapstructure:"SECRET_KEY" validate:"required"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

type MQTTConfig struct {
	Broker   string `mapstructure:"BROKER" validate:"required"`
	Port     string `mapstructure:"PORT" validate:"required"`
	ClientID string `mapstructure:"CLIENT_ID" validate:"required"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	QoS      byte   `mapstructure:"QOS"`
	APIKey   string `mapstructure:"API_KEY" validate:"required"`
}
