package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL      string `mapstructure:"DATABASE_URL"`
	DatabaseHost     string `mapstructure:"DATABASE_HOST"`
	DatabasePort     string `mapstructure:"DATABASE_PORT"`
	DatabaseUser     string `mapstructure:"DATABASE_USER"`
	DatabasePassword string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseName     string `mapstructure:"DATABASE_NAME"`
	DatabaseSllMode  string `mapstructure:"DATABASE_SSLMODE"`
	DatabaseMaxOpen  int    `mapstructure:"DATABASE_MAX_CONNECTIONS"`
	DatabaseMaxIdle  int    `mapstructure:"DATABASE_MAX_IDLE_CONNECTIONS"`
	ElasticsearchURL string `mapstructure:"ELASTICSEARCH_URL"`
	NatsURL          string `mapstructure:"NATS_URL"`
	SMTPHost         string `mapstructure:"SMTP_HOST"`
	SMTPPort         string `mapstructure:"SMTP_PORT"`
	SMTPUser         string `mapstructure:"SMTP_USER"`
	SMTPPassword     string `mapstructure:"SMTP_PASSWORD"`
	APIPort          string `mapstructure:"API_PORT"`
	GRPCPort         string `mapstructure:"GRPC_PORT"`
	DOMAIN           string `mapstructure:"DOMAIN"`
}

func (v *Config) GetConnectionString() string {
	if v.DatabaseURL == "" {
		// return connection string by postgres driver
		disabled := "disable"
		if v.DatabaseSllMode != "" && (v.DatabaseSllMode == "true" || v.DatabaseSllMode == "require") {
			disabled = "enable"
		}
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", v.DatabaseUser, v.DatabasePassword, v.DatabaseHost, v.DatabasePort, v.DatabaseName, disabled)
	}
	return v.DatabaseURL
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
