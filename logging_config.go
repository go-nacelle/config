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

type loggingConfig struct {
	Config
	logger     Logger
	maskedKeys []string
}

// NewLoggingConfig wraps a config object with logging. After each successful load,
// the populated configuration object is serialized as fields and output at the info
// level.
func NewLoggingConfig(config Config, logger Logger, maskedKeys []string) Config {
	return &loggingConfig{
		Config:     config,
		logger:     logger,
		maskedKeys: maskedKeys,
	}
}

func (c *loggingConfig) Load(target interface{}, modifiers ...TagModifier) error {
	if err := c.Config.Load(target, modifiers...); err != nil {
		c.dumpSource()
		return err
	}

	chunk, err := dumpChunk(target)
	if err != nil {
		return fmt.Errorf("failed to serialize config (%s)", err.Error())
	}

	c.logger.Printf("Config loaded: %s", normalizeChunk(chunk))
	return nil
}

func (c *loggingConfig) MustLoad(target interface{}, modifiers ...TagModifier) {
	if err := c.Load(target, modifiers...); err != nil {
		panic(err.Error())
	}
}

func (c *loggingConfig) dumpSource() error {
	chunk := map[string]interface{}{}
	for key, value := range c.Config.Dump() {
		if c.isMasked(key) {
			chunk[key] = "*****"
		} else {
			chunk[key] = value
		}
	}

	c.logger.Printf("Config source assets: %s", strings.Join(c.Config.Assets(), ", "))
	c.logger.Printf("Config source contents: %s", normalizeChunk(chunk))
	return nil
}

func (c *loggingConfig) isMasked(target string) bool {
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
				return nil, fmt.Errorf("field '%s' has an invalid mask tag", fieldType.Name)
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
