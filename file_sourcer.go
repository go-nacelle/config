package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"
)

type fileSourcer struct {
	filename string
	parser   FileParser
	fs       FileSystem
	optional bool
	values   map[string]string
}

var _ Sourcer = &fileSourcer{}

type FileParser func(content []byte) (map[string]interface{}, error)

var parserMap = map[string]FileParser{
	".yaml": ParseYAML,
	".yml":  ParseYAML,
	".json": ParseYAML,
	".toml": ParseTOML,
}

// NewOptionalFileSourcer creates a file sourcer that does not error on init
// if the file does not exist. The underlying sourcer returns no values.
func NewOptionalFileSourcer(filename string, parser FileParser, configs ...FileSourcerConfigFunc) Sourcer {
	return &fileSourcer{
		filename: filename,
		parser:   parser,
		fs:       getFileSourcerConfigOptions(configs).fs,
		optional: true,
	}
}

// NewFileSourcer creates a sourcer that reads content from a file. The format
// of the file is read by the given FileParser. The content of the file must be
// an encoding of a map from string keys to JSON-serializable values. If a nil
// parser is supplied, one will be selected based on the extension of the file.
// JSON, YAML, and TOML files are supported.
func NewFileSourcer(filename string, parser FileParser, configs ...FileSourcerConfigFunc) Sourcer {
	return &fileSourcer{
		filename: filename,
		parser:   parser,
		fs:       getFileSourcerConfigOptions(configs).fs,
	}
}

func (s *fileSourcer) Init() error {
	if s.optional {
		exists, err := s.fs.Exists(s.filename)
		if err != nil || !exists {
			return err
		}

		if !exists {
			return nil
		}
	}

	values, err := readFile(s.filename, s.fs, s.parser)
	if err != nil {
		return err
	}

	jsonValues, err := serializeJSONValues(values)
	if err != nil {
		return err
	}

	s.values = jsonValues
	return nil
}

func (s *fileSourcer) Tags() []string {
	return []string{FileTag}
}

func (s *fileSourcer) Get(values []string) (string, SourcerFlag, error) {
	if values[0] == "" {
		return "", FlagSkip, nil
	}

	segments := strings.Split(values[0], ".")

	if val, ok := s.values[segments[0]]; ok {
		if val, ok := extractJSONPath(val, segments[1:]); ok {
			return val, FlagFound, nil
		}
	}

	return "", FlagMissing, nil
}

func (s *fileSourcer) Assets() []string {
	return []string{s.filename}
}

func (s *fileSourcer) Dump() map[string]string {
	return s.values
}

//
// Parsers

// NewYAMLFileSourcer creates a file sourcer that parses conent as YAML.
func NewYAMLFileSourcer(filename string) Sourcer {
	return NewFileSourcer(filename, ParseYAML)
}

// NewTOMLFileSourcer creates a file sourcer that parses conent as TOML.
func NewTOMLFileSourcer(filename string) Sourcer {
	return NewFileSourcer(filename, ParseTOML)
}

// ParseYAML parses the given content as YAML.
func ParseYAML(content []byte) (map[string]interface{}, error) {
	return commonParser(content, func(content []byte, values interface{}) error {
		return yaml.Unmarshal(content, values)
	})
}

// ParseTOML parses the given content as JSON.
func ParseTOML(content []byte) (map[string]interface{}, error) {
	return commonParser(content, toml.Unmarshal)
}

func commonParser(content []byte, unmarshaller func([]byte, interface{}) error) (map[string]interface{}, error) {
	values := map[string]interface{}{}
	if err := unmarshaller(content, &values); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config (%s)", err.Error())
	}

	return values, nil
}

//
// Helpers

func readFile(filename string, fs FileSystem, parser FileParser) (map[string]interface{}, error) {
	content, err := fs.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s' (%s)", filename, err.Error())
	}

	parser, err = chooseParser(filename, parser)
	if err != nil {
		return nil, err
	}

	return parser(content)
}

func chooseParser(filename string, parser FileParser) (FileParser, error) {
	if parser != nil {
		return parser, nil
	}

	if parser, ok := parserMap[filepath.Ext(filename)]; ok {
		return parser, nil
	}

	return nil, fmt.Errorf("failed to determine parser for file %s", filename)
}

func serializeJSONValues(values map[string]interface{}) (map[string]string, error) {
	jsonValues := map[string]string{}
	for key, value := range values {
		serialized, err := serializeJSONValue(value)
		if err != nil {
			return nil, fmt.Errorf("illegal configuration value for '%s' (%s)", key, err.Error())
		}

		jsonValues[key] = serialized
	}

	return jsonValues, nil
}

func serializeJSONValue(value interface{}) (string, error) {
	if str, ok := value.(string); ok {
		return str, nil
	}

	serialized, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(serialized), nil
}

func extractJSONPath(val string, path []string) (string, bool) {
	if len(path) == 0 {
		return val, true
	}

	for _, segment := range path {
		mapping := map[string]json.RawMessage{}
		if err := json.Unmarshal([]byte(val), &mapping); err != nil {
			return "", false
		}

		inner, ok := mapping[segment]
		if !ok {
			return "", false
		}

		val = string(inner)
	}

	return val, true
}
