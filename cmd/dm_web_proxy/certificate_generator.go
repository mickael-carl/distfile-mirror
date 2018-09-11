package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// CertificateGenerator implements a GetCertificate hook for tls.Config
// to generate SSL certificates on demand, signed by a given CA. For
// simplicity, certificates will use the same private key as the CA.
type CertificateGenerator struct {
	caParsedCertificate *x509.Certificate
	caEncodedPrivateKey []byte
	caParsedPrivateKey  *rsa.PrivateKey
}

func NewCertificateGenerator(caEncodedCertificate []byte, caEncodedPrivateKey []byte) (*CertificateGenerator, error) {
	// Parse CA certificate.
	caDecodedCertificate, _ := pem.Decode(caEncodedCertificate)
	if caDecodedCertificate == nil {
		return nil, errors.New("Failed to parse certificate")
	}
	caParsedCertificate, err := x509.ParseCertificate(caDecodedCertificate.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse certificate: %s", err)
	}

	// Parse CA private key.
	caDecodedPrivateKey, _ := pem.Decode(caEncodedPrivateKey)
	if caDecodedPrivateKey == nil {
		return nil, errors.New("Failed to parse private key")
	}
	caParsedPrivateKey, err := x509.ParsePKCS1PrivateKey(caDecodedPrivateKey.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key: %s", err)
	}

	return &CertificateGenerator{
		caParsedCertificate: caParsedCertificate,
		caEncodedPrivateKey: caEncodedPrivateKey,
		caParsedPrivateKey:  caParsedPrivateKey,
	}, nil
}

func (cg *CertificateGenerator) GetCertificate(clientHelloInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if clientHelloInfo.ServerName == "" {
		return nil, errors.New("SNI is required for this service")
	}

	// Determine certificate fields.
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate certificate serial number: %s", err)
	}
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24)

	// Generate certificate and encode it.
	decodedCertificate, err := x509.CreateCertificate(
		rand.Reader,
		&x509.Certificate{
			SerialNumber:          serialNumber,
			NotBefore:             notBefore,
			NotAfter:              notAfter,
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
			DNSNames:              []string{clientHelloInfo.ServerName},
		},
		cg.caParsedCertificate,
		&cg.caParsedPrivateKey.PublicKey,
		cg.caParsedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate certificate: %s", err)
	}
	encodedCertificate := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: decodedCertificate,
	})

	keyPair, err := tls.X509KeyPair(encodedCertificate, cg.caEncodedPrivateKey)
	return &keyPair, err
}
