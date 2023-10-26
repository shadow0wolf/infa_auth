package infa_auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/lwithers/minijks/jks"
	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension/auth"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
	"time"
)

type infaAuthExtension struct {
	cfg                  *Config
	logger               *zap.Logger
	sessionServiceClient *http.Client
}

func newExtension(ctx context.Context, cfg *Config, logger *zap.Logger) (auth.Server, error) {

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
	if strings.TrimSpace(e.cfg.ValidationURL) == "" {
		return errors.New("ValidationURL is empty")
	}

	if strings.TrimSpace(e.cfg.Headerkey) == "" {
		return errors.New("Headerkey is empty")
	}

	sessionServiceClient, err := getClient(e.cfg)
	e.sessionServiceClient = sessionServiceClient
	if err != nil {
		log.Debug("error while creating client in extension.start")
		return err
	}

	log.Debugf("finished executing extension.start , sessionServiceClient : %A", sessionServiceClient)
	return nil
}

func getJksKeystore(filename string, password string) (*jks.Keystore, error) {

	jksContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	log.Debug("read keystorefile : " + filename)

	var opts *jks.Options
	if strings.TrimSpace(password) == "" {
		log.Debug("password is empty , will use nil options")
	} else {
		opts = &jks.Options{
			Password: password,
		}
	}

	keyStore, err := jks.Parse(jksContent, opts)
	if err != nil {
		return nil, err
	}
	log.Debug("success read keystorefile : " + filename)
	return keyStore, nil
}

func getClientCert(ks jks.Keystore) (*tls.Certificate, error) {
	cert, err := tls.X509KeyPair(ks.Keypairs[0].CertChain[0].Cert.Raw, ks.Keypairs[0].RawKey)
	if err != nil {
		log.Debugf("getClientCert ERR %A", err)
	}
	return &cert, nil
}

func getCACertPool(ks jks.Keystore) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for i := 0; i < len(ks.Certs); i++ {
		log.Debugf("adding CA cert to pool : " + string(ks.Certs[i].Alias))
		pool.AddCert(ks.Certs[i].Cert)
	}
	return pool, nil
}

// authenticate checks whether the given context contains valid auth data. Successfully authenticated calls will always return a nil error and a context with the auth data.
// this associated to the Authenticate() method
func (e *infaAuthExtension) authenticate(ctx context.Context, headers map[string][]string) (context.Context, error) {
	log.Debug("executing extensions.authenticate() ")
	/*
		log.Debugf("headers :: %A", headers)
		log.Debugf("ctx :: %A", ctx)
		log.Debugf("e.cfg :: %A", *e.cfg)
	*/
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
	status, err := validateToken(e, token)
	if err != nil || !status {
		return ctx, err
	}

	//success
	return client.NewContext(ctx, cl), nil
}

func getClient(cfg *Config) (*http.Client, error) {
	var pool *x509.CertPool
	var clientCert *tls.Certificate
	client := &http.Client{
		Timeout: func() time.Duration {
			if cfg.TimeOut > 0 {
				return time.Duration(cfg.TimeOut) * time.Second
			}
			return 2 * time.Second
		}(),
	}
	tr := &http.Transport{}
	client.Transport = tr

	//create client with no TLS
	if cfg.InsecureSkipVerify {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		log.Info("returning client block1 ")
		return client, nil
	} else {
		//create client with TLS , CACert is mandatory clientCert is optional
		caKs, err := getJksKeystore(cfg.CAJksPath, cfg.CAJksPassword)

		if err != nil {
			log.Info("error while reading CA jksKeystore : " + cfg.CAJksPath)
			return nil, err
		}

		pool, err = getCACertPool(*caKs)
		if err != nil {
			log.Info("error while generating CACertPool : " + cfg.CAJksPath)
			return nil, err
		}

		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: false, RootCAs: pool}

		if cfg.ClientSideSsl {
			clientKs, err := getJksKeystore(cfg.ClientJksPath, cfg.ClientJksPassword)
			if err != nil {
				log.Info("error while reading Client jksKeystore : " + cfg.ClientJksPath)
				return nil, err
			}

			clientCert, err = getClientCert(*clientKs)
			if err != nil {
				log.Info("error while creating Client cert : " + cfg.ClientJksPath)
				return nil, err
			}

			tr.TLSClientConfig.Certificates = []tls.Certificate{*clientCert}
		}

		log.Info("returning client block2 ")
		client.Transport = tr
		return client, nil
	}
}

// this method creates a HTTP client and makes HTTP/S request to session service, if http status is 200 , it returns
// true with nil error , otherwise non-nil error is returned
func validateToken(e *infaAuthExtension, sessionToken string) (bool, error) {
	// Create an HTTP client
	log.Debug("calling extension.validateToken() ")
	req, err := http.NewRequest("GET", e.cfg.ValidationURL, nil)
	if err != nil {
		log.Debugf("Error creating request:", err)
		return false, err
	}

	req.Header.Add(e.cfg.Headerkey, sessionToken)

	// Make the request
	log.Debugf("invoking http request : %A", req)
	resp, err := e.sessionServiceClient.Do(req)
	log.Debugf("got response %A", resp)

	if err != nil {
		log.Debugf("Error making request:", err)
		return false, err
	}

	if resp.StatusCode != 200 {
		log.Debug("http status is not 200 in response")
		body := make([]byte, 1024)
		_, err := resp.Body.Read(body)
		if err != nil {
			log.Debugf("Error reading response body:", err)
			return false, err
		}
		return false, errors.New("http status is not 200 in response response body is " + string(body))
	}

	defer resp.Body.Close()
	return true, nil
}
