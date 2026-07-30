package api

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

type TLSConfig struct {
	CertFile string
	KeyFile  string
	AutoCert bool
}

func DefaultTLSConfig(homeDir string) *TLSConfig {
	return &TLSConfig{
		CertFile: filepath.Join(homeDir, "cert.pem"),
		KeyFile:  filepath.Join(homeDir, "key.pem"),
		AutoCert: true,
	}
}

func LoadTLSConfig(cfg *TLSConfig) (*tls.Config, error) {
	if cfg.CertFile == "" || cfg.KeyFile == "" {
		return nil, fmt.Errorf("tls: cert and key paths required")
	}

	if cfg.AutoCert {
		if err := ensureCert(cfg.CertFile, cfg.KeyFile); err != nil {
			return nil, fmt.Errorf("tls auto-cert: %w", err)
		}
	}

	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("tls load key pair: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}

func ensureCert(certFile, keyFile string) error {
	if _, err := os.Stat(certFile); err == nil {
		if _, err := os.Stat(keyFile); err == nil {
			return nil
		}
	}

	cert, key, err := generateSelfSignedCert()
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(certFile), 0700); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	if err := os.WriteFile(certFile, cert, 0644); err != nil {
		return fmt.Errorf("write cert: %w", err)
	}

	if err := os.WriteFile(keyFile, key, 0600); err != nil {
		return fmt.Errorf("write key: %w", err)
	}

	return nil
}

func generateSelfSignedCert() (certPEM, keyPEM []byte, err error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("generate key: %w", err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, fmt.Errorf("serial: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   "GhostStack",
			Organization: []string{"GhostStack"},
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
		DNSNames:              []string{"localhost", "ghoststack.local"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("create cert: %w", err)
	}

	certPEM = append([]byte("-----BEGIN CERTIFICATE-----\n"), encodeBase64(certDER)...)
	certPEM = append(certPEM, []byte("\n-----END CERTIFICATE-----\n")...)

	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal key: %w", err)
	}

	keyPEM = append([]byte("-----BEGIN EC PRIVATE KEY-----\n"), encodeBase64(b)...)
	keyPEM = append(keyPEM, []byte("\n-----END EC PRIVATE KEY-----\n")...)

	return certPEM, keyPEM, nil
}

func encodeBase64(data []byte) []byte {
	const enc = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	encoded := make([]byte, ((len(data)+2)/3)*4)
	for i, j := 0, 0; i < len(data); i += 3 {
		val := uint(data[i]) << 16
		if i+1 < len(data) {
			val |= uint(data[i+1]) << 8
		}
		if i+2 < len(data) {
			val |= uint(data[i+2])
		}
		encoded[j] = enc[(val>>18)&0x3F]
		encoded[j+1] = enc[(val>>12)&0x3F]
		if i+1 < len(data) {
			encoded[j+2] = enc[(val>>6)&0x3F]
		} else {
			encoded[j+2] = '='
		}
		if i+2 < len(data) {
			encoded[j+3] = enc[val&0x3F]
		} else {
			encoded[j+3] = '='
		}
		j += 4
	}
	return encoded
}
