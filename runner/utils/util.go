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
			return
		}
		_ = json.Unmarshal([]byte(value), &ss)
	case []string:
		return
	case []any:
		if len(value) == 0 {
			return
		}
		data, _ := json.Marshal(value)
		_ = json.Unmarshal(data, &ss)
	}
	return
}

func SplitCommand(command string) (name string, args []string) {
	fields := strings.Fields(command)
	if len(fields) == 1 {
		return fields[0], nil
	}
	return fields[0], fields[1:]
}