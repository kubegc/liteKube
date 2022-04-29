package certs

import (
	"github.com/Litekube/network-controller/config"
	"github.com/Litekube/network-controller/utils"
	"github.com/litekube/LiteKube/pkg/global"
	"github.com/rancher/dynamiclistener/cert"
	"net"
)

func CheckGrpcCertConfig(tlsConfig config.TLSConfig) error {
	// generate for grpc
	// generate CA
	regenGrpc, err := GenerateSigningCertKey(false, "network-controller-grpc", tlsConfig.CAFile, tlsConfig.CAKeyFile)
	if err != nil {
		return err
	}

	// generate server
	if _, _, _, err := GenerateServerCertKey(regenGrpc, "network-controller-grpc-server", nil,
		&cert.AltNames{
			DNSNames: append([]string{}, global.LocalHostDNSName),
			IPs:      append(append(global.LocalIPs, []net.IP{net.ParseIP(utils.QueryPublicIp())}...)),
		}, tlsConfig.CAFile, tlsConfig.CAKeyFile, tlsConfig.ServerCertFile, tlsConfig.ServerKeyFile); err != nil {
		return err
	}

	// generate client
	if _, _, _, err := GenerateClientCertKey(regenGrpc, "network-controller-grpc-client", []string{"network-controller-grpc"}, tlsConfig.CAFile, tlsConfig.CAKeyFile, tlsConfig.ClientCertFile, tlsConfig.ClientKeyFile); err != nil {
		return err
	}
	return nil
}

func CheckNetworkCertConfig(tlsConfig config.TLSConfig) error {
	//generate for network
	//generate CA
	regenGrpc, err := GenerateSigningCertKey(false, "network-controller", tlsConfig.CAFile, tlsConfig.CAKeyFile)
	if err != nil {
		return err
	}

	// generate server
	if _, _, _, err := GenerateServerCertKey(regenGrpc, "network-controller-server", nil,
		&cert.AltNames{
			DNSNames: append([]string{}, global.LocalHostDNSName),
			IPs:      append(append(global.LocalIPs, []net.IP{net.ParseIP(utils.QueryPublicIp())}...)),
		}, tlsConfig.CAFile, tlsConfig.CAKeyFile, tlsConfig.ServerCertFile, tlsConfig.ServerKeyFile); err != nil {
		return err
	}

	// generate client
	if _, _, _, err := GenerateClientCertKey(regenGrpc, "network-controller-client", []string{"network-controller"}, tlsConfig.CAFile, tlsConfig.CAKeyFile, tlsConfig.ClientCertFile, tlsConfig.ClientKeyFile); err != nil {
		return err
	}
	return nil
}
