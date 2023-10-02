package infa_auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/extension/extensiontest"
)

func TestCreateDefaultConfig(t *testing.T) {
	// prepare and test
	expected := &Config{
		TimeOut:            2,
		ClientSideSsl:      false,
		ValidationURL:      "http://localhost:8080/",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		ClientCertPath:     "/mnt/crt/client_cert.crt",
		CACertPath:         "/mnt/crt/t_store_def.crt",
		InsecureSkipVerify: false,
		ClientKeyPath:      "/mnt/crt/client_key.crt",
	}

	// test
	cfg := createDefaultConfig()

	// verify
	assert.Equal(t, expected, cfg)
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestCreateExtension(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	ext, err := createExtension(context.Background(), extensiontest.NewNopCreateSettings(), cfg)
	assert.NoError(t, err)
	assert.NotNil(t, ext)
}

func TestNewFactory(t *testing.T) {
	f := NewFactory()
	assert.NotNil(t, f)
}
