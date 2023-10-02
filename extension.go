package infa_auth

import (
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"

	"crypto/x509"
	"github.com/gofiber/fiber/v2/log"
	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension/auth"
	"go.uber.org/zap"
	"os"
	"time"
)

type infaAuthExtension struct {
	cfg    *Config
	logger *zap.Logger
}

func newExtension(ctx context.Context, cfg *Config, logger *zap.Logger) (auth.Server, error) {
	/*
		if cfg.ValidationURL == "" {
			return nil, errors.New("validation url is empty")
		}

		if cfg.Headerkey == "" {
			return nil, errors.New("header key is empty")
		}
	*/
	oe := &infaAuthExtension{
		cfg:    cfg,
		logger: logger,
	}
	return auth.NewServer(auth.WithServerStart(oe.start), auth.WithServerAuthenticate(oe.authenticate)), nil
}

// function that executes to initialize the auth extension , returns no error and nil in case of success
// associated to Start()
func (e *infaAuthExtension) start(context.Context, component.Host) error {
	log.Debug("begin executing extension.start")

	//validation url and header key must NOT be empty
	if e.cfg.ValidationURL == "" {
		return errors.New("ValidationURL is empty")
	}

	if e.cfg.Headerkey == "" {
		return errors.New("Headerkey is empty")
	}

	/*
		if e.cfg.CACertPath == "" {
			log.Debug("CACertPath is empty")
			return errors.New("CACertPath is empty")
		}

		log.Debugf("CACertPath is : %s ", e.cfg.CACertPath)
		_, error := os.Stat(e.cfg.CACertPath)
		if error != nil {
			log.Debugf("error reading CACertPath")
			return errors.New("error reading CACertPath")
		}
	*/

	//validate client cert
	log.Debugf("ClientSideSsl is : %s ", e.cfg.ClientSideSsl)
	if e.cfg.ClientSideSsl {
		log.Debugf("ClientCertPath is : %s ", e.cfg.ClientCertPath)
		if e.cfg.ClientCertPath == "" {
			log.Debug("ClientCertPath is empty")
			return errors.New("ClientCertPath is empty")

		} else {
			_, error := os.Stat(e.cfg.ClientCertPath)
			if error != nil {
				log.Debugf("error reading ClientCert")
				return errors.New("error reading ClientCert")
			}
		}

		log.Debugf("ClientKeyPath is : %s ", e.cfg.ClientKeyPath)
		if e.cfg.ClientKeyPath == "" {
			log.Debug("ClientKeyPath is empty")
			return errors.New("ClientKeyPath is empty")

		} else {
			_, error := os.Stat(e.cfg.ClientKeyPath)
			if error != nil {
				log.Debugf("error reading ClientKeyPath")
				return errors.New("error reading ClientKeyPath")
			}
		}

	}

	log.Debug("finished executing extension.start")
	return nil
}

// authenticate checks whether the given context contains valid auth data. Successfully authenticated calls will always return a nil error and a context with the auth data.
// this associated to the Authenticate() method
func (e *infaAuthExtension) authenticate(ctx context.Context, headers map[string][]string) (context.Context, error) {
	log.Debug("executing extensions.authenticate() ")
	log.Debugf("headers :: %A", headers)
	log.Debugf("ctx :: %A", ctx)
	log.Debugf("e.cfg :: %A", *e.cfg)

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
				return ctx, errors.New("ids-agent-session-id header does not exist")
			}
		}
	}

	token := h[0]
	log.Debugf("token is %A", token)

	if len(token) == 0 {
		return ctx, errors.New("ids-agent-session-id header value is empty string")
	}

	cl := client.FromContext(ctx)
	status, err := validateToken(e.cfg, token)
	if err != nil || status == false {
		return ctx, err
	}

	//success
	return client.NewContext(ctx, cl), nil
}

// this method creates a HTTP client and makes HTTP/S request to session service, if http status is 200 , it returns
// true with nil error , otherwise non-nil error is returned
func validateToken(cfg *Config, sessionToken string) (bool, error) {
	// Create an HTTP client
	log.Debug("calling extension.validateToken() ")

	tr := &http.Transport{}

	//read ca-cert file
	pool := x509.NewCertPool()
	if cfg.CACertPath != "" {
		cert, err := ioutil.ReadFile(cfg.CACertPath)
		if err != nil {
			log.Debugf("Error reading CA certificate: %A", err)
			return false, err
		} else {
			pool.AppendCertsFromPEM(cert)
			tr = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify, RootCAs: pool},
			}
		}
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify},
		}
	}

	// Load client certificate and private key
	if cfg.ClientSideSsl {
		cert, err := tls.LoadX509KeyPair(cfg.ClientCertPath, cfg.ClientKeyPath)
		if err != nil {
			log.Debugf("Error loading client certificate:", err)
			return false, err
		}
		config := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            pool,
			InsecureSkipVerify: cfg.InsecureSkipVerify,
		}
		tr = &http.Transport{
			TLSClientConfig: config,
		}
	}

	client := &http.Client{}
	if cfg.TimeOut > 0 {
		client = &http.Client{
			Timeout:   time.Duration(cfg.TimeOut) * time.Second,
			Transport: tr}
	} else {
		client = &http.Client{
			Timeout:   2 * time.Second,
			Transport: tr}
	}

	req, err := http.NewRequest("GET", cfg.ValidationURL, nil)
	if err != nil {
		log.Debugf("Error creating request:", err)
		return false, err
	}

	req.Header.Add(cfg.Headerkey, sessionToken)

	// Make the request
	log.Debugf("invoking http request : %A", req)
	resp, err := client.Do(req)
	log.Debugf("got response %A", resp)

	if err != nil {
		log.Debugf("Error making request:", err)
		return false, err
	}

	if resp.StatusCode != 200 {
		log.Debug("http status is not 200 in response")
		return false, errors.New("http status is not 200 in response")
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debugf("Error reading response body:", err)
		return false, err
	}

	// Print the response
	log.Debugf("response body is %A", string(body))
	return true, nil

}
