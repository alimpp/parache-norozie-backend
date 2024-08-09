package config

import (
	"ecom/pkg/constants"
	"errors"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Mode struct {
	Debug bool
}

type Log struct {
	Level int
}

type Server struct {
	Addr string
	Cert string
	Key  string
}

type DB struct {
	DriverName string
	Url        string
}

type Key struct {
	ID   string
	Cert string
}

type SMS struct {
	Url        string
	ApiKey     string
	TemplateId int
}

type OTP struct {
	TTE int
}

type Redis struct {
	Url string
}

type Swagger struct {
	Password string
}

type Observability struct {
	PProf      PProf
	Prometheus Prometheus
	Jaeger     Jaeger
}

type PProf struct {
	Enabled bool
}

type Prometheus struct {
	Enabled bool
}

type Jaeger struct {
	Enabled bool
}

type ConfStruct struct {
	Mode          Mode          `validate:"required"`
	Log           Log           `validate:"required"`
	Server        Server        `validate:"required"`
	DB            DB            `validate:"required"`
	SMS           SMS           `validate:"required"`
	OTP           OTP           `validate:"required"`
	Redis         Redis         `validate:"required"`
	Swagger       Swagger       `validate:"required"`
	Observability Observability `validate:"required"`
}

var Config ConfStruct

func (c ConfStruct) Validate() error {
	return validator.New().Struct(c)
}

func LoadConfig(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		constants.Logger.Error().Err(err).Msgf("config file not found")
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetEnvPrefix(constants.ServiceName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	return reload(file)
}

func reload(file string) bool {
	err := viper.MergeInConfig()
	if err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			constants.Logger.Error().Err(err).Msgf("config file not found %s", file)
		} else {
			constants.Logger.Error().Err(err).Msgf("config file read failed %s", file)
		}
		return false
	}

	err = viper.GetViper().UnmarshalExact(&Config)
	if err != nil {
		constants.Logger.Error().Err(err).Msgf("faild to unmarshal the conf %s", file)
		return false
	}

	if err = Config.Validate(); err != nil {
		constants.Logger.Error().Err(err).Msgf("invalid configuration %s", file)
	}

	constants.Logger.Info().Msgf("config file loaded %s", file)
	return true
}
