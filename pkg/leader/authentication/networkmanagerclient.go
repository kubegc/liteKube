package authentication

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/litekube/LiteKube/pkg/certificate"
	"github.com/litekube/LiteKube/pkg/global"
	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
)

type NetworkManagerClient struct {
	ManagerRootCertPath    string
	ManagerCertDir         string
	RegisterManagerCertDir string
	RegisterAddress        *string // value only tls-bootstrap without init
	RegisterPort           *uint16 // value only tls-bootstrap without init
	RegisterCACert         string
	RegisterClientCert     string
	RegisterClientkey      string
	JoinManagerCertDir     string
	JoinAddress            *string // value only tls-bootstrap without init
	JoinPort               *uint16 // value only tls-bootstrap without init
	JoinCACert             string
	JoinClientCert         string
	JoinClientkey          string
	Token                  string
	NodeTokenPath          string
}

func NewNetworkManagerClient(rootCertPath string, token string, registerAddress *string, registerPort *uint16, joinAddress *string, joinPort *uint16) *NetworkManagerClient {
	if token == "" {
		token = "unknown"
	}

	if rootCertPath == "" {
		rootCertPath = filepath.Join(globaloptions.DefaultGO.WorkDir, "tls/")
	}

	managerRootCertPath := filepath.Join(rootCertPath, "network-manager/")
	managerCertDir := filepath.Join(managerRootCertPath, strings.SplitN(token, "@", 2)[0])
	registerManagerCertDir := filepath.Join(managerCertDir, "register")
	joinManagerCertDir := filepath.Join(managerCertDir, "join")

	return &NetworkManagerClient{
		ManagerRootCertPath:    managerRootCertPath,
		ManagerCertDir:         managerCertDir,
		RegisterManagerCertDir: registerManagerCertDir,
		RegisterAddress:        registerAddress,
		RegisterPort:           registerPort,
		RegisterCACert:         filepath.Join(registerManagerCertDir, "ca.crt"),
		RegisterClientCert:     filepath.Join(registerManagerCertDir, "client.crt"),
		RegisterClientkey:      filepath.Join(registerManagerCertDir, "client.key"),
		JoinManagerCertDir:     joinManagerCertDir,
		JoinAddress:            joinAddress,
		JoinPort:               joinPort,
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
func (na *NetworkManagerClient) GenerateOrSkip() error {
	if na.Token == "unknown" {
		return fmt.Errorf("token is unknown")
	}

	if na == nil {
		return fmt.Errorf("nil network authentication")
	}

	if na.Token == "local" {
		return na.CreatelinkForClient()
	} else {
		Endpoint := strings.SplitN(strings.SplitN(na.Token, "@", 2)[1], ":", 2)
		var port int
		if p, err := strconv.Atoi(Endpoint[1]); err != nil {
			return fmt.Errorf("bad network-bootstrap port")
		} else {
			port = p
		}

		if ip := net.ParseIP(Endpoint[0]); ip == nil {
			return fmt.Errorf("bad network-bootstrap ip")
		}
		return na.TLSBootStrap(Endpoint[0], port)
	}
}

// download certificates and get node-token from network manager
func (na *NetworkManagerClient) TLSBootStrap(address string, port int) error {
	if na.Check() {
		// file exist, bootstrap ok.
		return nil
	}

	if address == "" || port < 1 || port > 65535 {
		return fmt.Errorf("none tls bootstrap address and port for network-manager")
	}

	err := os.MkdirAll(na.ManagerCertDir, os.ModePerm)
	if err != nil {
		return err
	}

	// generate certificate and node-token here.
	// need to value address and port here like:
	// *RegisterAddress="127.0.0.1", *RegisterPort=6440
	// *JoinAddress="127.0.0.1", *JoinPort=6441
	return nil
}

func (na *NetworkManagerClient) Check() bool {
	return certificate.Exists(na.RegisterCACert, na.RegisterClientCert, na.RegisterClientkey, na.JoinCACert, na.JoinClientCert, na.JoinClientkey, na.NodeTokenPath)
}

func (na *NetworkManagerClient) CreatelinkForClient() error {
	registerCACert := filepath.Join(na.ManagerRootCertPath, "register/ca.crt")
	registerClientCert := filepath.Join(na.ManagerRootCertPath, "register/client.crt")
	registerClientKey := filepath.Join(na.ManagerRootCertPath, "register/client.key")

	joinCACert := filepath.Join(na.ManagerRootCertPath, "join/ca.crt")
	joinClienCert := filepath.Join(na.ManagerRootCertPath, "join/client.crt")
	joinClienKey := filepath.Join(na.ManagerRootCertPath, "join/client.key")

	// clear old link
	if !global.Exists(registerCACert, registerClientCert, registerClientKey, joinCACert, joinClienCert, joinClienKey) {
		return fmt.Errorf("bad token to TLS bootstrap for network-manager")
	}

	if err := os.RemoveAll(na.ManagerCertDir); err != nil {
		return err
	}

	if err := os.MkdirAll(na.RegisterManagerCertDir, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(na.JoinManagerCertDir, os.ModePerm); err != nil {
		return err
	}

	// create cert symlink for Register
	if err := os.Symlink(registerCACert, na.RegisterCACert); err != nil {
		return fmt.Errorf("fail to create link for network-manager certificat err:%s", err.Error())
	}
	if err := os.Symlink(registerClientCert, na.RegisterClientCert); err != nil {
		return fmt.Errorf("fail to create link for network-manager certificat err:%s", err.Error())
	}
	if err := os.Symlink(registerClientKey, na.RegisterClientkey); err != nil {
		return fmt.Errorf("fail to create link for network-manager certificat err:%s", err.Error())
	}

	// create cert symlink for join
	if err := os.Symlink(joinCACert, na.JoinCACert); err != nil {
		return fmt.Errorf("fail to create link for network-manager certificat err:%s", err.Error())
	}
	if err := os.Symlink(joinClienCert, na.JoinClientCert); err != nil {
		return fmt.Errorf("fail to create link for network-manager certificat err:%s", err.Error())
	}
	if err := os.Symlink(joinClienKey, na.JoinClientkey); err != nil {
		return fmt.Errorf("fail to create link for network-manager certificat err:%s", err.Error())
	}

	if err := ioutil.WriteFile(na.NodeTokenPath, []byte(global.ReservedNodeToken), os.FileMode(0644)); err != nil {
		return fmt.Errorf("fail to create node token")
	}

	if !na.Check() {
		return fmt.Errorf("fail to create symlink for network-manager certificate")
	}

	return nil
}
