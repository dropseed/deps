package env

import (
	"encoding/json"
	"fmt"
	"strings"
)

func SettingToEnviron(key string, value interface{}) (string, error) {
	envKey := settingKeyToEnv(key)
	envVal, err := settingValToEnv(value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s=%s", envKey, envVal), nil
}

func settingKeyToEnv(k string) string {
	return fmt.Sprintf("DEPS_SETTING_%s", strings.ToUpper(k))
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
