package a3spow

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/yaml.v2"
)

type Config struct {
	HD     bool               `yaml:"hd"`
	Owner  common.Address     `yaml:"owner"`
	Filter LongRepeatedFilter `yaml:"filter"`
	Number int32              `yaml:"number"`
}

func MustReadConfig(filename string) Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	return cfg
}
