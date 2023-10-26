package infa_auth

import (
	"context"
	"github.com/shadow0wolf/infa_auth/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

const (
	defaultTimeOut            = 7
	defaultClientSideSsl      = true
	defaultValidationURL      = "http://localhost:8080/"
	defaultHeaderkey          = "IDS-AGENT-SESSION-ID"
	defaultClientJksPath      = "/mnt/crt/client_cert.jks"
	defaultClientJksPassword  = "change_it_1"
	defaultCAJksPath          = "/mnt/crt/trust_store.jks"
	defaultCAJksPassword      = "change_it_2"
	defaultInsecureSkipVerify = false
)

func NewFactory() extension.Factory {
	return extension.NewFactory(
		metadata.Type,
		createDefaultConfig,
		createExtension,
		metadata.ExtensionStability,
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		TimeOut:            defaultTimeOut,
		ClientSideSsl:      defaultClientSideSsl,
		ValidationURL:      defaultValidationURL,
		Headerkey:          defaultHeaderkey,
		ClientJksPath:      defaultClientJksPath,
		ClientJksPassword:  defaultClientJksPassword,
		CAJksPath:          defaultCAJksPath,
		CAJksPassword:      defaultCAJksPassword,
		InsecureSkipVerify: defaultInsecureSkipVerify,
	}
}

func createExtension(ctx context.Context, set extension.CreateSettings, cfg component.Config) (extension.Extension, error) {
	return newExtension(ctx, cfg.(*Config), set.Logger)
}
