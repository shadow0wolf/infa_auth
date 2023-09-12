package infa_auth

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAuthenticationSucceeded(t *testing.T) {

	// prepare
	//startNewMockSessionService()

	config := &Config{
		ValidationURL: "http://127.0.0.1:9898/session-service/api/v1/session/Agent",
		Headerkey:     "IDS-AGENT-SESSION-ID",
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
		ValidationURL: "http://127.0.0.1:9898/session-service/api/v1/session/Agent",
		Headerkey:     "IDS-AGENT-SESSION-ID",
	}

	p, err := newExtension(context.Background(), config, zap.NewNop())
	log.Println(p)
	require.NoError(t, err)
	ctx, err1 := p.Authenticate(context.Background(), map[string][]string{config.Headerkey: {"123123123_"}})
	log.Println(ctx)
	require.Error(t, err1)

}
