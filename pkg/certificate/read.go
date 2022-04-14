package certificate

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	certutil "github.com/rancher/dynamiclistener/cert"
)

func LoadCertificate(certPath string) (*x509.Certificate, error) {
	bytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("fail to decode pem cert")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func LoadCertificates(certPath string) ([]*x509.Certificate, error) {
	certBytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	certificates, err := certutil.ParseCertsPEM(certBytes)
	if err != nil {
		return nil, err
	}

	return certificates, nil
}

func LoadCertPool(caCertPath string) (*x509.CertPool, error) {
	// read CA
	caBytes, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caBytes) {
		return nil, fmt.Errorf("fail to parse ca certificates")
	}

	return pool, nil
}
