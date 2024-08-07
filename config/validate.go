package config

import "github.com/go-playground/validator/v10"

func (c *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(c)

	if err != nil {
		return err
	}

	return nil
}
