# Nacelle Config [![GoDoc](https://godoc.org/github.com/go-nacelle/config?status.svg)](https://godoc.org/github.com/go-nacelle/config) [![CircleCI](https://circleci.com/gh/go-nacelle/config.svg?style=svg)](https://circleci.com/gh/go-nacelle/config) [![Coverage Status](https://coveralls.io/repos/github/go-nacelle/config/badge.svg?branch=master)](https://coveralls.io/github/go-nacelle/config?branch=master)

Configuration loading and validation for [nacelle](https://nacelle.dev).

---

Often, [initializers and processes](https://nacelle.dev/docs/core/process) will need external configuration during their startup process. These values can be pulled from a **configuration loader** backed by a particular [source](#sourcers) (e.g. environment or disk) and assigned to tagged fields of a configuration struct.

You can see an additional example of loading configuration in the [example repository](https://github.com/go-nacelle/example): [definition](https://github.com/go-nacelle/example/blob/843979aaa86786784a1ca3646e8d0d1f69e29c65/internal/redis_initializer.go#L13) and [loading](https://github.com/go-nacelle/example/blob/843979aaa86786784a1ca3646e8d0d1f69e29c65/internal/redis_initializer.go#L36).

### Configuration Struct Definition

Configuration structs are defined by the application or library developer with the fields needed by the package in which they are defined. Each field is tagged with a *source hint* (e.g. an environment variable name, a key in a YAML file) and, optionally, default values and basic validation. Tagged fields must be exported in order for this package to assign to them.

The following example defines configuration for a hypothetical worker process. For the application to start successfully, the address of an API must be supplied. All other configuration values are optional.

```go
type Config struct {
    APIAddr        string   `env:"api_addr" required:"true"`
    CassandraHosts []string `env:"cassandra_hosts"`
    NumWorkers     int      `env:"num_workers" default:"10"`
    BufferSize     int      `env:"buffer_size" default:"1024"`
}
```

### Configuration Loading

At initialization time of an application component, the particular subset of configuration variables should be populated, validated, and stored on the service that will later require them.

```go
func (p *Process) Init(config nacelle.Config) error {
    appConfig := &Config{}
    if err := config.Load(appConfig); err != nil {
        return err
    }

    // Use populated appConfig
    return nil
}
```

The `Load` method fails if a value from the source cannot be converted into the correct type, a value from the source cannot be decoded as JSON (if the target is a non-string type), or is required and not supplied. After each successful load of a configuration struct, the loaded configuration values will are logged. This, however, may be a concern for application secrets. In order to hide sensitive configuration values, add the `mask:"true"` struct tag to the field. This will omit that value from the log message. Additionally, configuration loader object can be initialized with a blacklist of values that should be masked (values printed as `*****` rather than their real value) instead of omitted. These values can be configured in the [bootstrapper](https://nacelle.dev/docs/core).

#### Conversion and Validation

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

#### Anonymous Structs

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

#### Sourcers

A sourcer reads values from a particular source based on a configuration struct's tags. Sourcers declare the struct tags that determine their behavior when loading configuration structs. The examples above only work with the environment sourcer. All sources support the `default` and `required` tags (which are mutually exclusive). Tagged fields must be exported. The following six sourcers are supplied. Additional behavior can be added by conforming to the *Sourcer* interface.

<dl>
  <dt>Environment Sourcer</dt>
  <dd>An <a href="https://godoc.org/github.com/go-nacelle/config#NewEnvSourcer">environment sourcer</a> reads the <code>env</code> tag and looks up the corresponding value in the process's environment. An expected prefix may be supplied in order to namespace application configuration from the rest of the system. A sourcer instantiated with <code>NewEnvSourcer("APP")</code> will load the env tag <code>fetch_limit</code> from the environment variable <code>APP_FETCH_LIMIT</code> and falling back to the environment variable <code>FETCH_LIMIT</code>.</dd>

  <dt>Test Environment Sourcer</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewTestEnvSourcer">test environment sourcer</a> reads the <code>env</code> tag but looks up the corresponding value from a literal map. This sourcer can be used in unit tests where the full construction of a nacelle [process](https://nacelle.dev/docs/core/process) is too burdensome.</dd>

  <dt>Flag Sourcer</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewFlagSourcer">flag sourcer</a> reads the <code>flag</code> tag and looks up the corresponding value attached to the process's command line arguments.</dd>

  <dt>File Sourcer</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewFileSourcer">file sourcer</a> reads the <code>file</code> tag and returns the value at the given path. A filename and a file parser musts be supplied on instantiation. Both <a href="https://godoc.org/github.com/go-nacelle/config#ParseYAML">ParseYAML</a> and <a href="https://godoc.org/github.com/go-nacelle/config#ParseTOML">ParseTOML</a> are supplied file parsers -- note that as JSON is a subset of YAML, <code>ParseYAML</code> will also correctly parse JSON files. If a <code>nil</code> file parser is supplied, one is chosen by the filename extension. A file sourcer will load the file tag <code>api.timeout</code> from the given file by parsing it into a map of values and recursively walking the (keys separated by dots). This can return a primitive type or a structured map, as long as the target field has a compatible type. The constructor <a href="https://godoc.org/github.com/go-nacelle/config#NewOptionalFileSourcer">NewOptionalFileSourcer</a> will return a no-op sourcer if the filename does not exist.</dd>

  <dt>Multi sourcer</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewMultiSourcer">multi-sourcer</a> is a sourcer wrapping one or more other sourcers. For each configuration struct field, each sourcer is queried in reverse order of registration and the first value to exist is returned. This is useful to allow a chain of configuration files in which some files or directories take precedence over others, or to allow environment variables to take precedence over files.</dd>

  <dt>Directory Sourcer</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewDirectorySourcer">directory sourcer</a> creates a multi-sourcer by reading each file in a directory in alphabetical order. The constructor <a href="https://godoc.org/github.com/go-nacelle/config#NewOptionalDirectorySourcer">NewOptionalDirectorySourcer</a> will return a no-op sourcer if the directory does not exist.</dd>

  <dt>Glob Sourcer</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewGlobSourcer">glob sourcer</a> creates a multi-sourcer by reading each file that matches a given glob pattern. Each matching file creates a distinct file sourcer and does so in alphabetical order.</dd>
</dl>

### Tag Modifiers

A tag modifier dynamically alters the tags of a configuration struct. The following five tag modifiers are supplied. Additional behavior can be added by conforming to the *TagModifier* interface.

<dl>
  <dt>Default Tag Setter</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewDefaultTagSetter">default tag setter</a> sets the <code>default</code> tag for a particular field. This is useful when the default values supplied by a library are inappropriate for a particular application. This would otherwise require a source change in the library.</dd>

  <dt>Display Tag Setter</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewDisplayTagSetter">display tag setter</a> sets the <code>display</code> tag to the value of the <code>env</code> tag. This tag modifier can be used to provide sane defaults to the tag without doubling the length of the struct tag definition.</dd>

  <dt>Flag Tag Setter</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewFlagTagSetter">flag tag setter</a> sets the <code>flag</code> tag to the value of the <code>env</code> tag. This tag modifier can be used to provide sane defaults to the tag without doubling the length of the struct tag definition.</dd>

  <dt>File Tag Setter</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewFileTagSetter">file tag setter</a> sets the <code>file</code> tag to the value of the <code>env</code> tag. This tag modifier can be used to provide sane defaults to the tag without doubling the length of the struct tag definition.</dd>

  <dt>Env Tag Prefixer</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewEnvTagPrefixer">environment tag prefixer</a> inserts a prefix on each <code>env</code> tags. This is useful when two distinct instances of the same configuration are required, and each one should be configured independently from the other (for example, using the same abstraction to consume from two different event busses with the same consumer code).</dd>

  <dt>Flag Tag Prefixer</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewFlagTagPrefixer">flag tag prefixer</a> inserts a prefix on each <code>flag</code> tag. This effectively looks in a distinct top-level namespace in the parsed configuration. This is similar to the env tag prefixer.</dd>

  <dt>File Tag Prefixer</dt>
  <dd>A <a href="https://godoc.org/github.com/go-nacelle/config#NewFileTagPrefixer">file tag prefixer</a> inserts a prefix on each <code>file</code> tag. This effectively looks in a distinct top-level namespace in the parsed configuration. This is similar to the env tag prefixer.</dd>
</dl>

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
