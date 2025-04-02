package utils

import (
	"encoding/json"
	"os"
	"strings"
)

func IsFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		return false
	}

	return true
}

func GetStringArray(s any) (ss []string) {
	switch value := s.(type) {
	case string:
		if value == "" || value == "[]" {
			return ss
		}

		err := json.Unmarshal([]byte(value), &ss)
		if err != nil {
			return nil
		}

		return ss
	case []string:
		return value
	case []any:
		if len(value) == 0 {
			return ss
		}

		data, err := json.Marshal(value)
		if err != nil {
			return nil
		}

		err = json.Unmarshal(data, &ss)
		if err != nil {
			return nil
		}

		return ss
	}

	return nil
}

func CastToBool(v any) bool {
	switch v := v.(type) {
	case string:
		return v == "true"
	case bool:
		return v
	}

	return false
}

func SplitCommand(command string) (name string, args []string) {
	fields := strings.Fields(command)
	if len(fields) == 1 {
		return fields[0], nil
	}

	return fields[0], fields[1:]
}
