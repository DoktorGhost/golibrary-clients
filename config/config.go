package config

import (
	"fmt"
	"github.com/DoktorGhost/platform/storage/psg"
	"github.com/spf13/viper"
	"sync"
)

var (
	once   sync.Once
	Config config
)

type config struct {
	ProviderConfig ProviderConfig `mapstructure:"ProviderConfig"`
	DBConfig       DBConfig       `mapstructure:"DBConfig"`
}

type ProviderConfig struct {
	Provider_port string `mapstructure:"PROVIDER_PORT"`
	Http_port     string `mapstructure:"HTTP_PORT"`
}

type DBConfig struct {
	DB_host  string `mapstructure:"DB_HOST"`
	DB_port  string `mapstructure:"DB_PORT"`
	DB_name  string `mapstructure:"DB_NAME"`
	DB_login string `mapstructure:"DB_LOGIN"`
	DB_pass  string `mapstructure:"DB_PASS"`
}

func LoadConfig() config {
	once.Do(func() {
		// Декодируем значения в структуру Config
		viper.BindEnv("DBConfig.DB_HOST", "DB_HOST")
		viper.BindEnv("DBConfig.DB_PORT", "DB_PORT")
		viper.BindEnv("DBConfig.DB_NAME", "DB_NAME")
		viper.BindEnv("DBConfig.DB_LOGIN", "DB_LOGIN")
		viper.BindEnv("DBConfig.DB_PASS", "DB_PASS")

		viper.BindEnv("ProviderConfig.PROVIDER_PORT", "PROVIDER_PORT")
		viper.BindEnv("ProviderConfig.HTTP_PORT", "HTTP_PORT")

		if err := viper.Unmarshal(&Config); err != nil {
			panic(fmt.Errorf("ошибка декодирования конфигурации: %w", err))
		}
	})

	return Config
}

// Конвертация из config.DBConfig в psg.DBConfig
func ConvertToPsgDBConfig(conf DBConfig) psg.DBConfig {
	return psg.DBConfig{
		DbHost:  conf.DB_host,
		DbPort:  conf.DB_port,
		DbName:  conf.DB_name,
		DbLogin: conf.DB_login,
		DbPass:  conf.DB_pass,
	}
}
