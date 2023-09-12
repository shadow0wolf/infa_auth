package infa_auth

import (
	"path/filepath"
	"testing"

	"github.com/shadow0wolf/infa_auth/internal/metadata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

/*
func TestLoadConfig2(t *testing.T) {

	//factories, err := componenttest.ExampleComponents()
	assert.NoError(t, err)

	//factories.Extension[typeStr] = NewFactory()

	//cfg, err := configtest.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)
	require.NoError(t, err)
	//require.NotNil(t, cfg)
}
*/

type Aaa struct {
	id          component.ID
	expected    component.Config
	expectedErr bool
}

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	var tt = Aaa{
		id: component.NewID(metadata.Type),
		expected: &Config{
			ValidationURL: "https://pod.ics.dev:444/session-service/api/v1/session/Agent",
			Headerkey:     "IDS-AGENT-SESSION-ID",
		},
	}

	sub, err := cm.Sub(tt.id.String())
	require.NoError(t, err)
	require.NoError(t, component.UnmarshalConfig(sub, cfg))
	assert.NoError(t, component.ValidateConfig(cfg))
	assert.Equal(t, tt.expected, cfg)
}
