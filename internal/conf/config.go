package conf

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/fengyuanluo/go-backend-template/internal/utils/log"
	"github.com/spf13/viper"
)

type Server struct {
	Host string `json:"host" mapstructure:"host"`
	Port int    `json:"port" mapstructure:"port"`
}

type Log struct {
	Level string `json:"level" mapstructure:"level"`
}

type Database struct {
	Type string `json:"type" mapstructure:"type"`
	Path string `json:"path" mapstructure:"path"`
}

type Auth struct {
	Secret string `json:"secret" mapstructure:"secret"`
}

type Config struct {
	Server   Server   `mapstructure:"server"`
	Log      Log      `mapstructure:"logging"`
	Database Database `mapstructure:"database"`
	Auth     Auth     `mapstructure:"auth"`
}

var AppConfig Config

func Load(path string) error {
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		viper.AddConfigPath("data")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix(APP_NAME)

	setDefaults()

	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Infof("Config file not found, creating default config")
			if err := os.MkdirAll("data", 0755); err != nil {
				log.Errorf("Failed to create data directory: %v", err)
			}
			if err := ensureAuthSecret(); err != nil {
				return err
			}
			if err := viper.SafeWriteConfigAs("data/config.json"); err != nil {
				log.Errorf("Failed to create default config: %v", err)
			}
		} else {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	if err := ensureAuthSecret(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}
	return nil
}

func setDefaults() {
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.path", "data/data.db")
	viper.SetDefault("logging.level", "info")
}

func ensureAuthSecret() error {
	secret := strings.TrimSpace(viper.GetString("auth.secret"))
	if secret != "" {
		return nil
	}
	generated, err := generateAuthSecret()
	if err != nil {
		return err
	}
	viper.Set("auth.secret", generated)
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		return nil
	}
	if _, err := os.Stat(configPath); err == nil {
		if err := viper.WriteConfig(); err != nil {
			log.Warnf("Failed to persist auth secret to config: %v", err)
		}
	}
	return nil
}

func generateAuthSecret() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(raw), nil
}
