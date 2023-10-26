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
		TimeOut:            7,
		ClientSideSsl:      true,
		ValidationURL:      "http://localhost:8080/",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		ClientJksPath:      "/mnt/crt/client_cert.jks",
		ClientJksPassword:  "change_it_1",
		CAJksPath:          "/mnt/crt/trust_store.jks",
		CAJksPassword:      "change_it_2",
		InsecureSkipVerify: false,
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
