package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Routes  []RouteConfig `yaml:"routes"`
	Logging LoggingConfig `yaml:"logging"`
}

type ServerConfig struct {
	Listen string `yaml:"listen"`
}

type RouteConfig struct {
	Host     string `yaml:"host"`
	Upstream string `yaml:"upstream"`
}

type LoggingConfig struct {
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("config file not loaded: %v", err)
		//return nil, err
	}

	var cfg Config
	err = parseYAML(data, &cfg)
	if err != nil {
		log.Fatalf("invalid yaml parse: %v", err)
	}

	err = validateConfig(&cfg)
	if err != nil {
		log.Fatalf("INVALID CONFIG: %v", err)
	}
	//returnConfig(&cfg)

	return &cfg, nil
}

func parseYAML(data []byte, cfg *Config) error {
	return yaml.Unmarshal(data, &cfg)
}

func validateConfig(cfg *Config) error {
	if cfg.Server.Listen == "" {
		return errors.New("Server Listen address not found")
	}
	if len(cfg.Routes) == 0 {
		return errors.New("No routes found")
	}
	for i := range cfg.Routes {
		route := &cfg.Routes[i]
		if route.Host == "" {
			return fmt.Errorf("No host found for route: %v", i+1)
		}
		if route.Upstream == "" {
			return fmt.Errorf("No backend found for route: %v", i)
		}
	}

	return nil
}

// func returnConfig(cfg *Config) {
// 	fmt.Printf("SERVER LISTEN: %v\n", cfg.Server.Listen)
// 	fmt.Printf("ROUTES: \n")
// 	for i, r := range cfg.Routes {
// 		fmt.Printf("ROUTE %v\n", i)
// 		fmt.Printf("ROUTE HOST: %v\n", r.Host)
// 		fmt.Printf("ROUTE BACKEND: %v\n", r.Backend)
// 	}
// 	fmt.Printf("LOGGING FORMAT: %v\n", cfg.Logging.Format)
// 	fmt.Printf("LOGGING OUTPUT: %v\n", cfg.Logging.Output)
// }
