package certificate

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"

	certutil "github.com/rancher/dynamiclistener/cert"
)

func LoadCertificate(certPath string) (*x509.Certificate, error) {
	certificates, err := LoadCertificates(certPath)
	if err != nil || certificates == nil || len(certificates) < 1 {
		return nil, err
	} else {
		return certificates[0], err
	}

	// bytes, err := ioutil.ReadFile(certPath)
	// if err != nil {
	// 	return nil, err
	// }

	// block, _ := pem.Decode(bytes)
	// if block == nil {
	// 	return nil, fmt.Errorf("fail to decode pem cert")
	// }

	// cert, err := x509.ParseCertificate(block.Bytes)
	// if err != nil {
	// 	return nil, err
	// }
}

// if client/server certificate generate by this package, return[0] is client/server certificate, return[1] is CA certificate
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
