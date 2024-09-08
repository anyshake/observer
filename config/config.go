package config

import (
	"encoding/json"
	"os"
)

func (c *Config) Read(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return err
	}

	return nil
}
