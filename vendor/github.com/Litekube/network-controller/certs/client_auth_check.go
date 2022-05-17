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
	newTlsDir := filepath.Join(tlsDir, "/grpc")
	utils.CreateDir(newTlsDir)
	caPath := filepath.Join(newTlsDir, "ca.crt")
	certPath := filepath.Join(newTlsDir, "client.crt")
	keyPath := filepath.Join(newTlsDir, "client.key")

	// soft link with tlsDir
	defualtDir := filepath.Join(utils.GetHomeDir(), ".litekube/nc/certs/grpc")
	caFile := filepath.Join(defualtDir, contant.CAFile)
	clientCertFile := filepath.Join(defualtDir, contant.ClientCertFile)
	clientKeyFile := filepath.Join(defualtDir, contant.ClientKeyFile)

	// in tls dir: generate default client certs for grpc
	if !utils.Exists(caPath) {
		err := utils.CopyFile(tlsConfig.CAFile, caPath)
		if err != nil {
			return err
		}
	}

	if _, _, _, err := GenerateClientCertKey(false, true, "register-client", []string{"lknm:register"}, tlsConfig.CAFile, tlsConfig.CAKeyFile, certPath, keyPath); err != nil {
		return err
	}

	if err := os.RemoveAll(defualtDir); err != nil {
		return err
	}
	utils.CreateDir(defualtDir)

	// in ncadm certs dir: create cert symlink
	if err := os.Symlink(caPath, caFile); err != nil {
		return fmt.Errorf("fail to create link for network-controller adm certs err:%s", err.Error())
	}
	if err := os.Symlink(certPath, clientCertFile); err != nil {
		return fmt.Errorf("fail to create link for network-controller adm certs err:%s", err.Error())
	}
	if err := os.Symlink(keyPath, clientKeyFile); err != nil {
		return fmt.Errorf("fail to create link for network-controller adm certs err:%s", err.Error())
	}

	return nil
}
