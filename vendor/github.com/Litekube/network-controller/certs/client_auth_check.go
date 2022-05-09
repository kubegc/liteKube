package certs

import (
	"github.com/Litekube/network-controller/config"
	"github.com/Litekube/network-controller/contant"
	"github.com/Litekube/network-controller/utils"
	"os"
	"path/filepath"
)

var homeDir string = func() string {
	if home, err := os.UserHomeDir(); err != nil {
		return ""
	} else {
		return home
	}
}()

func CheckGrpcClientCertConfig(tlsConfig config.TLSConfig) error {
	// generate default client certs for grpc
	defualtDir := filepath.Join(homeDir, ".litekube/nc/certs/grpc")
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

	//if _, _, _, err := GenerateClientCertKey(false, true, "network-controller-grpc-client", []string{"network-controller-grpc"}, tlsConfig.CAFile, tlsConfig.CAKeyFile, clientCertFile, clientKeyFile); err != nil {
	if _, _, _, err := GenerateClientCertKey(false, true, "register-client", []string{"lknm:register"}, tlsConfig.CAFile, tlsConfig.CAKeyFile, clientCertFile, clientKeyFile); err != nil {
		return err
	}
	return nil
}
