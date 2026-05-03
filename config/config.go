package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Servers      []ServerConfig `yaml:"servers"`
	ActiveServer string         `yaml:"active_server"`
	PollInterval int            `yaml:"poll_interval_ms"`
}

type ServerConfig struct {
	Name       string     `yaml:"name"`
	URL        string     `yaml:"url"`
	Auth       AuthConfig `yaml:"auth"`
	APIVersion string     `yaml:"api_version"`
}

type AuthConfig struct {
	Type     string `yaml:"type"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Token    string `yaml:"token,omitempty"`
	Profile  string `yaml:"profile,omitempty"`
	Region   string `yaml:"region,omitempty"`
}

func ConfigPath() string {
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		usr, err := user.Current()
		if err != nil {
			home, _ := os.UserHomeDir()
			xdgConfig = filepath.Join(home, ".config")
		} else {
			xdgConfig = filepath.Join(usr.HomeDir, ".config")
		}
	}
	return filepath.Join(xdgConfig, "airflow-tui", "config.yaml")
}

func LoadConfig(path string) *Config {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{PollInterval: 2000}
	}

	if err := v.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		return &Config{PollInterval: 2000}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config: %v\n", err)
		return &Config{PollInterval: 2000}
	}

	return &cfg
}

func SaveConfig(path string, cfg *Config) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(path)
	v.Set("servers", cfg.Servers)
	v.Set("active_server", cfg.ActiveServer)
	v.Set("poll_interval_ms", cfg.PollInterval)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		return
	}

	if err := v.WriteConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing config: %v\n", err)
	}
}
