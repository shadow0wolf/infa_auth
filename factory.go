// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package infa_auth

import (
	"context"

	"github.com/shadow0wolf/infa_auth/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

const (
	defaultValidationUrl = "http://localhost:8080/"
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
		ValidationURL: defaultValidationUrl,
	}
}

func createExtension(ctx context.Context, set extension.CreateSettings, cfg component.Config) (extension.Extension, error) {
	return newExtension(ctx, cfg.(*Config), set.Logger)
}
