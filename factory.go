// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package infa_auth

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

const (
	defaultValidationUrl = "http://localhost:8080/"
	Type                 = "infa_auth"
	ExtensionStability   = component.StabilityLevelBeta
)

// NewFactory creates a factory for the OIDC Authenticator extension.
func NewFactory() extension.Factory {
	return extension.NewFactory(
		//metadata.Type,
		Type,
		createDefaultConfig,
		createExtension,
		//metadata.ExtensionStability,
		ExtensionStability,
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		ValidationURL: defaultValidationUrl,
	}
}

func createExtension(ctx context.Context, set extension.CreateSettings, cfg component.Config) (extension.Extension, error) {
	return newExtension(ctx, cfg.(*Config), set.Logger)
}
