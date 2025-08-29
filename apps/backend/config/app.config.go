package config

import (
	"strings"

	"github.com/spf13/viper"
)

type (
	I18N struct {
	}
	MAIL struct {
		HOST     string `yaml:"HOST"`
		PORT     int    `yaml:"port"`
		USERNAME string `yaml:"username"`
		PASSWORD string `yaml:"password"`
	}

	SERVER struct {
		HOST_NAME string `yaml:"HOST_NAME" `
		HOST      string `yaml:"HOST"`
		PORT      string `yaml:"PORT"`
	}

	DB struct {
		DSN string `yaml:"DSN"`
	}

	CORS struct {
		ALLOW_ORIGIN      []string `yaml:"ALLOW_ORIGIN"`
		ALLOW_METHODS     []string `yaml:"ALLOW_METHODS"`
		ALLOW_HEADERS     []string `yaml:"ALLOW_HEADERS"`
		EXPOSE_HEADERS    []string `yaml:"EXPOSE_HEADERS"`
		ALLOW_CREDENTIALS bool     `yaml:"ALLOW_CREDENTIALS"`
		MAX_AGE           string   `yaml:"MAX_AGE"`
	}

	JWT struct {
		SECRET         string `yaml:"SECRET"`
		REFRESH_SECRET string `yaml:"REFRESH_SECRET"`
		EMAIL_SECRET   string `yaml:"EMAIL_SECRET"`
	}

	CARBIN struct {
		MODEL string `yaml:"MODEL"`
	}
	PPROF struct {
		ENABLE bool `yaml:"ENABLE"`
	}
	LabODT struct {
		Url   string `yaml:"URL"`
		Token string `yaml:"TOKEN"`
	}

	Config struct {
		SERVER SERVER `yaml:"SERVER"`
		DB     DB     `yaml:"DB"`
		CORS   CORS   `yaml:"CORS"`
		JWT    JWT    `yaml:"JWT"`
		PPROF  PPROF  `yaml:"PPROF"`
		MAIL   MAIL   `yaml:"mail"`
		CARBIN CARBIN `yaml:"CARBIN"`
		LabODT LabODT `yaml:"LabODT"`
	}
)

func newAppConfig() *Config {
	conf := &Config{}

	viper.SetConfigName("configs")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic("panic in config parser : " + err.Error())
	}

	// default config
	viper.SetDefault("SERVER", SERVER{})
	viper.SetDefault("DB", DB{})
	viper.SetDefault("CORS", CORS{})
	viper.SetDefault("JWT", JWT{})
	viper.SetDefault("PPROF", PPROF{})
	viper.SetDefault("MAIL", MAIL{})
	viper.SetDefault("CARBIN", CARBIN{})
	viper.SetDefault("LabODT", LabODT{})

	viper.WriteConfig()

	err = viper.Unmarshal(conf)
	if err != nil {
		panic(err)
	}
	return conf
}
