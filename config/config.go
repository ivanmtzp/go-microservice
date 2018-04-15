package config

import (
	"strings"

	"github.com/spf13/viper"
)

type ConfigFileType string

const (
	Yaml	ConfigFileType = "yaml"
	Toml    ConfigFileType = "toml"
)

const (
	local	string = "local"
	dev    string = "dev"
	stage  string = "stage"
	prod string = "prod"
)

var environment string = local

func Environment() string {
	return environment
}

func Read(envPrefix, configFilePath, configFileName string, configFileType ConfigFileType) error {
	viper.SetDefault("environment", local)
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetConfigType(string(configFileType))
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	environment = viper.GetString("environment")
	return nil
}

func GetString(keys ...string) string {
	return viper.GetString(environment + "." + strings.Join(keys, "."))
}

func GetInt(keys ...string) int {
	return viper.GetInt(environment + "." + strings.Join(keys, "."))
}

func GetBool(keys ...string) bool {
	return viper.GetBool(environment + "." + strings.Join(keys, "."))
}