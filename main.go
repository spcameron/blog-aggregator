// main.go

package main

import (
	"fmt"
	"log"

	"github.com/spcameron/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Read() error: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	err = cfg.SetUser("Sean")
	if err != nil {
		log.Fatalf("SetUser() error: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Read() error: %v", err)
	}
	fmt.Printf("Read config after update: %+v", cfg)
}
