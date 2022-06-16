package viper

import (
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type Config struct {
	path string
	v    *viper.Viper
}

//NewConfig return a new Config.
func NewConfig(path string) *Config {
	v := viper.New()
	ss := strings.Split(filepath.Base(path), ".")
	if len(ss) != 2 {
		panic("invalid config file.")
	}
	v.SetConfigName(ss[0])
	v.SetConfigType(ss[1])
	v.SetConfigFile(path)
	return &Config{
		path: path,
		v:    v,
	}
}

//LoadInto loads config into val.
func (c *Config) LoadInto(val interface{}) error {
	if err := c.v.ReadInConfig(); err != nil {
		return err
	}
	return c.v.Unmarshal(&val)
}
