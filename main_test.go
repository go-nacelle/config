package config

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Clearenv()
	os.Exit(m.Run())
}
