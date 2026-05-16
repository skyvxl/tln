package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

// LoadTLSConfig reads CA-certificate from caFile, puts it in RootCAs
// and returns a tls.Config, ready for tls.Dial
func LoadTLSConfig(caFile, serverName string) (*tls.Config, error) {
	bytes, err := os.ReadFile(caFile) //nolint:gosec // G304: config-driven path
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(bytes)
	if !ok {
		return nil, errors.New("failed to parse root certificate")
	}
	return &tls.Config{
		RootCAs:    pool,
		ServerName: serverName,
		MinVersion: tls.VersionTLS12,
	}, nil
}
