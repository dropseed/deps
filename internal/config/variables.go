package config

import (
	"fmt"
	"reflect"
	"strings"
)

var tempConfigVariables map[string]interface{}

func resetVariables() {
	tempConfigVariables = nil
}

func variableMapDecode(from reflect.Kind, to reflect.Kind, value interface{}) (interface{}, error) {
	if from == reflect.String {
		if s, ok := value.(string); ok && strings.HasPrefix(s, "$") {
			variableValue, ok := tempConfigVariables[s[1:]]
			if !ok {
				return nil, fmt.Errorf("%s variable not found", s)
			}
			return variableValue, nil
		}
	}

	return value, nil
}
