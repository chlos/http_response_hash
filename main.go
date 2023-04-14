package main

import (
	"fmt"
	"os"

	"http_response_hash/config"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Failed to get configuration: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("limit: %d\n", config.ParallelLimit) // FIXME:
	fmt.Printf("args: %s\n", config.URLs)           // FIXME:
}
