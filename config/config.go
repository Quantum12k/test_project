package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

// Config contains all configurable vars for apps.
type Config struct {

	ServerPort string `mapstructure:"server_port"`
	DBName string     `mapstructure:"db_name"`
	DBPass string     `mapstructure:"db_pass"`
	DBUser string     `mapstructure:"db_user"`
	DBHost string     `mapstructure:"db_host"`
	DBPort string     `mapstructure:"db_port"`
}

func NewConfig() (Config, error) {

	configInstance := Config{}
	if err := viper.Unmarshal(&configInstance); err != nil {

		return Config{}, errors.New("can't load config structure")
	}

	return configInstance, nil
}

func init() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	_ = err

	viper.SetDefault("server_port", ":8000")
	viper.SetDefault("db_name", "test_db")
	viper.SetDefault("db_pass", "54678")
	viper.SetDefault("db_user", "postgres")
	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", "5434")
}