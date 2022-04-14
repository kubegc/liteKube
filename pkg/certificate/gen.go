package certificate

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	certutil "github.com/rancher/dynamiclistener/cert"
)

// write signing Certkey with `CN: {prefix}-ca@{time.Now().Unix()}` to certFile
// if keyFile valid, it will be use, or generate.
// if generate new ,then return (true,nil); or return (false,nil)
func GenerateSigningCertKey(regen bool, prefix, certFile, keyFile string) (bool, error) {
	// file exist and valid
	if !regen && ValidateTLSPair(certFile, keyFile) {
		return false, nil
	}

	caKeyBytes, _, err := certutil.LoadOrGenerateKeyFile(keyFile, false)
	if err != nil {
		return false, err
	}

	caKey, err := certutil.ParsePrivateKeyPEM(caKeyBytes)
	if err != nil {
		return false, err
	}

	cfg := certutil.Config{
		CommonName: fmt.Sprintf("%s-ca@%d", prefix, time.Now().Unix()),
	}

	cert, err := certutil.NewSelfSignedCACert(cfg, caKey.(crypto.Signer))
	if err != nil {
		return false, err
	}

	if err := certutil.WriteCert(certFile, certutil.EncodeCertPEM(cert)); err != nil {
		return false, err
	}
	return true, nil
}

func GenerateServerCertKey(regen bool, commonName string, organization []string, altNames *certutil.AltNames, caCertPath, caKeyPath, certPath, keyPath string) (bool, error) {
	if !ValidateTLSPair(caCertPath, caKeyPath) {
		return false, fmt.Errorf("bad CA")
	}

	// already exist and valid
	if !regen && !Exists(certPath, keyPath) && ValidateCA(certPath, caCertPath) {
		return false, nil
	}

	if _, err := GenerateCertKey(regen, commonName, organization, altNames, []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, caCertPath, caKeyPath, certPath, keyPath); err != nil {
		return false, err
	}

	return true, nil

}

func GenerateClientCertKey(regen bool, commonName string, organization []string, caCertPath, caKeyPath, certPath, keyPath string) (bool, error) {
	if !ValidateTLSPair(caCertPath, caKeyPath) {
		return false, fmt.Errorf("bad CA")
	}

	// already exist and valid
	if !regen && !Exists(certPath, keyPath) && ValidateCA(certPath, caCertPath) {
		return false, nil
	}

	if _, err := GenerateCertKey(regen, commonName, organization, nil, []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, caCertPath, caKeyPath, certPath, keyPath); err != nil {
		return false, err
	}

	return true, nil
}

// set regen=true to force gen new cert
func GenerateCertKey(regen bool, commonName string, organization []string, altNames *certutil.AltNames, extKeyUsage []x509.ExtKeyUsage, caCertPath, caKeyPath, certPath, keyPath string) (bool, error) {
	caBytes, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return false, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caBytes)

	// check for certificate expiration
	if !regen {
		regen = Expired(certPath, pool, 10)
	}

	if !regen {
		if Exists(certPath, keyPath) {
			return false, nil
		}
	}

	caKeyBytes, err := ioutil.ReadFile(caKeyPath)
	if err != nil {
		return false, err
	}

	caKey, err := certutil.ParsePrivateKeyPEM(caKeyBytes)
	if err != nil {
		return false, err
	}

	caCert, err := certutil.ParseCertsPEM(caBytes)
	if err != nil {
		return false, err
	}

	keyBytes, _, err := certutil.LoadOrGenerateKeyFile(keyPath, regen)
	if err != nil {
		return false, err
	}

	key, err := certutil.ParsePrivateKeyPEM(keyBytes)
	if err != nil {
		return false, err
	}

	cfg := certutil.Config{
		CommonName:   commonName,
		Organization: organization,
		Usages:       extKeyUsage,
	}
	if altNames != nil {
		cfg.AltNames = *altNames
	}
	cert, err := certutil.NewSignedCert(cfg, key.(crypto.Signer), caCert[0], caKey.(crypto.Signer))
	if err != nil {
		return false, err
	}

	return true, certutil.WriteCert(certPath, append(certutil.EncodeCertPEM(cert), certutil.EncodeCertPEM(caCert[0])...))
}
