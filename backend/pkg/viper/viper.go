package viper

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Path string
	Name string
	Type Type
}

type Type string

const (
	YML  Type = "yml"
	JSON Type = "json"
	TOML Type = "toml"
)

func Init[T any](c Config, o *T) error {
	if os.Getenv("LOAD_DEV_ENV") != "0" {
		err := godotenv.Load(".env.dev")
		if err == nil {
			fmt.Println("Loaded dev env")
		}
	}

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("env file could not be loaded: %s", err.Error())
	}

	v := viper.New()

	v.AddConfigPath(c.Path)
	v.SetConfigName(c.Name)
	v.SetConfigType(string(c.Type))

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("config file could not be read: %s", err.Error())
	}

	for _, k := range v.AllKeys() {
		val := v.GetString(k)
		v.Set(k, os.ExpandEnv(val))
	}

	if err := v.Unmarshal(o); err != nil {
		return fmt.Errorf("config file could not be unmarshaled: %s", err.Error())
	}

	validate := validator.New()
	err = validate.Struct(o)
	if err != nil {
		return err
	}

	return nil
}
