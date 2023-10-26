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

type Aaa struct {
	id       component.ID
	expected component.Config
}

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	factory := NewFactory()
	defaultCfg := factory.CreateDefaultConfig()

	var tt = Aaa{
		id: component.NewID(metadata.Type),
		expected: &Config{
			TimeOut:            1,
			ClientSideSsl:      true,
			ValidationURL:      "https://pod.ics.dev:444/session-service/api/v1/session/Agent",
			Headerkey:          "IDS-AGENT-SESSION-ID",
			ClientJksPath:      "/mnt/a/c1Client.crt",
			ClientJksPassword:  "changeit1",
			CAJksPath:          "/mnt/a/ca_cert_path.crt",
			CAJksPassword:      "changeit2",
			InsecureSkipVerify: false,
		},
	}

	sub, err := cm.Sub(tt.id.String())
	require.NoError(t, err)
	require.NoError(t, component.UnmarshalConfig(sub, defaultCfg))
	assert.NoError(t, component.ValidateConfig(defaultCfg))
	assert.Equal(t, tt.expected, defaultCfg)
}
