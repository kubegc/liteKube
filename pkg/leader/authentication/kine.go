package authentication

import (
	"fmt"
	"net"
	"path/filepath"

	"github.com/litekube/LiteKube/pkg/certificate"
	"github.com/litekube/LiteKube/pkg/global"
	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
	kineoptions "github.com/litekube/LiteKube/pkg/options/leader/kine"
	"github.com/rancher/dynamiclistener/cert"
)

type KineAuthentication struct {
	KineCertDir string
	BindAddress string
	CACert      string
	CAKey       string
	ServerCert  string
	Serverkey   string
	ClientCert  string
	Clientkey   string
}

func NewKineAuthentication(rootCertPath string, bindAddress string) *KineAuthentication {
	if rootCertPath == "" {
		rootCertPath = filepath.Join(globaloptions.DefaultGO.WorkDir, "tls/")
	}

	kineCertPath := filepath.Join(rootCertPath, "kine")

	// check bind-address
	if ip := net.ParseIP(bindAddress); ip == nil {
		bindAddress = kineoptions.DefaultKO.BindAddress
	}

	return &KineAuthentication{
		KineCertDir: kineCertPath,
		BindAddress: bindAddress,
		CACert:      filepath.Join(kineCertPath, "ca.crt"),
		CAKey:       filepath.Join(kineCertPath, "ca.key"),
		ServerCert:  filepath.Join(kineCertPath, "server.crt"),
		Serverkey:   filepath.Join(kineCertPath, "server.key"),
		ClientCert:  filepath.Join(kineCertPath, "client.crt"),
		Clientkey:   filepath.Join(kineCertPath, "client.key"),
	}
}

func (kine *KineAuthentication) GenerateOrSkip() error {
	if kine == nil {
		return fmt.Errorf("nil kine")
	}

	// generate CA
	regen, err := certificate.GenerateSigningCertKey(false, "kine", kine.CACert, kine.CAKey)
	if err != nil {
		return err
	}

	// generate server
	if _, err := certificate.GenerateServerCertKey(regen, "kine-server", nil,
		&cert.AltNames{
			DNSNames: append(kine.QueryRemoteDNSNames(), global.LocalHostDNSName),
			IPs:      append(global.LocalIPs, kine.QueryRemoteIps()...),
		}, kine.CACert, kine.CAKey, kine.ServerCert, kine.Serverkey); err != nil {
		return err
	}

	// generate client
	if _, err := certificate.GenerateClientCertKey(regen, "kine-client", []string{"kine:client"}, kine.CACert, kine.CAKey, kine.ClientCert, kine.Clientkey); err != nil {
		return err
	}

	return nil
}

func (kine *KineAuthentication) QueryRemoteIps() []net.IP {
	return []net.IP{}
}

func (kine *KineAuthentication) QueryRemoteDNSNames() []string {
	return []string{}
}
