// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package infa_auth

import "go.opentelemetry.io/collector/client"

var _ client.AuthData = (*authData)(nil)

type authData struct {
	token string
}

func (a *authData) GetAttribute(name string) interface{} {
	return a.token
}

func (*authData) GetAttributeNames() []string {
	return []string{"token"}
}
