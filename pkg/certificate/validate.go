package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	certutil "github.com/rancher/dynamiclistener/cert"
)

// check if certificate and key are one pair
func ValidateTLSPair(certPath string, keyPath string) bool {
	if !Exists(certPath, keyPath) {
		return false
	}

	// validate pair for certificate and key
	if _, err := tls.LoadX509KeyPair(certPath, keyPath); err != nil {
		return false
	}

	return true
}

// verifies that the signature on cert is a valid signature from issuer.
func ValidateIssuer(child *x509.Certificate, issuer *x509.Certificate) bool {
	if child == nil || issuer == nil || child.CheckSignatureFrom(issuer) != nil {
		return false
	}

	return true
}

// validate if cert is valid to this ca
func ValidateCA(certPath string, caCertPath string) bool {
	certPool, err := LoadCertPool(caCertPath)
	if err != nil {
		return false
	}

	certs, err := LoadCertificates(certPath)
	if err != nil {
		return false
	}

	_, err = certs[0].Verify(x509.VerifyOptions{
		Roots: certPool,
		KeyUsages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageAny,
		},
	})
	if err != nil {
		return true
	}
	return true
}

// less than {days} days will be treat as exired
func ValidateExpired(certPath string, caCertPath string, days int) bool {
	if !ValidateCA(certPath, caCertPath) {
		return false
	}

	certs, err := LoadCertificates(certPath)
	if err != nil {
		return false
	}

	return certutil.IsCertExpired(certs[0], days)
}

// less than {days} days will be treat as exired
func Expired(certPath string, pool *x509.CertPool, days int) bool {
	certBytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		return false
	}
	certificates, err := certutil.ParseCertsPEM(certBytes)
	if err != nil {
		return false
	}
	_, err = certificates[0].Verify(x509.VerifyOptions{
		Roots: pool,
		KeyUsages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageAny,
		},
	})
	if err != nil {
		return true
	}
	return certutil.IsCertExpired(certificates[0], days)
}
