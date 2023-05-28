package envconfig_test

import (
	"fmt"
	"testing"

	"github.com/bongnv/sen/pkg/plugins/envconfig"
	"github.com/bongnv/sen/pkg/sen"
)

type mockConfig struct {
	MockName string `env:"MOCK_NAME" envDefault:"defaultName"`
}

func TestConfig(t *testing.T) {
	t.Run("should use env to load environment variables", func(t *testing.T) {
		cfg := &mockConfig{}
		app := sen.New()
		err := app.With(envconfig.Config("mock-config", cfg))
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}

		if cfg.MockName != "defaultName" {
			t.Errorf("Expected the default value is used but got %v", cfg.MockName)
		}
	})

	t.Run("should return an error if unable to load the config", func(t *testing.T) {
		app := sen.New()
		err := app.With(envconfig.Config("mock-config", mockConfig{}))
		if fmt.Sprintf("%v", err) != "env: expected a pointer to a Struct" {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
