package config

import (
	"errors"
	"flag"
	"fmt"
)

// _defaultParallelLimit is a default value for the number of parallel requests.
const _defaultParallelLimit = 10

var limit int

func init() {
	flag.IntVar(&limit, "parallel", _defaultParallelLimit, "The limit number of parallel requests")
}

// Config is a struct containing tool's settings.
type Config struct {
	// ParallelLimit is a value for the number of parallel requests.
	ParallelLimit int
	// URLs is a slice of URLs of websites which should be hashed.
	URLs []string
}

// NewConfig is a constructor for the Config.
func NewConfig() (*Config, error) {
	config := &Config{}
	err := parseFlags(config)
	if err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	return config, nil
}

// parseFlags parses cmd line flags and args and updates config structure with the corresponding values.
func parseFlags(config *Config) error {
	flag.Parse()

	config.ParallelLimit = limit

	if len(flag.Args()) == 0 {
		return errors.New("no URLs specified")
	}
	config.URLs = append(config.URLs, flag.Args()...)

	return nil
}
