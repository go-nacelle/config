# Nacelle Config [![GoDoc](https://godoc.org/github.com/go-nacelle/config?status.svg)](https://godoc.org/github.com/go-nacelle/config) [![CircleCI](https://circleci.com/gh/go-nacelle/config.svg?style=svg)](https://circleci.com/gh/go-nacelle/config) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/config/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/config?branch=master)

Configuration loading and validation for [nacelle](https://github.com/go-nacelle/nacelle).

---

Nacelle assigns values to the fields of **configuration structs** by reading a particular source (e.g. environment or disk) and correlating them with struct tags. Basic validation (types and required values) are included, as well as an extension point to allow arbitrary validation in-code.

### Usage

Configuration is loaded into the application by assigning the fields of tagged structs. Tagged fields must be exported in order for this package to assign to them. The following example defines configuration for a hypothetical worker process. The address of an API must be supplied, but all other configuration values are optional.

```go
type Config struct {
    APIAddr        string   `env:"api_addr" required:"true"`
    CassandraHosts []string `env:"cassandra_hosts"`
    NumWorkers     int      `env:"num_workers" default:"10"`
    BufferSize     int      `env:"buffer_size" default:"1024"`
}
```

A configuration loader object is created with a **sourcer**, which knows how to acquire values based on struct tags. In the following, we create a sourcer that reads environment variables, then inject the values from the environment into an instance of `Config`. Additional [sourcers](#Sourcers) are supplied.

```go
config := NewConfig(NewEnvSourcer("APP"))

appConfig := &Config{}
if err := config.Load(appConfig); err != nil {
    // handle error
}
```

The `Load` method fails if a value from the source cannot be converted into the correct type, a value from the source cannot be decoded as JSON (if the target is a non-string type), or is required and not supplied.

All sources support the `default` and `required` tags (which are mutually exclusive). Tagged fields must be exported.

### Post Load Hook

After successful loading of a configuration struct, the method named `PostLoad` will be called if it is defined. This allows a place for additional validation (such as mutually exclusive settings, regex value matching, etc) and deserialization of more complex types (such enums from strings, durations from integers, etc). The following example parses and stores a `text/template` from a user-supplied string.

```go
import "text/template"

type Config struct {
    RawTemplate    string `env:"template" default:"Hello, {{.Name}}!"`
    ParsedTemplate *template.Template
}

func (c *Config) PostLoad() (err error) {
    c.ParsedTemplate, err = template.New("ConfigTemplate").Parse(c.RawTemplate)
    return
}
```

An error returned by `PostLoad` will be returned via the `Load` method.

### Anonymous Structs

Loading configuration values also works with structs containing composite fields. The following example shows the definition of multiple configuration structs with a set of shared fields.

```go
type StreamConfig struct {
    StreamName string `env:"stream_name" required:"true"`
}

type StreamProducerConfig struct {
    StreamConfig
    PublishAttempts int `env:"publish_attempts" default:"3"`
    PublishDelay    int `env:"publish_delay" default:"1"`
}

type StreamConsumerConfig struct {
    StreamConfig
    FetchLimit int `env:"fetch_limit" default:"100"`
}
```

### Logging Config

A `LoggingConfig` wraps a configuration object as well as a nacelle [logger](https://nacelle.dev/docs/core/log). After each successful load of a configuration struct, the loaded configuration values will be logged. This, however, may be a concern for application secrets. In order to hide sensitive configuration values, add the `mask:"true"` struct tag to the field. This will omit that value from the log message. Additionally, the logging config keeps a blacklist of values which should be masked (values printed as `*****` rather than their real value) instead of omitted. This blacklist iss given at the time of construction.

### Sourcers

A sourcer reads values from a particular source based on a configuration struct's tags. Sourcers declare the struct tags that determine their behavior when loading configuration structs. The examples above only work with the environment sourcer. The following six sourcers are supplied. Additional behavior can be added by conforming to the *Sourcer* interface.

**Environment Sourcer** reads the `env` tag and looks up the corresponding value in the process's environment. An expected prefix may be supplied in order to namespace application configuration from the rest of the system. A sourcer instantiated with `NewEnvSourcer("APP")` will load the env tag `fetch_limit` from the environment variable `APP_FETCH_LIMIT` and falling back to the environment variable `FETCH_LIMT`.

**Test Environment Sourcer** reads the `env` tag but looks up the corresponding value from a literal map. This sourcer is meant to be used in unit tests where the full construction of a nacelle [process](https://nacelle.dev/docs/core/process) is beneficial.

**File Sourcer** reads the `file` tag and returns the value at the given path. A filename and a file parser musts be supplied on instantiation. Both `ParseYAML` and `ParseTOML` are supplied file parsers -- note that as JSON is a subset of YAML, `ParseYAML` will also correctly parse JSON files. If a `nil` file parser is supplied, one is chosen by the filename extension.

A file sourcer will load the file tag `api.timeout` from the given file by parsing it into a map of values and recursively walking the (keys separated by dots). This can return a primitive type or a structured map, as long as the target field has a compatible type.

The constructor `NewOptionalFileSourcer` will return a no-op sourcer if the filename does not exist.

**Directory Sourcer** creates a multi-sourcer by reading each file in a directory in alphabetical order. The constructor `NewOptionalDirectorySourcer` will return a no-op sourcer if the directory does not exist.

**Glob Sourcer** creates a multi-sourcer by reading each file that matches a given glob pattern. Each matching file creates a distinct file sourcer and does so in alphabetical order.

**Multi sourcer** is a sourcer wrapping one or more other sourcers. For each configuration struct field, each sourcer is queried in reverse order of registration and the first value to exist is returned.

### Tag Modifiers

A tag modifier dynamically alters the tags of a configuration struct. The following five tag modifiers are supplied. Additional behavior can be added by conforming to the *TagModifier* interface.

**Default Tag Setter** sets the `default` tag for a particular field. This is useful when the default values supplied by a library are inappropriate for a particular application. This would otherwise require a source change in the library.

**Display Tag Setter** sets the `display` tag to the value of the `env` tag. This tag modifier can be used to provide sane defaults to the tag without doubling the length of the struct tag definition.

**File Tag Setter** sets the `file` tag to the value of the `env` tag. This tag modifier can be used to provide sane defaults to the tag without doubling the length of the struct tag definition.

**Env Tag Prefixer** inserts a prefix on each `env` tags. This is useful when two distinct instances of the same configuration are required, and each one should be configured independently from the other (for example, using the same abstraction to consume from two different event busses with the same consumer code).

**File Tag Prefixer** inserts a prefix on each `file` tag. This effectively looks in a distinct top-level namespace in the parsed configuration. This is similar to the env tag prefixer.

Tag modifiers are supplied at the time that a configuration struct is loaded. In the following example, each env tag is prefixed with `ACME_`, and the CassandraHosts field is given a default. Notice that you supply the *field* name to the tag modifier (not a tag value) when targeting a particular field value.

```go
if err := config.Load(
    appConfig,
    NewEnvTagPrefixer("ACME"),
    NewDefaultTagSetter("CassandraHosts", "[127.0.0.1:9042]"),
); err != nil {
    // handle error
}
```
