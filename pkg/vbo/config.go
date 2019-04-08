package vbo

import (
	"fmt"

	"github.com/pelletier/go-toml"
)

func LoadConfig(path string) *Config {
	tree, err := toml.LoadFile(path)
	if err != nil {
		fmt.Printf("Unable to load config from %s, some features will not work", path)
		return nil
	}

	circuitTools := tree.Get("circuittools.path").(string)
	port := tree.Get("server.port").(int64)

	return &Config{
		CTPath: circuitTools,
		Port:   port}
}

type Config struct {
	CTPath string
	Port   int64
}
