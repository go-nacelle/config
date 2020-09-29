package config

import (
	"fmt"
	"time"
)

type TestSimpleConfig struct {
	X string   `env:"x"`
	Y int      `env:"y"`
	Z []string `env:"w" display:"Q"`
}

type TestEmbeddedJSONConfig struct {
	P1 *TestJSONPayload `env:"p1"`
	P2 *TestJSONPayload `env:"p2"`
}

type TestJSONPayload struct {
	V1 int     `json:"v_int"`
	V2 float64 `json:"v_float"`
	V3 bool    `json:"v_bool"`
}

type TestRequiredConfig struct {
	X string `env:"x" required:"true"`
}

type TestBadRequiredConfig struct {
	X string `env:"x" required:"yup"`
}

type TestDefaultConfig struct {
	X string   `env:"x" default:"foo"`
	Y []string `env:"y" default:"[\"bar\", \"baz\", \"bonk\"]"`
}

type TestBadDefaultConfig struct {
	X int `env:"x" default:"foo"`
}

type TestUnsettableConfig struct {
	x int `env:"s"`
}

type TestPostLoadConfig struct {
	X int `env:"X"`
}

func (c *TestPostLoadConfig) PostLoad() error {
	if c.X < 0 {
		return fmt.Errorf("X must be positive")
	}

	return nil
}

type TestPostLoadConversion struct {
	RawDuration int `env:"duration"`
	Duration    time.Duration
}

func (c *TestPostLoadConversion) PostLoad() error {
	c.Duration = time.Duration(c.RawDuration) * time.Second
	return nil
}

type TestMaskConfig struct {
	X string   `env:"x"`
	Y int      `env:"y" mask:"true"`
	Z []string `env:"w" mask:"true"`
}

type TestBadMaskTagConfig struct {
	X string `env:"x" mask:"34"`
}

type TestParentConfig struct {
	ChildConfig
	X int `env:"x"`
	Y int `env:"y"`
}

type ChildConfig struct {
	A int `env:"a"`
	B int `env:"b"`
	C int `env:"c"`
}

func (c *ChildConfig) PostLoad() error {
	if c.A >= c.B || c.B >= c.C {
		return fmt.Errorf("fields must be increasing")
	}

	return nil
}

type TestBadParentConfig struct {
	*ChildConfig
	X int `env:"x"`
	Y int `env:"y"`
}

type BasicConfig struct {
	X string `env:"a" default:"q"`
	Y string
}

type BasicFileConfig struct {
	X string `file:"a" default:"q"`
	Y string
}

type BasicFlagConfig struct {
	X string `flag:"a" default:"q"`
	Y string
}

type ParentConfig struct {
	BasicConfig
}

type ParentFileConfig struct {
	BasicFileConfig
}

type ParentFlagConfig struct {
	BasicFlagConfig
}

type TestStringer struct{}

func (TestStringer) String() string {
	return "bar"
}
