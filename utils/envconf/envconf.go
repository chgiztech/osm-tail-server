package envconf

import (
	"errors"
	"fmt"
	"log"
	"os"
	"osm-tail/utils/validation"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	AppName = "osm-tail"
)

var (
	App AppConfig
)

type AppConfig struct {
	Port                   int  `env:"PORT" default:"3000"`
	EnableHeaderValidation bool `env:"ENABLE_HEADER_VALIDATION" default:"false"`
	PostgreSQL             PostgreSQLConfig
}

type PostgreSQLConfig struct {
	Host     string `env:"POSTGRES_HOST" default:"localhost"`
	Port     string `env:"POSTGRES_PORT" default:"5432"`
	User     string `env:"POSTGRES_USER" default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" default:"1234567"`
	Database string `env:"POSTGRES_DB" default:"postgres"`
	LogLevel string `env:"POSTGRES_LOGLEVEL" default:"warn"`
}

func LoadAppConfig() error {
	// Load service config primitives
	if err := loadFromEnv(&App); err != nil {
		return err
	}

	// Load nested configs
	loaders := []func() error{
		func() error { return loadFromEnv(&App.PostgreSQL) },
	}

	for _, load := range loaders {
		if err := load(); err != nil {
			return err
		}
	}

	err := validation.Validate.Struct(App)
	if err != nil {
		return err
	}

	return nil
}

// loadFromEnv populates a struct from environment variables using struct tags
func loadFromEnv(config interface{}) error {
	configValue := reflect.ValueOf(config)
	if configValue.Kind() != reflect.Ptr || configValue.Elem().Kind() != reflect.Struct {
		return errors.New("config must be a pointer to a struct")
	}
	configElem := configValue.Elem()
	configType := configElem.Type()

	err := godotenv.Load()

	if err != nil {
		fmt.Errorf("Error loading .env file: %v", err)
	}

	for i := 0; i < configElem.NumField(); i++ {
		field := configElem.Field(i)
		fieldType := configType.Field(i)

		// Skip if field is a struct (for nested config structs)
		if field.Kind() == reflect.Struct {
			continue
		}

		envKey := fieldType.Tag.Get("env")
		defaultValue := fieldType.Tag.Get("default")

		// Skip if no env tag
		if envKey == "" {
			continue
		}

		// Get value from environment or use default
		value := os.Getenv(envKey)

		if value == "" {
			value = defaultValue
			log.Printf("Config: %s=%s (default)", envKey, maskSensitiveValue(envKey, value))
		} else {
			log.Printf("Config: %s=%s (env)", envKey, maskSensitiveValue(envKey, value))
		}

		// Set field value based on its type
		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid integer value for %s: %w", envKey, err)
			}
			field.SetInt(intValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid boolean value for %s: %w", envKey, err)
			}
			field.SetBool(boolValue)
		default:
			return fmt.Errorf("unsupported field type for %s", envKey)
		}
	}

	return nil
}

func maskSensitiveValue(key, value string) string {
	// List of keys that contain sensitive information
	sensitiveKeys := []string{
		"PASSWORD", "SECRET", "KEY", "TOKEN", "CREDENTIALS",
	}

	// Check if the key contains any sensitive words
	for _, sensitiveKey := range sensitiveKeys {
		if strings.Contains(strings.ToUpper(key), sensitiveKey) {
			if value == "" {
				return ""
			}
			return "********"
		}
	}

	return value
}
