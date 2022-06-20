package config

import "github.com/feimingxliu/quicksearch/pkg/viper"

var Global = new(Config)

//Init init the config.
func Init(path string) error {
	c := viper.NewConfig(path)
	return c.LoadInto(Global)
}

type Config struct {
	Storage   Storage   `mapstructure:"storage" json:"storage" yaml:"storage"`
	Tokenizer Tokenizer `mapstructure:"tokenizer" json:"tokenizer" yaml:"tokenizer"`
}

type Storage struct {
	Type    string `mapstructure:"type" json:"type" yaml:"type"`
	DataDir string `mapstructure:"data-dir" json:"data_dir" yaml:"data-dir"`
}

type Tokenizer struct {
	Type string `mapstructure:"type" json:"type" yaml:"type"`
}
