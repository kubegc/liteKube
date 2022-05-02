package authentication

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/litekube/LiteKube/pkg/certificate"
	"github.com/litekube/LiteKube/pkg/global"
	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
	"github.com/rancher/dynamiclistener/cert"
)

type NetworkControllerAuthentication struct {
	ManagerCertDir         string
	NetworkControllerDir   string
	RegisterManagerCertDir string
	RegisterBindAddress    string
	RegisterCACert         string
	RegisterCAKey          string
	RegisterServerCert     string
	RegisterServerkey      string
	RegisterClientCert     string
	RegisterClientkey      string

	JoinManagerCertDir string
	JoinBindAddress    string
	JoinCACert         string
	JoinCAKey          string
	JoinServerCert     string
	JoinServerkey      string
	JoinClientCert     string
	JoinClientkey      string
}

func NewNetworkControllerAuthentication(workDir string, rootCertPath string, registerBindAddress, joinBindAddress string) *NetworkControllerAuthentication {
	if rootCertPath == "" {
		rootCertPath = filepath.Join(globaloptions.DefaultGO.WorkDir, "tls/")
	}

	if workDir == "" {
		workDir = globaloptions.DefaultGO.WorkDir
	}

	networkControllerDir := filepath.Join(workDir, "network-controller")
	managerCertDir := filepath.Join(rootCertPath, "network-controller")
	registerManagerCertDir := filepath.Join(managerCertDir, "register")
	joinManagerCertDir := filepath.Join(managerCertDir, "join")

	if err := os.MkdirAll(networkControllerDir, os.FileMode(0644)); err != nil {
		panic(err)
	}

	// check bind-address
	if ip := net.ParseIP(registerBindAddress); ip == nil || registerBindAddress == "127.0.0.1" {
		registerBindAddress = "0.0.0.0"
	}

	if ip := net.ParseIP(joinBindAddress); ip == nil || joinBindAddress == "127.0.0.1" {
		joinBindAddress = "0.0.0.0"
	}

	return &NetworkControllerAuthentication{
		NetworkControllerDir:   networkControllerDir,
		ManagerCertDir:         managerCertDir,
		RegisterManagerCertDir: registerManagerCertDir,
		RegisterBindAddress:    registerBindAddress,
		RegisterCACert:         filepath.Join(registerManagerCertDir, "ca.crt"),
		RegisterCAKey:          filepath.Join(registerManagerCertDir, "ca.key"),
		RegisterServerCert:     filepath.Join(registerManagerCertDir, "server.crt"),
		RegisterServerkey:      filepath.Join(registerManagerCertDir, "server.key"),
		RegisterClientCert:     filepath.Join(registerManagerCertDir, "client.crt"),
		RegisterClientkey:      filepath.Join(registerManagerCertDir, "client.key"),
		JoinBindAddress:        joinBindAddress,
		JoinCACert:             filepath.Join(joinManagerCertDir, "ca.crt"),
		JoinCAKey:              filepath.Join(joinManagerCertDir, "ca.key"),
		JoinServerCert:         filepath.Join(joinManagerCertDir, "server.crt"),
		JoinServerkey:          filepath.Join(joinManagerCertDir, "server.key"),
		JoinClientCert:         filepath.Join(joinManagerCertDir, "client.crt"),
		JoinClientkey:          filepath.Join(joinManagerCertDir, "client.key"),
	}
}

// generate X.509 certificate for network-manager
func (na *NetworkControllerAuthentication) GenerateOrSkip() error {
	if na == nil {
		return fmt.Errorf("nil network authentication")
	}

	// generate for register
	// generate CA
	regenRegister, err := certificate.GenerateSigningCertKey(false, "lknm-register", na.RegisterCACert, na.RegisterCAKey)
	if err != nil {
		return err
	}

	// generate server
	if _, err := certificate.GenerateServerCertKey(regenRegister, "register-server", nil,
		&cert.AltNames{
			DNSNames: append(na.QueryRemoteDNSNames(), global.LocalHostDNSName),
			IPs:      global.RemoveRepeatIps(append(append(global.LocalIPs, []net.IP{net.ParseIP(na.RegisterBindAddress), net.ParseIP(na.JoinBindAddress)}...), na.QueryRemoteIps()...)),
		}, na.RegisterCACert, na.RegisterCAKey, na.RegisterServerCert, na.RegisterServerkey); err != nil {
		return err
	}

	// generate client
	if _, err := certificate.GenerateClientCertKey(regenRegister, "register-client", []string{"lknm:register"}, na.RegisterCACert, na.RegisterCAKey, na.RegisterClientCert, na.RegisterClientkey); err != nil {
		return err
	}

	// generate for join
	// generate CA
	regenJoin, err := certificate.GenerateSigningCertKey(false, "lknm-join", na.JoinCACert, na.JoinCAKey)
	if err != nil {
		return err
	}

	// generate server
	if _, err := certificate.GenerateServerCertKey(regenJoin, "join-server", nil,
		&cert.AltNames{
			DNSNames: append(na.QueryRemoteDNSNames(), global.LocalHostDNSName),
			IPs:      append(append(global.LocalIPs, []net.IP{net.ParseIP(na.RegisterBindAddress), net.ParseIP(na.JoinBindAddress)}...), na.QueryRemoteIps()...),
		}, na.JoinCACert, na.JoinCAKey, na.JoinServerCert, na.JoinServerkey); err != nil {
		return err
	}

	// generate client
	if _, err := certificate.GenerateClientCertKey(regenJoin, "join-client", []string{"lknm:join"}, na.JoinCACert, na.JoinCAKey, na.JoinClientCert, na.JoinClientkey); err != nil {
		return err
	}
	return nil
}

func (na *NetworkControllerAuthentication) QueryRemoteIps() []net.IP {
	networkControllerServerIp := global.GetDefaultServiceIp(global.NetworkControllerCIDR)
	if networkControllerServerIp != nil {
		return []net.IP{networkControllerServerIp}
	}
	return []net.IP{}
}

func (na *NetworkControllerAuthentication) QueryRemoteDNSNames() []string {
	return []string{}
}
