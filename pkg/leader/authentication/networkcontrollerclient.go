package authentication

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Litekube/network-controller/grpc/grpc_client"
	"github.com/Litekube/network-controller/grpc/pb_gen"
	certutil "github.com/rancher/dynamiclistener/cert"

	"github.com/litekube/LiteKube/pkg/certificate"
	"github.com/litekube/LiteKube/pkg/global"
	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

type NetworkControllerClientAuthentication struct {
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
	NodeToken              string
	InfoPath               string
}

type RemoteHostInfo struct {
	RegisterAddress string // value only tls-bootstrap without init
	RegisterPort    uint16 // value only tls-bootstrap without init
	JoinAddress     string // value only tls-bootstrap without init
	JoinPort        uint16 // value only tls-bootstrap without init
	NodeToken       string
}

func NewControllerClientAuthentication(rootCertPath string, token string, registerAddress *string, registerPort *uint16, joinAddress *string, joinPort *uint16) *NetworkControllerClientAuthentication {
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

	return &NetworkControllerClientAuthentication{
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
		NodeToken:              "",
		InfoPath:               filepath.Join(managerCertDir, "info.yaml"),
	}
}

func (na *NetworkControllerClientAuthentication) LoadInfo() error {
	if !global.Exists(na.InfoPath) {
		return fmt.Errorf("info file not exist")
	}

	if bytes, err := ioutil.ReadFile(na.InfoPath); err != nil {
		return err
	} else {
		data := RemoteHostInfo{}
		if err := yaml.Unmarshal(bytes, &data); err != nil {
			return err
		} else {
			if data.NodeToken != "" {
				na.NodeToken = data.NodeToken
			} else {
				return fmt.Errorf("fail to load network controller server key information")
			}

			if data.RegisterAddress != "" {
				*na.RegisterAddress = data.RegisterAddress
			}
			if data.RegisterPort != 0 {
				*na.RegisterPort = data.RegisterPort
			}
			if data.JoinAddress != "" {
				*na.JoinAddress = data.JoinAddress
			}
			if data.JoinPort != 0 {
				*na.JoinPort = data.JoinPort
			}
		}
	}

	return nil
}

func (na *NetworkControllerClientAuthentication) Nodetoken() (string, error) {
	return na.NodeToken, nil
}

// generate X.509 certificate for network-manager
func (na *NetworkControllerClientAuthentication) GenerateOrSkip() error {
	if na.Token == "unknown" {
		return fmt.Errorf("token is unknown")
	}

	if na == nil {
		return fmt.Errorf("nil network authentication")
	}

	if na.Check() {
		// file exist, bootstrap ok.
		return nil
	}

	if na.Token == "local" {
		return na.CreatelinkForClient()
	} else {
		token := strings.SplitN(na.Token, "@", 2)[0]
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
		return na.TLSBootStrap(Endpoint[0], port, token)
	}
}

// download certificates and get node-token from network manager
func (na *NetworkControllerClientAuthentication) TLSBootStrap(address string, port int, bootstrapToken string) error {
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

	bootClient := &grpc_client.GrpcBootStrapClient{
		Ip:            address,
		BootstrapPort: strconv.FormatUint(uint64(port), 10),
		//BootstrapPort: "6439",
	}

	if bootClient.BootstrapC == nil {
		if err := bootClient.InitGrpcBootstrapClientConn(); err != nil {
			panic(err)
		}
	}

	// start in 5s
	for i := 0; i < 10; i++ {
		resp, err := bootClient.BootstrapC.HealthCheck(context.Background(), &pb_gen.HealthCheckRequest{})
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if resp.Code == "200" {
			break
		}
		if i == 9 {
			panic(err)
		}
	}

	req := &pb_gen.GetTokenRequest{
		BootStrapToken: bootstrapToken,
	}
	resp, err := bootClient.BootstrapC.GetToken(context.Background(), req)
	if err != nil {
		return err
	}

	// assign value address and port
	grpcPort, _ := strconv.ParseUint(resp.GrpcServerPort, 10, 16)
	networkPort, _ := strconv.ParseUint(resp.NetworkServerPort, 10, 16)
	*na.RegisterAddress = resp.GrpcServerIp
	*na.RegisterPort = uint16(grpcPort)
	*na.JoinAddress = resp.NetworkServerIp
	*na.JoinPort = uint16(networkPort)

	if bytes, err := yaml.Marshal(RemoteHostInfo{
		RegisterAddress: *na.RegisterAddress,
		RegisterPort:    *na.RegisterPort,
		JoinAddress:     *na.JoinAddress,
		JoinPort:        *na.JoinPort,
		NodeToken:       resp.Token,
	}); err != nil {
		return fmt.Errorf("fail to marshal host info")
	} else {
		if err := ioutil.WriteFile(na.InfoPath, bytes, os.FileMode(0644)); err != nil {
			return fmt.Errorf("fail to create node token")
		}
	}

	// register cert file
	caBytes, _ := base64.StdEncoding.DecodeString(resp.GrpcCaCert)
	certBytes, _ := base64.StdEncoding.DecodeString(resp.GrpcClientCert)
	keyBytes, _ := base64.StdEncoding.DecodeString(resp.GrpcClientKey)
	certutil.WriteCert(na.RegisterCACert, caBytes)
	certutil.WriteCert(na.RegisterClientCert, certBytes)
	certutil.WriteKey(na.RegisterClientkey, keyBytes)

	// join cert file
	caBytes, _ = base64.StdEncoding.DecodeString(resp.NetworkCaCert)
	certBytes, _ = base64.StdEncoding.DecodeString(resp.NetworkClientCert)
	keyBytes, _ = base64.StdEncoding.DecodeString(resp.NetworkClientKey)
	certutil.WriteCert(na.JoinCACert, caBytes)
	certutil.WriteCert(na.JoinClientCert, certBytes)
	certutil.WriteKey(na.JoinClientkey, keyBytes)

	return nil
}

func (na *NetworkControllerClientAuthentication) Check() bool {
	if !certificate.Exists(na.RegisterCACert, na.RegisterClientCert, na.RegisterClientkey, na.JoinCACert, na.JoinClientCert, na.JoinClientkey, na.InfoPath) {
		return false
	}

	if err := na.LoadInfo(); err != nil {
		klog.Warningf("%v", err)
		return false
	}

	return true
}

func (na *NetworkControllerClientAuthentication) CreatelinkForClient() error {
	registerCACert := filepath.Join(na.ManagerRootCertPath, "register/ca.crt")
	registerClientCert := filepath.Join(na.ManagerRootCertPath, "register/client.crt")
	registerClientKey := filepath.Join(na.ManagerRootCertPath, "register/client.key")

	joinCACert := filepath.Join(na.ManagerRootCertPath, "join/ca.crt")
	joinClienCert := filepath.Join(na.ManagerRootCertPath, "join/client.crt")
	joinClienKey := filepath.Join(na.ManagerRootCertPath, "join/client.key")

	// clear old link
	if !global.Exists(registerCACert, registerClientCert, registerClientKey, joinCACert, joinClienCert, joinClienKey) {
		return fmt.Errorf("TLS bootstrap for network-manager set token='local' only be allowed while worker run in leader to ")
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

	if bytes, err := yaml.Marshal(RemoteHostInfo{
		RegisterAddress: *na.RegisterAddress,
		RegisterPort:    *na.RegisterPort,
		JoinAddress:     *na.JoinAddress,
		JoinPort:        *na.JoinPort,
		NodeToken:       global.ReservedNodeToken,
	}); err != nil {
		return fmt.Errorf("fail to marshal host info")
	} else {
		if err := ioutil.WriteFile(na.InfoPath, bytes, os.FileMode(0644)); err != nil {
			return fmt.Errorf("fail to create node token")
		}
	}

	if !na.Check() {
		return fmt.Errorf("fail to create symlink for network-manager certificate or record info")
	}

	return nil
}
