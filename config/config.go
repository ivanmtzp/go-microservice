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

type Config struct {
	viper *viper.Viper
}

func New() *Config {
	return &Config{viper.New()}
}

func (c *Config) Read(envPrefix, configFilePath, configFileName string, configFileType ConfigFileType) error {
	c.viper.SetEnvPrefix(envPrefix)
	c.viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	c.viper.SetEnvKeyReplacer(replacer)

	c.viper.SetConfigType(string(configFileType))
	c.viper.SetConfigName(configFileName)
	c.viper.AddConfigPath(configFilePath)

	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func (c *Config) GetString(keys ...string) string {
	return c.viper.GetString(strings.Join(keys, "."))
}

func (c *Config) GetInt(keys ...string) int {
	return c.viper.GetInt(strings.Join(keys, "."))
}

func (c *Config) GetBool(keys ...string) bool {
	return c.viper.GetBool(strings.Join(keys, "."))
}