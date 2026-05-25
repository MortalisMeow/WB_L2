package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Port string
}

func LoadConfig() Config {
	port := os.Getenv("CALENDAR_PORT")
	if port == "" {
		port = "8080"
	}

	flagPort := flag.String("port", port, "server port")
	flag.Parse()

	return Config{
		Port: *flagPort,
	}
}

func (c Config) Address() string {
	return fmt.Sprintf(":%s", c.Port)
}
