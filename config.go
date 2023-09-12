// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package infa_auth

// Config has the configuration for the OIDC Authenticator extension.
type Config struct {
	ValidationURL string `mapstructure:"validation_url"`
	Headerkey     string `mapstructure:"header_key"`
}
