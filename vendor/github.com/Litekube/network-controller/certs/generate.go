package certs

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"github.com/Litekube/network-controller/utils"
	certutil "github.com/rancher/dynamiclistener/cert"
	"io/ioutil"
	"time"
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

func GenerateServerCertKey(regen bool, commonName string, organization []string, altNames *certutil.AltNames, caCertPath, caKeyPath, certPath, keyPath string) ([]byte, []byte, bool, error) {
	if !ValidateTLSPair(caCertPath, caKeyPath) {
		return []byte{}, []byte{}, false, fmt.Errorf("bad CA")
	}

	// already exist and valid
	if !regen && utils.Exists(certPath, keyPath) && ValidateCA(certPath, caCertPath) {
		keyBytes, _, _ := certutil.LoadOrGenerateKeyFile(keyPath, regen)
		cert, _ := LoadCertificate(caCertPath)
		return keyBytes, certutil.EncodeCertPEM(cert), false, nil
	}

	return GenerateCertKey(regen, commonName, organization, altNames, []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, caCertPath, caKeyPath, certPath, keyPath)
}

func GenerateClientCertKey(regen bool, commonName string, organization []string, caCertPath, caKeyPath, certPath, keyPath string) ([]byte, []byte, bool, error) {
	if !ValidateTLSPair(caCertPath, caKeyPath) {
		return []byte{}, []byte{}, false, fmt.Errorf("bad CA")
	}

	// already exist and valid
	//if !regen && utils.Exists(certPath, keyPath) && ValidateCA(certPath, caCertPath) {
	//	keyBytes, _, _ := certutil.LoadOrGenerateKeyFile(keyPath, regen)
	//	cert, _ := LoadCertificate(caCertPath)
	//	return keyBytes, certutil.EncodeCertPEM(cert), false, nil
	//}
	// always re-generate for new client
	return GenerateCertKey(true, commonName, organization, nil, []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, caCertPath, caKeyPath, certPath, keyPath)
}

// set regen=true to force gen new cert
func GenerateCertKey(regen bool, commonName string, organization []string, altNames *certutil.AltNames, extKeyUsage []x509.ExtKeyUsage, caCertPath, caKeyPath, certPath, keyPath string) ([]byte, []byte, bool, error) {
	caBytes, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return []byte{}, []byte{}, false, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caBytes)

	var flag bool
	// check for certificate expiration
	if !regen {
		regen = Expired(certPath, pool, 10)
		flag = regen
	}

	if !regen {
		if utils.Exists(certPath, keyPath) {
			return []byte{}, []byte{}, false, nil
		}
	}

	caKeyBytes, err := ioutil.ReadFile(caKeyPath)
	if err != nil {
		return []byte{}, []byte{}, false, err
	}

	caKey, err := certutil.ParsePrivateKeyPEM(caKeyBytes)
	if err != nil {
		return []byte{}, []byte{}, false, err
	}

	caCert, err := certutil.ParseCertsPEM(caBytes)
	if err != nil {
		return []byte{}, []byte{}, false, err
	}

	keyBytes, _, err := certutil.LoadOrGenerateKeyFile(keyPath, regen)
	if err != nil {
		return []byte{}, []byte{}, false, err
	}

	key, err := certutil.ParsePrivateKeyPEM(keyBytes)
	if err != nil {
		return []byte{}, []byte{}, false, err
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
		return []byte{}, []byte{}, false, err
	}

	err = certutil.WriteCert(certPath, append(certutil.EncodeCertPEM(cert), certutil.EncodeCertPEM(caCert[0])...))
	if err != nil {
		return []byte{}, []byte{}, false, err
	}

	// for local admin use
	if !utils.Exists(certPath, keyPath) || flag {
		fmt.Println("1111")
		certutil.WriteCert(certPath, certutil.EncodeCertPEM(cert))
	}

	return keyBytes, certutil.EncodeCertPEM(cert), true, nil
}
