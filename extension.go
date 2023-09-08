package infa_auth

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension/auth"
	"go.uber.org/zap"
)

var (
	errNoAudienceProvided = errors.New("no Audience provided for the OIDC configuration")
	errNotAuthenticated   = errors.New("authentication didn't succeed")
)

type infaAuthExtension struct {
	cfg    *Config
	logger *zap.Logger
}

func newExtension(ctx context.Context, cfg *Config, logger *zap.Logger) (auth.Server, error) {
	if cfg.ValidationURL == "" {
		return nil, errNoAudienceProvided
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
	metadata := client.NewMetadata(headers)
	validationURL := metadata.Get(e.cfg.ValidationURL)
	//@#@#123 read this from request , or context
	token := "64Vjmeewe81iwbIPfgUmqu"
	if len(validationURL) == 0 {
		return ctx, errNotAuthenticated
	}

	cl := client.FromContext(ctx)
	cl.Auth = &authData{
		token: token,
	}

	status := false
	status = validateToken("", "")
	if status == false {

	}
	return client.NewContext(ctx, cl), nil
}

func validateToken(url string, sessionToken string) bool {
	// Create an HTTP client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// Create a client with the custom transport
	client := &http.Client{Transport: tr}

	//client := &http.Client{}

	// Create a GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}
	req.Header.Add("IDS-AGENT-SESSION-ID", sessionToken)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return false
	}
	if resp.StatusCode != 200 {
		fmt.Println("http status is not 200:")
		return false
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	}

	// Print the response
	fmt.Println(string(body))
	return true

}
