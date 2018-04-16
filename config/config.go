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

func Read(envPrefix, configFilePath, configFileName string, configFileType ConfigFileType) error {
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

	return nil
}

func GetString(keys ...string) string {
	return viper.GetString(strings.Join(keys, "."))
}

func GetInt(keys ...string) int {
	return viper.GetInt(strings.Join(keys, "."))
}

func GetBool(keys ...string) bool {
	return viper.GetBool(strings.Join(keys, "."))
}