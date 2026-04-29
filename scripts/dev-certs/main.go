// Package main generates a certificate/key pair for CA
// and certificate/key pairs for server, signed by the CA ( only for development )
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

const certsDir = "certs"

const (
	caValidity     = 10 * 365 * 24 * time.Hour
	serverValidity = 365 * 24 * time.Hour
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	err := os.MkdirAll(certsDir, 0o750)
	if err != nil {
		return fmt.Errorf("create certs dir: %w", err)
	}
	certCA, certDER, key, err := generateCA()
	if err != nil {
		return fmt.Errorf("generate CA: %w", err)
	}
	serverCertDER, serverKey, err := generateServer(certCA, key)
	if err != nil {
		return fmt.Errorf("generate server certificate: %w", err)
	}
	err = writeCertPEM(filepath.Join(certsDir, "ca.crt"), certDER)
	if err != nil {
		return fmt.Errorf("write CA certificate: %w", err)
	}
	err = writeCertPEM(filepath.Join(certsDir, "server.crt"), serverCertDER)
	if err != nil {
		return fmt.Errorf("write server certificate: %w", err)
	}
	err = writeKeyPEM(filepath.Join(certsDir, "server.key"), serverKey)
	if err != nil {
		return fmt.Errorf("write server's key: %w", err)
	}
	return nil
}

func generateCA() (*x509.Certificate, []byte, *ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, nil, err
	}

	serialMax := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, err := rand.Int(rand.Reader, serialMax)
	if err != nil {
		return nil, nil, nil, err
	}

	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "TLN dev CA", Organization: []string{"TLN"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(caValidity),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	certCA, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, nil, nil, err
	}
	return certCA, certDER, privateKey, nil
}

func generateServer(certCA *x509.Certificate, caKey *ecdsa.PrivateKey) ([]byte, *ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	template := &x509.Certificate{
		Subject:     pkix.Name{CommonName: "tln-dev-server"},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(serverValidity),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"localhost"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, certCA, &privateKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}
	return certDER, privateKey, nil
}

func writeCertPEM(path string, certDER []byte) error {
	block := &pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	pemData := pem.EncodeToMemory(block)
	return os.WriteFile(path, pemData, 0o600)
}

func writeKeyPEM(path string, key *ecdsa.PrivateKey) error {
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	block := &pem.Block{Type: "PRIVATE KEY", Bytes: keyDER}
	pemData := pem.EncodeToMemory(block)
	return os.WriteFile(path, pemData, 0o600)
}
