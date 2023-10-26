package infa_auth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestStartPass(t *testing.T) {
	config := &Config{
		TimeOut:            5,
		ValidationURL:      "http://localhost:8080/",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		InsecureSkipVerify: true,
	}
	extension, _ := newExtension(context.Background(), config, zap.NewNop())
	err := extension.Start(context.Background(), nil)
	assert.NoError(t, err)

	config = &Config{
		TimeOut:            5,
		ClientSideSsl:      true,
		ValidationURL:      "url1",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		ClientJksPath:      "testdata/scheduler-service-keystore.jks",
		ClientJksPassword:  "changeit",
		CAJksPath:          "testdata/truststore.jks",
		CAJksPassword:      "changeit",
		InsecureSkipVerify: false,
	}

	extension, _ = newExtension(context.Background(), config, zap.NewNop())
	err = extension.Start(context.Background(), nil)
	assert.NoError(t, err)

	config = &Config{
		TimeOut:            5,
		ClientSideSsl:      false,
		ValidationURL:      "url1",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		CAJksPath:          "testdata/truststore.jks",
		CAJksPassword:      "changeit",
		InsecureSkipVerify: false,
	}

	extension, _ = newExtension(context.Background(), config, zap.NewNop())
	err = extension.Start(context.Background(), nil)
	assert.NoError(t, err)

}

func TestStartFail1(t *testing.T) {
	config := &Config{
		TimeOut:            5,
		ClientSideSsl:      true,
		ValidationURL:      "  ",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		ClientJksPath:      "/mnt/crt/client_cert.crt",
		ClientJksPassword:  "change_it_1",
		CAJksPath:          "/mnt/crt/t_store_def.crt",
		CAJksPassword:      "change_it_2",
		InsecureSkipVerify: true,
	}
	extension, _ := newExtension(context.Background(), config, zap.NewNop())
	err := extension.Start(context.Background(), nil)
	assert.ErrorContains(t, err, "ValidationURL is empty")

	config = &Config{
		TimeOut:            5,
		ClientSideSsl:      true,
		ValidationURL:      "url1",
		ClientJksPath:      "/mnt/crt/client_cert.crt",
		ClientJksPassword:  "change_it_1",
		CAJksPath:          "/mnt/crt/t_store_def.crt",
		CAJksPassword:      "change_it_2",
		InsecureSkipVerify: true,
	}
	extension, _ = newExtension(context.Background(), config, zap.NewNop())
	err = extension.Start(context.Background(), nil)
	assert.ErrorContains(t, err, "Headerkey is empty")
}

func TestStartFail2(t *testing.T) {

	config := &Config{
		TimeOut:            5,
		ClientSideSsl:      true,
		ValidationURL:      "url1",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		CAJksPath:          "testdata/truststore.jks",
		CAJksPassword:      "changeit",
		InsecureSkipVerify: false,
	}

	extension, _ := newExtension(context.Background(), config, zap.NewNop())
	err := extension.Start(context.Background(), nil)
	assert.ErrorContains(t, err, "The system cannot find the file specified")

	config = &Config{
		TimeOut:            5,
		ClientSideSsl:      false,
		ValidationURL:      "url1",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		CAJksPath:          "testdata/truststore222.jks",
		CAJksPassword:      "changeit",
		InsecureSkipVerify: false,
	}

	extension, _ = newExtension(context.Background(), config, zap.NewNop())
	err = extension.Start(context.Background(), nil)
	assert.ErrorContains(t, err, "The system cannot find the file specified")
}

/*
func TestStartFail3(t *testing.T) {
	config := &Config{
		TimeOut:            5,
		ClientSideSsl:      true,
		ValidationURL:      "http://localhost:8080/",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		ClientJksPath:      "/mnt/crt/client_cert.crt",
		ClientJksPassword:  "change_it_1",
		CAJksPath:          "/mnt/crt/t_store_def.crt",
		CAJksPassword:      "change_it_2",
		InsecureSkipVerify: true,
	}
	extension, _ := newExtension(context.Background(), config, zap.NewNop())
	err := extension.Start(context.Background(), nil)
	assert.ErrorContains(t, err, "ClientCertPath is empty")
}
*/

func TestAuthenticationSucceeded1(t *testing.T) {
	config := &Config{
		TimeOut:            5,
		ClientSideSsl:      false,
		ValidationURL:      "http://127.0.0.1:9898/session-service/api/v1/session/Agent",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		ClientJksPath:      "testdata/scheduler-service-keystore.jks",
		ClientJksPassword:  "changeit",
		CAJksPath:          "testdata/truststore.jks",
		CAJksPassword:      "changeit",
		InsecureSkipVerify: true,
	}

	extension, err := newExtension(context.Background(), config, zap.NewNop())
	require.NoError(t, err)
	err = extension.Start(context.Background(), nil)
	assert.NoError(t, err)
	_, err = extension.Authenticate(context.Background(), map[string][]string{config.Headerkey: {"123123123"}})
	require.NoError(t, err)

}

func TestAuthenticationFailed1(t *testing.T) {

	config := &Config{
		TimeOut:            5,
		ClientSideSsl:      false,
		ValidationURL:      "http://127.0.0.1:9898/session-service/api/v1/session/Agent",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		ClientJksPath:      "testdata/scheduler-service-keystore.jks",
		ClientJksPassword:  "changeit",
		CAJksPath:          "testdata/truststore.jks",
		CAJksPassword:      "changeit",
		InsecureSkipVerify: true,
	}

	extension, err := newExtension(context.Background(), config, zap.NewNop())
	require.NoError(t, err)
	err = extension.Start(context.Background(), nil)
	assert.NoError(t, err)
	_, err = extension.Authenticate(context.Background(), map[string][]string{config.Headerkey: {"123123123_"}})
	require.Error(t, err)

}

func TestAuthenticationFailedTimeOut(t *testing.T) {

	config := &Config{
		TimeOut:            5,
		ClientSideSsl:      false,
		ValidationURL:      "http://127.0.0.1:9898/session-service/api/v1/session/Agent",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		ClientJksPath:      "testdata/scheduler-service-keystore.jks",
		ClientJksPassword:  "changeit",
		CAJksPath:          "testdata/truststore.jks",
		CAJksPassword:      "changeit",
		InsecureSkipVerify: true,
	}

	extension, err := newExtension(context.Background(), config, zap.NewNop())
	require.NoError(t, err)
	err = extension.Start(context.Background(), nil)
	assert.NoError(t, err)
	_, err = extension.Authenticate(context.Background(), map[string][]string{config.Headerkey: {"999"}})
	require.Error(t, err)

}
