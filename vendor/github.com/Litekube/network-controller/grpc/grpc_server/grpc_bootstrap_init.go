package grpc_server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Litekube/network-controller/grpc/pb_gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"net"
)

func (s *GrpcServer) StartBootstrapServerTcp() error {
	defer logger.Debug("StartBootstrapServerTcp done")

	tcpAddr := fmt.Sprintf(":%d", s.bootstrapPort)
	lis, err := net.Listen("tcp", tcpAddr)
	defer lis.Close()
	if err != nil {
		logger.Errorf("tcp failed to listen: %v", err)
		return err
	}

	gopts := []grpc.ServerOption{}
	if len(s.grpcTlsConfig.ServerCertFile) != 0 && len(s.grpcTlsConfig.ServerKeyFile) != 0 {
		creds, err := credentials.NewServerTLSFromFile(s.grpcTlsConfig.ServerCertFile, s.grpcTlsConfig.ServerKeyFile)
		if err != nil {
			logger.Error(err)
			return err
		}
		gopts = append(gopts, grpc.Creds(creds))
	}
	cert, err := tls.LoadX509KeyPair(s.grpcTlsConfig.ServerCertFile, s.grpcTlsConfig.ServerKeyFile)
	//cert, err := certificate.LoadCertificate(s.CertFile)
	if err != nil {
		logger.Errorf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(s.grpcTlsConfig.CAFile)
	if err != nil {
		logger.Errorf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		logger.Errorf("certPool.AppendCertsFromPEM err")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		// fix here
		ClientAuth: tls.VerifyClientCertIfGiven,
		ClientCAs:  certPool,
	})

	gopts = append(gopts, []grpc.ServerOption{grpc.Creds(creds)}...)
	server := grpc.NewServer(gopts...)
	// register reflection for grpcurl service
	reflection.Register(server)
	// register service
	pb_gen.RegisterLiteKubeNCBootstrapServiceServer(server, s)
	logger.Infof("grpc server bootstrap ready to serve at %+v", tcpAddr)

	go func() {
		for {
			select {
			case <-s.stopCh:
				server.GracefulStop()
				return
			}
		}
	}()

	if err := server.Serve(lis); err != nil {
		logger.Errorf("grpc server failed to serve: %v", err)
		return err
	}
	return nil
}
