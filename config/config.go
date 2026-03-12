package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var cfg Config

const JWTSecret = "4Rtg8BPKwixXy2ktDPxoMMAhRzmo9mmuZjvKONGPZZQSaJWNLijxR42qRgq0iBb5"

type Config struct {
	Debug    bool     `json:"debug"`
	Log      Log      `mapstructure:"log"`
	Db       Db       `mapstructure:"db"`
	Server   Server   `mapstructure:"server"`
	Merchant Merchant `mapstructure:"merchant"`
}

type Server struct {
	Port           int `mapstructure:"port"`
	JWTExpireHours int `mapstructure:"jwt_expire_hours"`
}
type Merchant struct {
	DefaultPassword string `mapstructure:"default_password"`
	//WithdrawalFee   float64 `mapstructure:"withdrawal_fee"`
	FrozenHours int64 `mapstructure:"frozen_hours"`
}
type Log struct {
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	LogLevel   string `mapstructure:"log_level"`
}

type Db struct {
	Driver       string `mapstructure:"driver"`
	Name         string `mapstructure:"name"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Charset      string `mapstructure:"charset"`
	Ssl          string `mapstructure:"ssl"`
	Schema       string `mapstructure:"schema"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxLifeTime  int    `mapstructure:"max_life_time"`
	LogLevel     string `mapstructure:"log_level"`
}

func LoadConfig(debug bool) {
	v := viper.New()
	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config error failed to read the configuration file: %s", err))
	}

	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("unmarshal config err : %s", err.Error()))
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 9000
	}

	if cfg.Merchant.DefaultPassword == "" {
		cfg.Merchant.DefaultPassword = "123456"
	}
	cfg.Debug = debug

}

func GetCfg() Config {
	return cfg
}
