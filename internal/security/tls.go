package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"gproc/pkg/types"
)

type TLSManager struct {
	config *types.TLSConfig
	cert   tls.Certificate
	caCert *x509.Certificate
}

func NewTLSManager(config *types.TLSConfig) (*TLSManager, error) {
	if !config.Enabled {
		return &TLSManager{config: config}, nil
	}

	tm := &TLSManager{config: config}
	
	// Load certificate and key
	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %v", err)
	}
	tm.cert = cert

	// Load CA certificate if provided
	if config.CAFile != "" {
		caCertPEM, err := ioutil.ReadFile(config.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %v", err)
		}
		
		caCert, err := x509.ParseCertificate(caCertPEM)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CA certificate: %v", err)
		}
		tm.caCert = caCert
	}

	return tm, nil
}

func (tm *TLSManager) GetTLSConfig() *tls.Config {
	if !tm.config.Enabled {
		return nil
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tm.cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// Configure client certificate verification if CA is provided
	if tm.caCert != nil {
		caCertPool := x509.NewCertPool()
		caCertPool.AddCert(tm.caCert)
		tlsConfig.ClientCAs = caCertPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return tlsConfig
}

func (tm *TLSManager) WrapHTTPServer(server *http.Server) {
	if tm.config.Enabled {
		server.TLSConfig = tm.GetTLSConfig()
	}
}

func (tm *TLSManager) GetHTTPClient() *http.Client {
	if !tm.config.Enabled {
		return http.DefaultClient
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Add CA certificate for server verification
	if tm.caCert != nil {
		caCertPool := x509.NewCertPool()
		caCertPool.AddCert(tm.caCert)
		tlsConfig.RootCAs = caCertPool
	}

	// Add client certificate for mutual TLS
	if tm.cert.Certificate != nil {
		tlsConfig.Certificates = []tls.Certificate{tm.cert}
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &http.Client{Transport: transport}
}