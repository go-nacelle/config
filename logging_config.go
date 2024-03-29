package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Logger is an interface to the logger where config values are printed.
type Logger interface {
	// Printf logs a message. Arguments should be handled in the manner of fmt.Printf.
	Printf(format string, args ...interface{})
}

type nilLogger struct{}

func (nilLogger) Printf(format string, args ...interface{}) {}

func (c *Config) dumpSource() error {
	chunk := map[string]interface{}{}
	for key, value := range c.Dump() {
		if c.isMasked(key) {
			chunk[key] = "*****"
		} else {
			chunk[key] = value
		}
	}

	c.logger.Printf("Config source assets: %s", strings.Join(c.Assets(), ", "))
	c.logger.Printf("Config source contents: %s", normalizeChunk(chunk))
	return nil
}

func (c *Config) isMasked(target string) bool {
	for _, key := range c.maskedKeys {
		if strings.ToLower(key) == strings.ToLower(target) {
			return true
		}
	}

	return false
}

func dumpChunk(obj interface{}) (map[string]interface{}, error) {
	objValue, objType, err := getIndirect(obj)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	for i := 0; i < objType.NumField(); i++ {
		fieldType := objType.Field(i)
		fieldValue := objValue.Field(i)
		maskTagValue := fieldType.Tag.Get(MaskTag)
		displayTagValue := fieldType.Tag.Get(DisplayTag)
		displayName := fieldType.Name

		if displayTagValue != "" {
			displayName = displayTagValue
		}

		if maskTagValue != "" {
			val, err := strconv.ParseBool(maskTagValue)
			if err != nil {
				return nil, newSerializeError(fieldType.Name, ErrInvalidMaskTag)
			}

			if val {
				continue
			}
		}

		if fieldValue.Kind() == reflect.String {
			m[displayName] = fmt.Sprintf("%s", fieldValue)
		} else {
			data, err := json.Marshal(fieldValue.Interface())
			if err != nil {
				return nil, err
			}

			m[displayName] = string(data)
		}
	}

	return m, nil
}

func normalizeChunk(chunk map[string]interface{}) string {
	if len(chunk) == 0 {
		return "<no values>"
	}

	values := []string{}
	for key, value := range chunk {
		values = append(values, fmt.Sprintf("%s=%v", key, value))
	}

	sort.Strings(values)
	return "\n" + strings.Join(values, "\n")
}
