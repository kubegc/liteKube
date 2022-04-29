package grpc_client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/Litekube/network-controller/grpc/pb_gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"path/filepath"
)

type GrpcBootStrapClient struct {
	BootstrapC    pb_gen.LiteKubeNCBootstrapServiceClient
	Ip            string
	BootstrapPort string
	GrpcCertDir   string
	CAFile        string
	CertFile      string
	KeyFile       string
}

func (c *GrpcBootStrapClient) InitGrpcBootstrapClientConn() error {
	// Set up a connection to the server.
	var bootAddress string
	if len(c.Ip) == 0 || len(c.BootstrapPort) == 0 {
		logger.Error("ip and port can't be empty")
		return errors.New("ip and port can't be empty")
	}
	bootAddress = fmt.Sprintf("%s:%s", c.Ip, c.BootstrapPort)

	var dialOpt []grpc.DialOption
	var creds credentials.TransportCredentials
	if c.GrpcCertDir != "" && c.CertFile != "" && c.KeyFile != "" && c.CAFile != "" {
		cert, err := tls.LoadX509KeyPair(filepath.Join(c.GrpcCertDir, c.CertFile), filepath.Join(c.GrpcCertDir, c.KeyFile))
		if err != nil {
			logger.Errorf("tls.LoadX509KeyPair err: %v", err)
			return err
		}

		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(filepath.Join(c.GrpcCertDir, c.CAFile))
		if err != nil {
			logger.Errorf("ioutil.ReadFile err: %v", err)
			return err
		}
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			logger.Errorf("certPool.AppendCertsFromPEM err")
			return err
		}
		creds = credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ServerName:   c.Ip,
			RootCAs:      certPool,
		})
	} else {
		creds = credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})
	}
	dialOpt = append(dialOpt, []grpc.DialOption{grpc.WithTransportCredentials(creds)}...)

	bootConn, err := grpc.Dial(bootAddress, dialOpt...)
	if err != nil {
		logger.Errorf("can't connect: %v", err)
		return err
	}
	c.BootstrapC = pb_gen.NewLiteKubeNCBootstrapServiceClient(bootConn)
	return nil
}
