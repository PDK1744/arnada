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
	Host  string       `yaml:"host"`
	Paths []PathConfig `yaml:"paths"`
}

type PathConfig struct {
	Path     string `yaml:"path"`
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
	for _, r := range cfg.Routes {
		if r.Host == "" {
			return errors.New("route configuration contains an empty or missing host domain")
		}
		if len(r.Paths) == 0 {
			return fmt.Errorf("host %q must have at least one path configured", r.Host)
		}
		for _, p := range r.Paths {
			if p.Path == "" {
				return fmt.Errorf("host %q contains a route configuration with a missing path string", r.Host)
			}
			if p.Upstream == "" {
				return fmt.Errorf("path %q under host %q is missing its upstream target URL", p.Path, r.Host)
			}
		}
	}

	return nil
}
