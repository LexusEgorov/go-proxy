package config

//TODO: read env

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"

	"github.com/LexusEgorov/go-proxy/internal/models"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type Interval struct {
	MinMilliseconds int `yaml:"min"`
	MaxMilliseconds int `yaml:"max"`
}

type ClientConfig struct {
	Interval Interval `yaml:"interval"`
	URL      string   `yaml:"url"`
	Factor   int      `yaml:"factor"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Client ClientConfig `yaml:"client"`
}

func New() (*Config, error) {
	var (
		cfg *Config
		err error
	)

	configPath, err := fetchConfigPath()

	if err != nil {
		if !errors.Is(err, models.ErrConfigPathNotProvided) {
			return nil, fmt.Errorf("read config error: %v", err)
		}

		cfg, err = readEnvConfig()
	} else {
		cfg, err = readFileConfig(configPath)
	}

	if err != nil {
		return nil, err
	}

	err = checkConfig(cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Валидирует конфиг
func checkConfig(cfg *Config) error {
	if err := checkServerConfig(&cfg.Server); err != nil {
		return err
	}

	return checkClientConfig(&cfg.Client)
}

// Валидирует клиентский конфиг
func checkClientConfig(cfg *ClientConfig) error {
	if cfg.Factor <= 1 {
		return models.ErrBadConfigFactor
	}

	if cfg.URL == "" {
		return models.ErrBadConfigURL
	}

	if cfg.Interval.MinMilliseconds <= 0 {
		return models.ErrBadConfigMinInterval
	}

	if cfg.Interval.MaxMilliseconds < 0 {
		return models.ErrBadConfigMaxInterval
	}

	if cfg.Interval.MaxMilliseconds < cfg.Interval.MinMilliseconds {
		return models.ErrBadConfigMinMaxInterval
	}

	return nil
}

// Валидирует серверный конфиг
func checkServerConfig(cfg *ServerConfig) error {
	if cfg.Port <= 0 {
		return models.ErrBadConfigPort
	}

	return nil
}

// Читает конфиг из env
func readEnvConfig() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))

	if err != nil {
		return nil, err
	}

	fmt.Printf("port: %d\n", port)

	minInterval, err := strconv.Atoi(os.Getenv("MIN_INTERVAL"))

	if err != nil {
		return nil, err
	}

	maxInterval, err := strconv.Atoi(os.Getenv("MAX_INTERVAL"))

	if err != nil {
		return nil, err
	}

	factor, err := strconv.Atoi(os.Getenv("FACTOR"))

	if err != nil {
		return nil, err
	}

	return &Config{
		Server: ServerConfig{
			Port: port,
		},
		Client: ClientConfig{
			Interval: Interval{
				MinMilliseconds: minInterval,
				MaxMilliseconds: maxInterval,
			},
			Factor: factor,
			URL:    os.Getenv("DESTINATION_URL"),
		},
	}, nil
}

// Читает конфиг из файла
func readFileConfig(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, err
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Получает путь до конфигурационного файла через флаг или env
func fetchConfigPath() (string, error) {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		err := godotenv.Load()

		if err != nil {
			return "", err
		}

		path = os.Getenv("CONFIG_PATH")

		if path == "" {
			return "", models.ErrConfigPathNotProvided
		}
	}

	return path, nil
}
