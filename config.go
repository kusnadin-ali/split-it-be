package main

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type listenConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (l listenConfig) Addr() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

type dbConfig struct {
	DBName   string `yaml:"db_name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode"`
}

func (d dbConfig) ConnStr() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}

type appConfig struct {
	Listen    listenConfig `yaml:"listen"`
	DBConfig  dbConfig     `yaml:"db"`
	JWTSecret string       `yaml:"jwt_secret"`
}

func defaultConfig() appConfig {
	return appConfig{
		Listen: listenConfig{Host: "127.0.0.1", Port: 8080},
		DBConfig: dbConfig{
			DBName:   "splitit",
			Host:     "127.0.0.1",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			SSLMode:  "disable",
		},
		JWTSecret: "change-me-in-production",
	}
}

// loadFromEnv overrides config values dari environment variable.
// Env var lebih prioritas dari config file — cocok untuk deployment.
func (c *appConfig) loadFromEnv() {
	if v := os.Getenv("LISTEN_HOST"); v != "" {
		c.Listen.Host = v
	}
	if v := os.Getenv("LISTEN_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			c.Listen.Port = p
		}
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		c.DBConfig.DBName = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		c.DBConfig.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			c.DBConfig.Port = p
		}
	}
	if v := os.Getenv("DB_USER"); v != "" {
		c.DBConfig.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		c.DBConfig.Password = v
	}
	if v := os.Getenv("DB_SSL_MODE"); v != "" {
		c.DBConfig.SSLMode = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		c.JWTSecret = v
	}
}

func loadConfigFromFile(filename string, cfg *appConfig) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(cfg)
}
