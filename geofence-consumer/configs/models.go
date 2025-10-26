package configs

import "time"

type Config struct {
	App      AppConfig      `mapstructure:"APP"`
	Database Database       `mapstructure:"DATABASE"`
	RabbitMQ RabbitMQConfig `mapstructure:"RABBIT_MQ"`
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

type AppConfig struct {
	Name           string        `mapstructure:"NAME" validate:"required"`
	ContextTimeout time.Duration `mapstructure:"CONTEXT_TIMEOUT" validate:"required"`
	Env            string        `mapstructure:"ENV" validate:"required"`
	LogLevel       string        `mapstructure:"LOG_LEVEL" validate:"omitempty,required"`
}

type RabbitMQConfig struct {
	Host     string                 `mapstructure:"HOST" validate:"required"`
	Port     string                 `mapstructure:"PORT" validate:"required"`
	User     string                 `mapstructure:"USER"`
	Password string                 `mapstructure:"PASSWORD"`
	Vhost    string                 `mapstructure:"VHOST"`
	Consumer RabbitMqConsumerConfig `mapstructure:"CONSUMER"`
}

type RabbitMqConsumerConfig struct {
	GeoFence struct {
		Enabled     bool   `mapstructure:"ENABLED" validate:"required"`
		Exchange    string `mapstructure:"EXCHANGE" validate:"required"`
		Queue       string `mapstructure:"QUEUE" validate:"required"`
		ConsumerTag string `mapstructure:"CONSUMER_TAG" validate:"required"`
		RoutingKey  string `mapstructure:"ROUTING_KEY" validate:"required"`
	} `mapstructure:"GEO_FENCE"`
}
