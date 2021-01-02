package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Config is a structure that can populate the exported fields of a
// struct based on the value of the field `env` tags.
type Config interface {
	// Init prepares state required by the registered sourcer. This
	// method should be called before calling any other method.
	Init() error

	// Load populates a configuration object. The given tag modifiers
	// are applied to the configuration object pre-load.
	Load(interface{}, ...TagModifier) error

	// Call the PostLoad method of the given target if it conforms to
	// the PostLoadConfig interface.
	PostLoad(interface{}) error

	// MustInject calls Injects and panics on error.
	MustLoad(interface{}, ...TagModifier)

	// Assets returns a list of names of assets that compose the
	// underlying sourcer. This can be a list of matched files that are
	// read, or a token that denotes a fixed source.
	Assets() []string

	// Dump returns the full content of the underlying sourcer. This
	// is used by the logging package to show the content of the
	// environment and config files when a value is missing or otherwise
	// illegal.
	Dump() map[string]string

	// Describe returns a description of the struct relevant to the given
	// config object. Field descriptions include the field name, the values
	// of struct tags matching the configured sourcer, whether or not the
	// field must be populated, and a default value (if any).
	Describe(interface{}, ...TagModifier) (*StructDescription, error)
}

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

type config struct {
	sourcer Sourcer
}

var _ Config = &config{}

// NewConfig creates a config loader with the given sourcer.
func NewConfig(sourcer Sourcer) Config {
	return &config{
		sourcer: sourcer,
	}
}

func (c *config) Init() error {
	return c.sourcer.Init()
}

func (c *config) Load(target interface{}, modifiers ...TagModifier) error {
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

	return loadError(errors)
}

func (c *config) PostLoad(target interface{}) error {
	if plc, ok := target.(PostLoadConfig); ok {
		return plc.PostLoad()
	}

	return nil
}

// MustLoad calls Load and panics on error.
func (c *config) MustLoad(target interface{}, modifiers ...TagModifier) {
	if err := c.Load(target, modifiers...); err != nil {
		panic(err.Error())
	}
}

func (c *config) Assets() []string {
	return c.sourcer.Assets()
}

func (c *config) Dump() map[string]string {
	return c.sourcer.Dump()
}

func (c *config) Describe(target interface{}, modifiers ...TagModifier) (*StructDescription, error) {
	config, err := ApplyTagModifiers(target, modifiers...)
	if err != nil {
		return nil, err
	}

	return c.describe(config)
}

func (c *config) load(target interface{}) []error {
	objValue, objType, err := getIndirect(target)
	if err != nil {
		return []error{err}
	}

	return c.loadStruct(objValue, objType)
}

func (c *config) loadStruct(objValue reflect.Value, objType reflect.Type) []error {
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

func (c *config) loadEnvField(
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

func (c *config) describe(target interface{}) (*StructDescription, error) {
	objValue, objType, err := getIndirect(target)
	if err != nil {
		return nil, err
	}

	return c.describeStruct(objValue, objType)
}

func (c *config) describeStruct(objValue reflect.Value, objType reflect.Type) (*StructDescription, error) {
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
