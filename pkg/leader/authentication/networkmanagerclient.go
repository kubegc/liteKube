package authentication

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"github.com/litekube/LiteKube/pkg/certificate"
	"github.com/litekube/LiteKube/pkg/global"
	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
	"k8s.io/klog/v2"
)

type NetworkManagerClient struct {
	ManagerRootCertPath    string
	ManagerCertDir         string
	RegisterManagerCertDir string
	RegisterCACert         string
	RegisterClientCert     string
	RegisterClientkey      string
	JoinManagerCertDir     string
	JoinCACert             string
	JoinClientCert         string
	JoinClientkey          string
	Token                  string
	NodeTokenPath          string
}

func NewNetworkManagerClient(rootCertPath string, token string) *NetworkManagerClient {
	if token == "" {
		token = "local"
	}

	if rootCertPath == "" {
		rootCertPath = filepath.Join(globaloptions.DefaultGO.WorkDir, "tls/")
	}

	managerRootCertPath := filepath.Join(rootCertPath, "network-manager/")
	managerCertDir := filepath.Join(managerRootCertPath, token)
	registerManagerCertDir := filepath.Join(managerCertDir, "register")
	joinManagerCertDir := filepath.Join(managerCertDir, "join")

	return &NetworkManagerClient{
		ManagerRootCertPath:    managerRootCertPath,
		ManagerCertDir:         managerCertDir,
		RegisterManagerCertDir: registerManagerCertDir,
		RegisterCACert:         filepath.Join(registerManagerCertDir, "ca.crt"),
		RegisterClientCert:     filepath.Join(registerManagerCertDir, "client.crt"),
		RegisterClientkey:      filepath.Join(registerManagerCertDir, "client.key"),
		JoinManagerCertDir:     joinManagerCertDir,
		JoinCACert:             filepath.Join(joinManagerCertDir, "ca.crt"),
		JoinClientCert:         filepath.Join(joinManagerCertDir, "client.crt"),
		JoinClientkey:          filepath.Join(joinManagerCertDir, "client.key"),
		Token:                  token,
		NodeTokenPath:          filepath.Join(managerCertDir, "node.token"),
	}
}

func (na *NetworkManagerClient) Nodetoken() (string, error) {
	bytes, err := ioutil.ReadFile(na.NodeTokenPath)
	return string(bytes), err
}

// generate X.509 certificate for network-manager
func (na *NetworkManagerClient) GenerateOrSkip(address string, port int) error {
	if na.Check() {
		// file exist
		return nil
	}

	if na == nil {
		return fmt.Errorf("nil network authentication")
	}

	if na.Token == "local" {
		return na.CreateSoftlinkForClient()
	} else {
		return na.TLSBootStrap(address, port)
	}
}

// to be finish
func (na *NetworkManagerClient) TLSBootStrap(address string, port int) error {
	if address == "" || port < 1 {
		return fmt.Errorf("none tls bootstrap address and port for network-manager")
	}

	err := os.MkdirAll(na.ManagerCertDir, os.ModePerm)
	if err != nil {
		return err
	}

	// generate certificate and node-token here.
	return nil
}

func (na *NetworkManagerClient) QueryRemoteIps() []net.IP {
	return []net.IP{}
}

func (na *NetworkManagerClient) QueryRemoteDNSNames() []string {
	return []string{}
}

func (na *NetworkManagerClient) Check() bool {
	return certificate.Exists(na.RegisterCACert, na.RegisterClientCert, na.RegisterClientkey, na.JoinCACert, na.JoinClientCert, na.JoinClientkey, na.NodeTokenPath)
}

func (na *NetworkManagerClient) CreateSoftlinkForClient() error {
	if na.Check() {
		return nil
	}

	if err := os.MkdirAll(na.RegisterManagerCertDir, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(na.JoinManagerCertDir, os.ModePerm); err != nil {
		return err
	}

	// create cert symlink for register
	if err := os.Symlink(filepath.Join(na.ManagerRootCertPath, "register/ca.key"), na.RegisterCACert); err != nil {
		klog.Warningf("fail to create symlink for network-manager certificat err:%s", err.Error())
	}
	if err := os.Symlink(filepath.Join(na.ManagerRootCertPath, "register/client.crt"), na.RegisterClientCert); err != nil {
		klog.Warningf("fail to create symlink for network-manager certificat err:%s", err.Error())
	}
	if err := os.Symlink(filepath.Join(na.ManagerRootCertPath, "register/client.key"), na.RegisterClientkey); err != nil {
		klog.Warningf("fail to create symlink for network-manager certificat err:%s", err.Error())
	}

	// create cert symlink for join
	if err := os.Symlink(filepath.Join(na.ManagerRootCertPath, "join/ca.key"), na.JoinCACert); err != nil {
		klog.Warningf("fail to create symlink for network-manager certificat err:%s", err.Error())
	}
	if err := os.Symlink(filepath.Join(na.ManagerRootCertPath, "join/client.crt"), na.JoinClientCert); err != nil {
		klog.Warningf("fail to create symlink for network-manager certificat err:%s", err.Error())
	}
	if err := os.Symlink(filepath.Join(na.ManagerRootCertPath, "join/client.key"), na.JoinClientkey); err != nil {
		klog.Warningf("fail to create symlink for network-manager certificat err:%s", err.Error())
	}

	if err := ioutil.WriteFile(na.NodeTokenPath, []byte(global.ReservedNodeToken), os.FileMode(0644)); err != nil {
		klog.Warningf("fail to create node token")
	}

	if !na.Check() {
		return fmt.Errorf("fail to create symlink for network-manager certificate")
	}

	return nil
}
