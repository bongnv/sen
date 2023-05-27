package zap_test

import (
	"testing"

	"go.uber.org/zap"

	"github.com/bongnv/sen/pkg/sen"

	zapPlugin "github.com/bongnv/sen/pkg/plugins/zap"
)

type mockPlugin struct {
	Logger *zap.Logger `inject:"logger"`
}

func (p mockPlugin) Initialize() error {
	return nil
}

func TestPlugin(t *testing.T) {
	t.Run("should inject logger to the app", func(t *testing.T) {
		app := sen.New()
		m := &mockPlugin{}

		err := app.With(&zapPlugin.Plugin{}, m)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}

		if m.Logger == nil {
			t.Errorf("Expected *zap.Logger to be populated")
		}
	})
}
