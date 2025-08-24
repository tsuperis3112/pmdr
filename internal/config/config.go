package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

const (
	ProjectName           = "pmdr"
	ConfigBaseName        = "config"
	DefaultConfigType     = "yaml"
	DefaultConfigFileName = ConfigBaseName + "." + DefaultConfigType
	LegacyConfigDirName   = ".pmdr"
	LegacyConfigBaseName  = ".pmdr"
)

// Duration is a wrapper around time.Duration for viper unmarshaling
// This is not strictly necessary with the decode hook, but can be useful
// for other purposes.
type Duration struct {
	time.Duration
}

// Config holds the application configuration
type Config struct {
	WorkDuration       time.Duration `mapstructure:"work_duration"`
	ShortBreakDuration time.Duration `mapstructure:"short_break_duration"`
	LongBreakDuration  time.Duration `mapstructure:"long_break_duration"`
	PomoCycles         int           `mapstructure:"pomo_cycles"`
	Hooks              Hook          `mapstructure:"hooks"`
}

// Hook represents a single hook command
type Hook struct {
	Work       []string `mapstructure:"work"`
	ShortBreak []string `mapstructure:"short_break"`
	LongBreak  []string `mapstructure:"long_break"`
}

// Load loads the configuration from viper
func Load() (*Config, error) {
	// Set default values
	vip := viper.GetViper()
	vip.SetDefault("work_duration", "25m")
	vip.SetDefault("short_break_duration", "5m")
	vip.SetDefault("long_break_duration", "15m")
	vip.SetDefault("pomo_cycles", 4)

	var config Config

	// Add the custom decode hook
	decodeHook := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)

	if err := vip.Unmarshal(&config, viper.DecodeHook(decodeHook)); err != nil {
		return nil, err
	}

	return &config, nil
}

// FindConfigFile finds the configuration file path.
func FindConfigFile(cfgFile string) (string, error) {
	if cfgFile != "" {
		return cfgFile, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		userConfigDir = filepath.Join(home, ".config")
	}

	configFilenames := []string{
		ConfigBaseName + ".yaml",
		ConfigBaseName + ".yml",
	}

	// Define search locations and the filenames to look for in each.
	searchLocations := []struct {
		dir       string
		filenames []string
	}{
		{".", []string{LegacyConfigBaseName + ".yaml", LegacyConfigBaseName + ".yml"}},
		{filepath.Join(home, LegacyConfigDirName), configFilenames},
		{filepath.Join(userConfigDir, ProjectName), configFilenames},
	}

	for _, loc := range searchLocations {
		for _, filename := range loc.filenames {
			path := filepath.Join(loc.dir, filename)
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}
	}

	return "", nil
}

// GetDefaultConfigPaths returns the default directory and file path for the configuration.
func GetDefaultConfigPaths() (string, string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", "", err
	}
	configDir := filepath.Join(userConfigDir, ProjectName)
	configFile := filepath.Join(configDir, DefaultConfigFileName)
	return configDir, configFile, nil
}

// GetConfigFilePath returns the path to the configuration file that viper is using.
// If no file is used, it returns the default path.
func GetConfigFilePath() (string, error) {
	if configFile := viper.ConfigFileUsed(); configFile != "" {
		return configFile, nil
	}

	// If no config file is in use, return the default path
	_, defaultPath, err := GetDefaultConfigPaths()
	if err != nil {
		return "", err
	}
	return defaultPath, nil
}
