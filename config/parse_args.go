package config

import (
	"errors"
	"flag"
	"fmt"
)

const _defaultParallelLimit = 10

var limit int

func init() {
	flag.IntVar(&limit, "parallel", _defaultParallelLimit, "The limit number of parallel requests")
}

type Config struct {
	ParallelLimit int
	URLs          []string
}

func NewConfig() (*Config, error) {
	config := &Config{}
	err := parseFlags(config)
	if err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	return config, nil
}

func parseFlags(config *Config) error {
	flag.Parse()

	config.ParallelLimit = limit

	if len(flag.Args()) == 0 {
		return errors.New("no URLs specified")
	}
	config.URLs = append(config.URLs, flag.Args()...)

	return nil
}
