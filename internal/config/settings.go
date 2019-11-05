package config

import (
	"strings"

	"github.com/dropseed/deps/internal/env"
)

type Settings map[string]interface{}

func (s Settings) Get(name string) interface{} {
	for k, v := range s {
		if strings.ToLower(k) == strings.ToLower(name) {
			return v
		}
	}

	return nil
}

func (s Settings) AsEnviron() []string {
	environ := []string{}

	for k, v := range s {
		environString, err := env.SettingToEnviron(k, v)
		if err != nil {
			panic(err)
		}
		environ = append(environ, environString)
	}

	return environ
}
