package config

import (
	"context"
	"fmt"
)

type configKeyType struct{}

var configKey = configKeyType{}

func WithConfig(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

func FromContext(ctx context.Context) *Config {
	if v, ok := ctx.Value(configKey).(*Config); ok {
		return v
	}
	return nil
}

func LoadFromContext(ctx context.Context, target interface{}, modifiers ...TagModifier) error {
	cfg := FromContext(ctx)
	if cfg == nil {
		return fmt.Errorf("config loader not found in context")
	}

	return cfg.Load(target, modifiers...)
}
