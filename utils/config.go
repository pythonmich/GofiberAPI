package utils

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	BuildVersion             string        `mapstructure:"BUILD_VERSION"`
	ServerAddress            string        `mapstructure:"SERVER_ADDRESS"`
	DBName                   string        `mapstructure:"DB_NAME"`
	DBDriver                 string        `mapstructure:"DB_DRIVER"`
	DBTimeout                time.Duration `mapstructure:"DB_TIMEOUT"`
	TokenSymmetricKey        string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenDuration            time.Duration `mapstructure:"TOKEN_DURATION"`
	RefreshTokenDuration     time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RefreshTokenSymmetricKey string        `mapstructure:"REFRESH_TOKEN_SYMMETRIC_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
