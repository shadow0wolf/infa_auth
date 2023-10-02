package infa_auth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"log"
	"testing"
)

func TestStartPass(t *testing.T) {
	config := &Config{
		TimeOut:            2,
		ClientSideSsl:      false,
		ValidationURL:      "http://localhost:8080/",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		ClientCertPath:     "/mnt/crt/client_cert.crt",
		CACertPath:         "/mnt/crt/t_store_def.crt",
		InsecureSkipVerify: false,
		ClientKeyPath:      "/mnt/crt/client_key.crt",
	}
	extension, _ := newExtension(context.Background(), config, zap.NewNop())
	err := extension.Start(context.Background(), nil)
	assert.NoError(t, err)

}

func TestStartFail1(t *testing.T) {
	config := &Config{
		TimeOut:       2,
		ClientSideSsl: false,
		//ValidationURL:      "http://localhost:8080/",
		Headerkey:          "IDS-AGENT-SESSION-ID",
		ClientCertPath:     "/mnt/crt/client_cert.crt",
		CACertPath:         "/mnt/crt/t_store_def.crt",
		InsecureSkipVerify: false,
		ClientKeyPath:      "/mnt/crt/client_key.crt",
	}
	extension, _ := newExtension(context.Background(), config, zap.NewNop())
	err := extension.Start(context.Background(), nil)
	assert.ErrorContains(t, err, "ValidationURL is empty")

}

func TestStartFail2(t *testing.T) {
	config := &Config{
		TimeOut:       2,
		ClientSideSsl: false,
		ValidationURL: "http://localhost:8080/",
		//Headerkey:          "IDS-AGENT-SESSION-ID",
		ClientCertPath:     "/mnt/crt/client_cert.crt",
		CACertPath:         "/mnt/crt/t_store_def.crt",
		InsecureSkipVerify: false,
		ClientKeyPath:      "/mnt/crt/client_key.crt",
	}
	extension, _ := newExtension(context.Background(), config, zap.NewNop())
	err := extension.Start(context.Background(), nil)
	assert.ErrorContains(t, err, "Headerkey is empty")
}

func TestStartFail3(t *testing.T) {
	config := &Config{
		TimeOut:       2,
		ClientSideSsl: true,
		ValidationURL: "http://localhost:8080/",
		Headerkey:     "IDS-AGENT-SESSION-ID",
		//ClientCertPath:     "/mnt/crt/client_cert.crt",
		CACertPath:         "/mnt/crt/t_store_def.crt",
		InsecureSkipVerify: false,
		ClientKeyPath:      "/mnt/crt/client_key.crt",
	}
	extension, _ := newExtension(context.Background(), config, zap.NewNop())
	err := extension.Start(context.Background(), nil)
	assert.ErrorContains(t, err, "ClientCertPath is empty")
}

func TestAuthenticationSucceeded(t *testing.T) {

	config := &Config{
		TimeOut:       2,
		ClientSideSsl: false,
		ValidationURL: "http://localhost:9898/session-service/api/v1/session/Agent",
		Headerkey:     "IDS-AGENT-SESSION-ID",
		//ClientCertPath:     "/mnt/crt/client_cert.crt",
		//CACertPath:         "/mnt/crt/t_store_def.crt",
		InsecureSkipVerify: true,
		//ClientKeyPath:      "/mnt/crt/client_key.crt",
	}

	p, err := newExtension(context.Background(), config, zap.NewNop())
	log.Println(p)
	require.NoError(t, err)
	ctx, err1 := p.Authenticate(context.Background(), map[string][]string{config.Headerkey: {"123123123"}})
	log.Println(ctx)
	require.NoError(t, err1)

}

func TestAuthenticationFailed(t *testing.T) {

	config := &Config{
		TimeOut:       2,
		ClientSideSsl: false,
		ValidationURL: "http://localhost:9898/session-service/api/v1/session/Agent",
		Headerkey:     "IDS-AGENT-SESSION-ID",
		//ClientCertPath:     "/mnt/crt/client_cert.crt",
		//CACertPath:         "/mnt/crt/t_store_def.crt",
		InsecureSkipVerify: true,
		//ClientKeyPath:      "/mnt/crt/client_key.crt",
	}

	p, err := newExtension(context.Background(), config, zap.NewNop())
	log.Println(p)
	require.NoError(t, err)
	ctx, err1 := p.Authenticate(context.Background(), map[string][]string{config.Headerkey: {"123123123_"}})
	log.Println(ctx)
	require.Error(t, err1)

}
