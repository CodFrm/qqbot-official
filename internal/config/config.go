package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AppID       uint64 `yaml:"appId"`
	AccessToken string `yaml:"accessToken"`
}

var AppConfig Config

func Init(filename string) error {
	b, _ := os.ReadFile(filename)
	if err := yaml.Unmarshal(b, &AppConfig); err != nil {
		return err
	}
	return nil
}
