package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// LoadTLSConfig reads CA-certificate from caFile, puts it in RootCAs
// and returns a tls.Config, ready for tls.Dial
func LoadTLSConfig(caFile, serverName string) (*tls.Config, error) {
	bytes, err := os.ReadFile(caFile) //nolint:gosec // G304: config-driven path
	if err != nil {
		return nil, fmt.Errorf("read root certificate %q: %w", caFile, err)
	}
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(bytes)
	if !ok {
		return nil, fmt.Errorf("failed to parse root certificate %q", caFile)
	}
	return &tls.Config{
		RootCAs:    pool,
		ServerName: serverName,
		MinVersion: tls.VersionTLS12,
	}, nil
}
