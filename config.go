package main

import (
	"fmt"
	"log"
	"time"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Initialize this variable to access the env values
var k = koanf.New(".")

type envConfigs struct {
	RuntimeEnv       string        `koanf:"RUNTIME_ENV"`
	DBType           string        `koanf:"DB_TYPE"`
	DBHost           string        `koanf:"DB_HOST"`
	DBPort           string        `koanf:"DB_PORT"`
	DBUsername       string        `koanf:"DB_USERNAME"`
	DBPassword       string        `koanf:"DB_PASSWORD"`
	DBName           string        `koanf:"DB_NAME"`
	RedisAddr        string        `koanf:"REDIS_ADDR"`
	RedisPassword    string        `koanf:"REDIS_PASSWORD"`
	RedisDB          int           `koanf:"REDIS_DB"`
	Port             string        `koanf:"PORT"`
	JWTSecret        string        `koanf:"JWT_SECRET"`
	JWTValid         time.Duration `koanf:"JWT_VALID"`
	JWTIssuer        string        `koanf:"JWT_ISSUER"`
	JWTRefreshSecret string        `koanf:"JWT_REFRESH_SECRET"`
	JWTRefreshValid  time.Duration `koanf:"JWT_REFRESH_VALID"`
}

func LoadEnvVariables(path string) (config *envConfigs, err error) {
	// Default Values
	err = k.Load(confmap.Provider(map[string]interface{}{
		"RUNTIME_ENV":        "production",
		"DB_TYPE":            "mysql",
		"DB_HOST":            "localhost",
		"DB_PORT":            "3306",
		"DB_USERNAME":        "root",
		"DB_PASSWORD":        "",
		"DB_NAME":            "gomailadmin",
		"REDIS_ADDR":         "localhost:6379",
		"REDIS_PASSWORD":     "",
		"REDIS_DB":           0,
		"PORT":               "3000",
		"JWT_SECRET":         "",
		"JWT_VALID":          60,
		"JWT_ISSUER":         "govmimbadmin",
		"JWT_REFRESH_SECRET": "",
		"JWT_REFRESH_VALID":  60 * 24 * 30,
	}, "."), nil)
	if err != nil {
		return nil, err
	}

	// read .env file
	if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
		log.Printf("Cannot load .env, continue: %v", err)
		return nil, err
	}
	// overwrite from environment
	if err := k.Load(env.Provider("", "=", nil), nil); err != nil {
		log.Fatalf("error loading environment: %v", err)
		return nil, err
	}
	err = k.Unmarshal("", &config)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Config: %+v\n", config)
	return config, nil
}
