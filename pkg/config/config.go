package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	AUTHTOKEN  string `mapstructure:"TWILIO_AUTH_TOKEN"`
	ACCOUNTSID string `mapstructure:"TWILIO_ACCOUNT_SID"`
	SERVICESID string `mapstructure:"TWILIO_SERVICE_SID"`

	JWT string `mapstructure:"JWT_CODE"`

	RazorpayAPIKeyID     string `mapstructure:"RAZORPAY_API_KEY_ID"`
	RazorpayAPIKeySecret string `mapstructure:"RAZORPAY_API_KEY_SECRET"`
}

// to hold all names of env variables
var envs = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD", // database
	"JWT_CODE",                                                      //JWT
	"TWILIO_AUTH_TOKEN", "TWILIO_ACCOUNT_SID", "TWILIO_SERVICE_SID", //twilio details
	"RAZORPAY_API_KEY_ID", "RAZORPAY_API_KEY_SECRET", //razor pay
}

var config Config

func LoadConfig() (Config, error) {
	//var config Config

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}

	fmt.Println("Config is", config)

	return config, nil
}

func GetConfig() Config {
	return config
}

// to get the secret code for jwt
func GetJWTConfig() string {
	return config.JWT
}
