package config

import "github.com/feimingxliu/quicksearch/pkg/viper"

var Global = new(Config)

//Init init the config.
func Init(path string) error {
	c := viper.NewConfig(path)
	return c.LoadInto(Global)
}

type Config struct {
	Env     string  `mapstructure:"env" json:"env" yaml:"env"`
	Storage Storage `mapstructure:"storage" json:"storage" yaml:"meta-storage"`
	Http    Http    `mapstructure:"http" json:"http" yaml:"http"`
}

type Storage struct {
	MetaType string `mapstructure:"meta-type" json:"meta_type" yaml:"meta-type"`
	DataDir  string `mapstructure:"data-dir" json:"data_dir" yaml:"data-dir"`
}

type Http struct {
	Addr string `mapstructure:"addr" json:"addr" yaml:"addr"`
}
