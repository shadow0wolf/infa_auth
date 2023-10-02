package infa_auth

// Config has the configuration for the INFA AUTH Authenticator extension.
type Config struct {
	TimeOut            int    `mapstructure:"time_out"`
	ClientSideSsl      bool   `mapstructure:"client_side_ssl"`
	ValidationURL      string `mapstructure:"validation_url"`
	Headerkey          string `mapstructure:"header_key"`
	ClientCertPath     string `mapstructure:"client_cert_path"`
	CACertPath         string `mapstructure:"ca_cert_path"`
	InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify"`
	ClientKeyPath      string `mapstructure:"client_key_path"`
}
