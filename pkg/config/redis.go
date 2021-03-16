package config

import "github.com/spf13/viper"

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func InitRedisConfig() (*RedisConfig, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("redis")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg RedisConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
