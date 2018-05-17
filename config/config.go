package config

import (
	"strings"

	"github.com/spf13/viper"
	"path/filepath"
	"fmt"
)

type ConfigFileType string

const (
	Yaml	ConfigFileType = "yaml"
)

type Config struct {
	viper *viper.Viper
}

func New() *Config {
	return &Config{viper.New()}
}

func (c *Config) Read(envPrefix, filename string) error {
	c.viper.SetEnvPrefix(envPrefix)
	c.viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	c.viper.SetEnvKeyReplacer(replacer)

	ext := filepath.Ext(filename)
	if ext != ".yml" && ext != "yaml"{
		return fmt.Errorf("configuration error, only yaml config file supported")
	}
	configDir, configFile := filepath.Split(filename)
	c.viper.SetConfigType(string(Yaml))
	c.viper.SetConfigName(strings.Replace(configFile, ext, "", 1))
	c.viper.AddConfigPath(configDir)

	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func (c* Config) HasKey(keys ...string) (interface{}, bool) {
	value := c.viper.Get(strings.Join(keys, "."))
	if value == nil {
		return nil, false
	}
	return value, true
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

func (c *Config) GetStringMap(keys ...string)  map[string]interface{} {
	return c.viper.GetStringMap(strings.Join(keys, "."))
}






