package certs

import (
	"errors"
	"fmt"
	"github.com/Litekube/network-controller/config"
	"github.com/Litekube/network-controller/contant"
	"github.com/Litekube/network-controller/utils"
	"os"
	"path/filepath"
)

func CheckGrpcClientCertConfig(tlsConfig config.TLSConfig, tlsDir string) error {
	if tlsDir == "" {
		return errors.New("tlsDir can't be empty")
	}
	// generate default client certs for grpc
	defualtDir := filepath.Join(utils.GetHomeDir(), ".litekube/nc/certs/grpc")
	utils.CreateDir(defualtDir)
	caFile := filepath.Join(defualtDir, contant.CAFile)

	if !utils.Exists(caFile) {
		err := utils.CopyFile(tlsConfig.CAFile, caFile)
		if err != nil {
			return err
		}
	}

	clientCertFile := filepath.Join(defualtDir, contant.ClientCertFile)
	clientKeyFile := filepath.Join(defualtDir, contant.ClientKeyFile)

	if _, _, _, err := GenerateClientCertKey(false, true, "register-client", []string{"lknm:register"}, tlsConfig.CAFile, tlsConfig.CAKeyFile, clientCertFile, clientKeyFile); err != nil {
		return err
	}

	// soft link with tlsDir
	newTlsDir := filepath.Join(tlsDir, "/grpc")
	utils.CreateDir(newTlsDir)
	caPath := filepath.Join(newTlsDir, "ca.crt")
	certPath := filepath.Join(newTlsDir, "client.crt")
	keyPath := filepath.Join(newTlsDir, "client.key")

	// create cert symlink
	if !utils.Exists(caPath) {
		if err := os.Symlink(caFile, caPath); err != nil {
			return fmt.Errorf("fail to create link for network-manager adm certs err:%s", err.Error())
		}
	}
	if !utils.Exists(certPath) {
		if err := os.Symlink(clientCertFile, certPath); err != nil {
			return fmt.Errorf("fail to create link for network-manager adm certs err:%s", err.Error())
		}
	}
	if !utils.Exists(keyPath) {
		if err := os.Symlink(clientKeyFile, keyPath); err != nil {
			return fmt.Errorf("fail to create link for network-manager adm certs err:%s", err.Error())
		}
	}

	return nil
}
