package logging

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/efritz/zubrin/config"
	"github.com/efritz/zubrin/internal/consts"
	"github.com/efritz/zubrin/internal/util"
	"github.com/efritz/zubrin/tags"
)

type (
	loggingConfig struct {
		config.Config
		logger Logger
	}

	// Logger is an interface to the logger where config values are printed.
	Logger interface {
		// Printf logs a mesasge. Arguments should be handled in the manner of fmt.Printf.
		Printf(format string, args ...interface{})
	}
)

// NewLoggingConfig wraps a config object with logging. After each successful load,
// the populated configuration object is serialized as fields and output at the info
// level.
func NewLoggingConfig(config config.Config, logger Logger) config.Config {
	return &loggingConfig{
		Config: config,
		logger: logger,
	}
}

func (c *loggingConfig) Load(target interface{}, modifiers ...tags.TagModifier) error {
	if err := c.Config.Load(target, modifiers...); err != nil {
		return err
	}

	m, err := dumpChunk(target)
	if err != nil {
		return fmt.Errorf("failed to serialize config (%s)", err.Error())
	}

	values := []string{}
	for key, value := range m {
		values = append(values, fmt.Sprintf("%s=%v", key, value))
	}

	sort.Strings(values)
	c.logger.Printf("Config loaded: %s", strings.Join(values, ", "))
	return nil
}

func (c *loggingConfig) MustLoad(target interface{}, modifiers ...tags.TagModifier) {
	if err := c.Load(target, modifiers...); err != nil {
		panic(err.Error())
	}
}

func dumpChunk(obj interface{}) (map[string]interface{}, error) {
	objValue, objType, err := util.GetIndirect(obj)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	for i := 0; i < objType.NumField(); i++ {
		var (
			fieldType       = objType.Field(i)
			fieldValue      = objValue.Field(i)
			maskTagValue    = fieldType.Tag.Get(consts.MaskTag)
			displayTagValue = fieldType.Tag.Get(consts.DisplayTag)
			displayName     = fieldType.Name
		)

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
