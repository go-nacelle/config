package config

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidMaskTag = errors.New("invalid mask tag")
)

// LoadError is an error returned when a config fails to load. It contains Errors, a slice of
// any errors that occurred during config load.
type LoadError struct {
	// Errors that were found while loading a config.
	Errors []error
}

func newLoadError(errs []error) *LoadError {
	return &LoadError{Errors: errs}
}

func (le *LoadError) Error() string {
	var errs []string
	for _, errStr := range le.Errors {
		errs = append(errs, errStr.Error())
	}

	return "failed to load config: " + strings.Join(errs, ", ")
}

// SerializeError is an error returned when a config loader fails to serialize a field
// in a config struct.
type SerializeError struct {
	FieldName string
	Err       error
}

func newSerializeError(fieldName string, serializeErr error) *SerializeError {
	return &SerializeError{
		FieldName: fieldName,
		Err:       serializeErr,
	}
}

func (se *SerializeError) Error() string {
	return fmt.Sprintf("field '%s': %s", se.FieldName, se.Err)
}

func (se *SerializeError) Unwrap() error {
	return se.Err
}

// PostLoadError is an error returned when a config's PostLoad method fails.
type PostLoadError struct {
	Err error
}

func newPostLoadError(err error) *PostLoadError {
	return &PostLoadError{Err: err}
}

func (ple *PostLoadError) Error() string {
	return "post load callback failed: " + ple.Err.Error()
}

func (ple *PostLoadError) Unwrap() error {
	return ple.Err
}
