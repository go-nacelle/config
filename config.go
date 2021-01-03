package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type StructDescription struct {
	Fields []FieldDescription
}

type FieldDescription struct {
	Name      string
	Default   string
	Required  bool
	TagValues map[string]string
}

// PostLoadConfig is a marker interface for configuration objects
// which should do some post-processing after being loaded. This
// can perform additional casting (e.g. ints to time.Duration) and
// more sophisticated validation (e.g. enum or exclusive values).
type PostLoadConfig interface {
	PostLoad() error
}

// Config is a structure that can populate the exported fields of a
// struct based on the value of the field `env` tags.
type Config struct {
	sourcer    Sourcer
	logger     Logger
	maskedKeys []string
}

// NewConfig creates a config loader with the given sourcer.
func NewConfig(sourcer Sourcer, configs ...ConfigOptionsFunc) *Config {
	options := getConfigOptions(configs)

	return &Config{
		sourcer:    sourcer,
		logger:     options.logger,
		maskedKeys: options.maskedKeys,
	}
}

// Init prepares state required by the registered sourcer. This
// method should be called before calling any other method.
func (c *Config) Init() error {
	return c.sourcer.Init()
}

// Load populates a configuration object. The given tag modifiers
// are applied to the configuration object pre-load.
func (c *Config) Load(target interface{}, modifiers ...TagModifier) error {
	config, err := ApplyTagModifiers(target, modifiers...)
	if err != nil {
		return err
	}

	errors := c.load(config)
	if len(errors) == 0 {
		sourceFields, _ := getExportedFields(config)
		targetFields, _ := getExportedFields(target)

		for i := 0; i < len(sourceFields); i++ {
			targetFields[i].Field.Set(sourceFields[i].Field)
		}
	}

	if err := loadError(errors); err != nil {
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

// Call the PostLoad method of the given target if it conforms to
// the PostLoadConfig interface.
func (c *Config) PostLoad(target interface{}) error {
	if plc, ok := target.(PostLoadConfig); ok {
		return plc.PostLoad()
	}

	return nil
}

// Assets returns a list of names of assets that compose the
// underlying sourcer. This can be a list of matched files that are
// read, or a token that denotes a fixed source.
func (c *Config) Assets() []string {
	return c.sourcer.Assets()
}

// Dump returns the full content of the underlying sourcer. This
// is used by the logging package to show the content of the
// environment and config files when a value is missing or otherwise
// illegal.
func (c *Config) Dump() map[string]string {
	return c.sourcer.Dump()
}

// Describe returns a description of the struct relevant to the given
// config object. Field descriptions include the field name, the values
// of struct tags matching the configured sourcer, whether or not the
// field must be populated, and a default value (if any).
func (c *Config) Describe(target interface{}, modifiers ...TagModifier) (*StructDescription, error) {
	config, err := ApplyTagModifiers(target, modifiers...)
	if err != nil {
		return nil, err
	}

	return c.describe(config)
}

func (c *Config) load(target interface{}) []error {
	objValue, objType, err := getIndirect(target)
	if err != nil {
		return []error{err}
	}

	return c.loadStruct(objValue, objType)
}

func (c *Config) loadStruct(objValue reflect.Value, objType reflect.Type) []error {
	if objType.Kind() != reflect.Struct {
		return []error{fmt.Errorf(
			"invalid embedded type in configuration struct",
		)}
	}

	errors := []error{}
	for i := 0; i < objType.NumField(); i++ {
		field := objValue.Field(i)
		fieldType := objType.Field(i)
		defaultTagValue := fieldType.Tag.Get(DefaultTag)
		requiredTagValue := fieldType.Tag.Get(RequiredTag)

		if fieldType.Anonymous {
			errors = append(errors, c.loadStruct(field, fieldType.Type)...)
			continue
		}

		tagValues := []string{}
		for _, tag := range c.sourcer.Tags() {
			tagValues = append(tagValues, fieldType.Tag.Get(tag))
		}

		err := c.loadEnvField(
			field,
			fieldType.Name,
			tagValues,
			defaultTagValue,
			requiredTagValue,
		)

		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (c *Config) loadEnvField(
	fieldValue reflect.Value,
	name string,
	tagValues []string,
	defaultTag string,
	requiredTag string,
) error {
	val, flag, err := c.sourcer.Get(tagValues)
	if err != nil {
		return err
	}

	if flag == FlagSkip && RequiredTag == "" && DefaultTag == "" {
		return nil
	}

	if !fieldValue.IsValid() {
		return fmt.Errorf("field '%s' is invalid", name)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("field '%s' can not be set", name)
	}

	if flag == FlagFound {
		if !toJSON([]byte(val), fieldValue.Addr().Interface()) {
			return fmt.Errorf("value supplied for field '%s' cannot be coerced into the expected type", name)
		}

		return nil
	}

	if requiredTag != "" {
		val, err := strconv.ParseBool(requiredTag)
		if err != nil {
			return fmt.Errorf("field '%s' has an invalid required tag", name)
		}

		if val {
			return fmt.Errorf("no value supplied for field '%s'", name)
		}
	}

	if defaultTag != "" {
		if !toJSON([]byte(defaultTag), fieldValue.Addr().Interface()) {
			return fmt.Errorf("default value for field '%s' cannot be coerced into the expected type", name)
		}

		return nil
	}

	return nil
}

func (c *Config) describe(target interface{}) (*StructDescription, error) {
	objValue, objType, err := getIndirect(target)
	if err != nil {
		return nil, err
	}

	return c.describeStruct(objValue, objType)
}

func (c *Config) describeStruct(objValue reflect.Value, objType reflect.Type) (*StructDescription, error) {
	if objType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid embedded type in configuration struct")
	}

	var fields []FieldDescription
	for i := 0; i < objType.NumField(); i++ {
		field := objValue.Field(i)
		fieldType := objType.Field(i)
		defaultTagValue := fieldType.Tag.Get(DefaultTag)
		requiredTagValue := fieldType.Tag.Get(RequiredTag)

		if fieldType.Anonymous {
			definition, err := c.describeStruct(field, fieldType.Type)
			if err != nil {
				return nil, err
			}

			fields = append(fields, definition.Fields...)
			continue
		}

		tagValues := map[string]string{}
		for _, tag := range c.sourcer.Tags() {
			if value := fieldType.Tag.Get(tag); value != "" {
				tagValues[tag] = value
			}
		}

		if len(tagValues) != 0 {
			fields = append(fields, FieldDescription{
				Name:      fieldType.Name,
				Default:   defaultTagValue,
				Required:  requiredTagValue != "",
				TagValues: tagValues,
			})
		}
	}

	return &StructDescription{Fields: fields}, nil
}

//
// Helpers

func loadError(errors []error) error {
	if len(errors) == 0 {
		return nil
	}

	messages := []string{}
	for _, err := range errors {
		messages = append(messages, err.Error())
	}

	return fmt.Errorf("failed to load config (%s)", strings.Join(messages, ", "))
}
