package infa_auth

import (
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2/log"
	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension/auth"
	"go.uber.org/zap"
)

var (
	errGenericError     = errors.New("errGenericError")
	errNotAuthenticated = errors.New("authentication didn't succeed")
)

type infaAuthExtension struct {
	cfg    *Config
	logger *zap.Logger
}

func newExtension(ctx context.Context, cfg *Config, logger *zap.Logger) (auth.Server, error) {
	if cfg.ValidationURL == "" {
		return nil, errGenericError
	}

	oe := &infaAuthExtension{
		cfg:    cfg,
		logger: logger,
	}
	return auth.NewServer(auth.WithServerStart(oe.start), auth.WithServerAuthenticate(oe.authenticate)), nil
}

func (e *infaAuthExtension) start(context.Context, component.Host) error {
	/*
		provider, err := getProviderForConfig(e.cfg)
		if err != nil {
			return fmt.Errorf("failed to get configuration from the auth server: %w", err)
		}
		e.provider = provider

		e.verifier = e.provider.Verifier(&oidc.Config{
			ClientID: e.cfg.Audience,
		})
	*/
	return nil
}

// authenticate checks whether the given context contains valid auth data. Successfully authenticated calls will always return a nil error and a context with the auth data.
func (e *infaAuthExtension) authenticate(ctx context.Context, headers map[string][]string) (context.Context, error) {
	log.Debug("executing authenticate")
	log.Debugf("headers ", headers)
	log.Debugf("ctx ", ctx)
	log.Debug("e.cfg ", *e.cfg)

	var h []string

	h = headers["Ids-Agent-Session-Id"]
	if h == nil {
		log.Debug("Ids-Agent-Session-Id header is null")
		h = headers["IDS-AGENT-SESSION-ID"]
		if h == nil {
			log.Debug("IDS-AGENT-SESSION-ID header is null")
			h = headers["ids-agent-session-id"]
			if h == nil {
				log.Debug("ids-agent-session-id header is null")
				return ctx, errGenericError
			}
		}
	}

	token := h[0]
	log.Debug("token is :" + token)
	if len(e.cfg.ValidationURL) == 0 {
		return ctx, errGenericError
	}

	if len(token) == 0 {
		return ctx, errGenericError
	}

	cl := client.FromContext(ctx)
	cl.Auth = &authData{
		token: token,
	}

	status := false
	status = validateToken(e.cfg.ValidationURL, token, e.cfg.Headerkey)
	if !status {
		return ctx, errNotAuthenticated
	}
	return client.NewContext(ctx, cl), nil
}

func validateToken(url string, sessionToken string, headerKey string) bool {
	// Create an HTTP client
	log.Debug("calling validateToken ")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// Create a client with the custom transport
	client := &http.Client{Transport: tr}

	//client := &http.Client{}

	// Create a GET request

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Debugf("Error creating request:", err)
		return false
	}
	//req.Header.Add("IDS-AGENT-SESSION-ID", sessionToken)
	req.Header.Add(headerKey, sessionToken)

	// Make the request
	log.Debug("calling making http request ")
	resp, err := client.Do(req)
	log.Debugf("resp ", resp)

	if err != nil {
		log.Debugf("Error making request:", err)
		return false
	}
	if resp.StatusCode != 200 {
		log.Debugf("http status is not 200:")
		return false
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debugf("Error reading response body:", err)
		return false
	}

	// Print the response
	log.Debugf(string(body))
	return true

}
