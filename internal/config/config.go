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
	Engine  Engine  `mapstructure:"engine" json:"engine" yaml:"engine"`
	Storage Storage `mapstructure:"storage" json:"storage" yaml:"storage"`
	Http    Http    `mapstructure:"http" json:"http" yaml:"http"`
}

type Engine struct {
	DefaultNumberOfShards   int `mapstructure:"default-number-of-shards" json:"default_number_of_shards" yaml:"default-number-of-shards"`
	DefaultBatchSize        int `mapstructure:"default-batch-size" json:"default_batch_size" yaml:"default-batch-size"`
	DefaultSearchResultSize int `mapstructure:"default-search-result-size" json:"default_search_result_size" yaml:"default-search-result-size"`
}

type Storage struct {
	MetaType string `mapstructure:"meta-type" json:"meta_type" yaml:"meta-type"`
	DataDir  string `mapstructure:"data-dir" json:"data_dir" yaml:"data-dir"`
}

type Http struct {
	ServerAddr string `mapstructure:"server-addr" json:"server_addr" yaml:"server-addr"`
}
