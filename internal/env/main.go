package env

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dropseed/deps/internal/output"
)

const SETTING_PREFIX = "DEPS_SETTING_"

func SettingToEnviron(name string, value interface{}) (string, error) {
	envKey := settingNameToKey(name)
	envVal, err := settingValToEnv(value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s=%s", envKey, envVal), nil
}

func SettingFromEnviron(name string) interface{} {
	v := os.Getenv(settingNameToKey(name))
	if v != "" {
		var data interface{}
		if err := json.Unmarshal([]byte(v), &data); err != nil {
			output.Error("Setting \"%s\" from env was not valid JSON")
			panic(err)
		}
		return data
	}
	return nil
}

func settingNameToKey(k string) string {
	return fmt.Sprintf("%s%s", SETTING_PREFIX, strings.ToUpper(k))
}

func settingValToEnv(v interface{}) (string, error) {
	envVarVal := v

	// if it's not already a string, json encode it
	// encoding a string seems to double quote
	if _, ok := v.(string); !ok {
		tmp, jsonErr := json.Marshal(v)
		if jsonErr != nil {
			return "", jsonErr
		}
		envVarVal = string(tmp)
	}

	return envVarVal.(string), nil
}
