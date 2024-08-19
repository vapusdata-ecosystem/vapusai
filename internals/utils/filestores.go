package utils

import (
	os "os"
	filepath "path/filepath"

	toml "github.com/pelletier/go-toml/v2"
	dmlogger "github.com/vapusdata-ecosystem/vapusai-studio/internals/logger"
)

// Function to write toml file
func WriteTomlFile(data interface{}, filename, path string) error {
	bytes, err := toml.Marshal(data)
	if err != nil {
		return err
	}

	file := filepath.Join(path, filename+DOT+DEFAULT_CONFIG_TYPE)
	dmlogger.CoreLogger.Info().Msgf("Writing to file: %v", file)
	err = os.WriteFile(file, bytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

// Function to read toml file
func ReadTomlFile(data interface{}, filename, path string) error {
	file := filepath.Join(path, filename+DOT+DEFAULT_CONFIG_TYPE)
	dmlogger.CoreLogger.Info().Msgf("Reading from file: %v", file)
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(bytes, data)
	if err != nil {
		return err
	}
	return nil
}
