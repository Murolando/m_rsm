package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	EmptyPath = "EMPTY_PATH"
)

type Config struct {
	Enviroment string `yaml:"env"`
	Line       int64  `yaml:"line"`
}

func ParseConfig() (*Config, error) {

	cfg := &Config{}
	// TODO: make lib for parsing?
	// Get config path from flags, if flags not exist get from env
	path := pflag.String("app-cfg", EmptyPath, "config path")

	pflag.Parse()
	if *path == EmptyPath {
		viper.AutomaticEnv()
		*path = viper.Get("APP_CFG_PATH").(string)
	}

	data, err := os.ReadFile(*path)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return cfg, nil
}
