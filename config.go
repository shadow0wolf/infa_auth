package infa_auth

import "net/http"

// Config has the configuration for the INFA AUTH Authenticator extension.
type Config struct {
	TimeOut            int    `mapstructure:"time_out"`
	ClientSideSsl      bool   `mapstructure:"client_side_ssl"`
	ValidationURL      string `mapstructure:"validation_url"`
	Headerkey          string `mapstructure:"header_key"`
	ClientJksPath      string `mapstructure:"client_jks_path"`
	ClientJksPassword  string `mapstructure:"client_jks_password"`
	CAJksPath          string `mapstructure:"ca_jks_path"`
	CAJksPassword      string `mapstructure:"ca_jks_password"`
	InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify"`
}

var sessionServiceClient *http.Client
