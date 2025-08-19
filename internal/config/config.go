package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   Server   `mapstructure:"server"`
	Database Database `mapstructure:"database"`
}

type Server struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type Database struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	SSLMode  string        `mapstructure:"sslmode"`
	User     string        `mapstructure:"user"`
	Password string        `mapstructure:"password"`
	Name     string        `mapstructure:"name"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

func LoadConfig(cfgPath string) (*Config, error) {
	setDefaultValues()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(cfgPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}
	log.Printf("Loaded configuration: %s", cfgPath)

	return &cfg, nil
}

func setDefaultValues() {
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 5*time.Second)
	viper.SetDefault("server.write_timeout", 5*time.Second)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.user", "user")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.name", "dbname")
	viper.SetDefault("database.timeout", 60*time.Second)
}
