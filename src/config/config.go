package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

type Config struct {
	Port         int          `yaml:"port"`
	Path         string       `yaml:"path"`
	MongoUrl     string       `yaml:"mongoUrl"`
	CacheOptions CacheOptions `yaml:"cacheOptions"`
}

type CacheOptions struct {
	Expiration int `yaml:"expiration"`
	Purge      int `yaml:"purge"`
}

const ConfigPathEnvName = "AGENT_MANAGER_CONFIG_PATH"

func NewConfig() Config {
	configPath := "./config.yaml"

	if os.Getenv(ConfigPathEnvName) != "" {
		configPath = os.Getenv(ConfigPathEnvName)
	}

	cfg, err := getConfig(configPath)
	if err != nil {
		log.Fatalf("could not parse config: %v", err)
	}
	return cfg
}

func (c *Config) validateAndFillDefaults() error {
	if c.Port <= 0 {
		c.Port = 8080
	}
	if c.MongoUrl == "" {
		return errors.New("no mongoUrl provided in config")
	}
	err := c.CacheOptions.validateAndFillDefaults()
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheOptions) validateAndFillDefaults() error {
	if c.Expiration <= 0 {
		c.Expiration = 5
	}
	if c.Purge <= 0 {
		c.Purge = 10
	}
	return nil
}

func getConfig(path string) (Config, error) {
	var res Config

	file, err := os.Open(path)
	if err != nil {
		return res, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return res, err
	}

	err = yaml.Unmarshal(bytes, &res)
	if err != nil {
		return res, err
	}

	err = res.validateAndFillDefaults()
	return res, err
}
