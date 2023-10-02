package infa_auth

import (
	"context"
	"github.com/shadow0wolf/infa_auth/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

const (
	defaultTimeOut            = 2
	defaultClientSideSsl      = false
	defaultValidationURL      = "http://localhost:8080/"
	defaultHeaderkey          = "IDS-AGENT-SESSION-ID"
	defaultClientCertPath     = "/mnt/crt/client_cert.crt"
	defaultCACertPath         = "/mnt/crt/t_store_def.crt"
	defaultInsecureSkipVerify = false
	defaultClientKeyPath      = "/mnt/crt/client_key.crt"
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
		ClientCertPath:     defaultClientCertPath,
		CACertPath:         defaultCACertPath,
		InsecureSkipVerify: defaultInsecureSkipVerify,
		ClientKeyPath:      defaultClientKeyPath,
	}
}

func createExtension(ctx context.Context, set extension.CreateSettings, cfg component.Config) (extension.Extension, error) {
	return newExtension(ctx, cfg.(*Config), set.Logger)
}
