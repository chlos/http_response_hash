package main

import (
	"fmt"
	"os"

	"http_response_hash/config"
	"http_response_hash/hashing"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Failed to get configuration: %v\n", err)
		os.Exit(1)
	}

	h, err := hashing.NewHashing(config.ParallelLimit, config.URLs)
	if err != nil {
		fmt.Printf("Failed to init hashing: %v\n", err)
		os.Exit(1)
	}
	h.Start()
	h.Print()
}
