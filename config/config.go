package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Env           string `mapstructure:"ENVIRONMENT"`
	APIPort       string `mapstructure:"API_PORT"`
	DbPort        string `mapstructure:"DB_PORT"`
	DbUsername    string `mapstructure:"DB_USERNAME"`
	DbPassword    string `mapstructure:"DB_PASSWORD"`
	DbName        string `mapstructure:"DB_NAME"`
	DbSSLMode     string `mapstructure:"DB_SSL_MODE"`
	KafkaAddress  string `mapstructure:"KAFKA_ADDRESS"`
	KafkaTopic    string `mapstructure:"KAFKA_TOPIC"`
	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
}

func Init() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read in config file: %w", err))
	}

	config := Config{}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshall config file: %w", err))
	}

	return &config
}
