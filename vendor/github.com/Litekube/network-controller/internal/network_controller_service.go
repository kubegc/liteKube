package internal

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/Litekube/network-controller/certs"
	"github.com/Litekube/network-controller/config"
	"github.com/Litekube/network-controller/contant"
	"github.com/Litekube/network-controller/grpc/pb_gen"
	"github.com/Litekube/network-controller/sqlite"
	"github.com/Litekube/network-controller/utils"
	certutil "github.com/rancher/dynamiclistener/cert"
)

type NetworkControllerService struct {
	unRegisterCh     chan string
	grpcTlsConfig    config.TLSConfig
	networkTlsConfig config.TLSConfig
	bootstrapIp      string
	port             string
}

var logger = utils.GetLogger()

func NewLiteNCService(unRegisterCh chan string, grpcTlsConfig config.TLSConfig, networkTlsConfig config.TLSConfig, bootstrapIp, port string) *NetworkControllerService {
	return &NetworkControllerService{
		unRegisterCh:     unRegisterCh,
		grpcTlsConfig:    grpcTlsConfig,
		networkTlsConfig: networkTlsConfig,
		bootstrapIp:      bootstrapIp,
		port:             port,
	}
}

func (service *NetworkControllerService) GetBootStrapToken(ctx context.Context, req *pb_gen.GetBootStrapTokenRequest) (*pb_gen.GetBootStrapTokenResponse, error) {

	wrappedResp := func(code, message, token string) (resp *pb_gen.GetBootStrapTokenResponse, err error) {
		if code != contant.STATUS_OK {
			err = errors.New(message)
		}
		resp = &pb_gen.GetBootStrapTokenResponse{
			Code:           code,
			Message:        message,
			BootStrapToken: token,
			CloudIp:        service.bootstrapIp,
			Port:           service.port,
		}
		logger.Debugf("resp: %+v", resp)
		return
	}

	if req.ExpireTime == 0 {
		req.ExpireTime = contant.IdleTokenExpireDuration
	}
	token := utils.GetUniqueToken()
	tm := sqlite.TokenMgr{}
	// no need
	//item, err := nm.QueryByToken(token)
	err := tm.Insert(sqlite.TokenMgr{
		Token: token,
		//ExpireTime: time.Now().Add(time.Duration(req.ExpireTime) * time.Minute),
	}, req.ExpireTime)
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}

	return wrappedResp(contant.STATUS_OK, contant.MESSAGE_OK, token)
}

func (service *NetworkControllerService) CheckConnState(ctx context.Context, req *pb_gen.CheckConnStateRequest) (*pb_gen.CheckConnResponse, error) {

	wrappedResp := func(code, message, bindIp string, state int32) (resp *pb_gen.CheckConnResponse, err error) {
		if code != contant.STATUS_OK {
			logger.Errorf("query token: %+v err: %+v", req.Token, err)
			err = errors.New(message)
		}
		resp = &pb_gen.CheckConnResponse{
			Message:   message,
			Code:      code,
			ConnState: state,
			BindIp:    bindIp,
		}
		logger.Debugf("resp: %+v", resp)
		return
	}

	if len(req.Token) == 0 {
		return wrappedResp(contant.STATUS_BADREQUEST, "token can't be empty", "", -1)
	}

	nm := sqlite.NetworkMgr{}
	item, err := nm.QueryByToken(req.Token)
	if item == nil {
		return wrappedResp(contant.STATUS_OK, err.Error(), "", -1)
	} else if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "", -1)
	}

	return wrappedResp(contant.STATUS_OK, contant.MESSAGE_OK, item.BindIp, int32(item.State))
}

func (service *NetworkControllerService) UnRegister(ctx context.Context, req *pb_gen.UnRegisterRequest) (*pb_gen.UnRegisterResponse, error) {

	wrappedResp := func(code, message string, result bool) (resp *pb_gen.UnRegisterResponse, err error) {
		if code != contant.STATUS_OK {
			logger.Errorf("query token: %+v err: %+v", req.Token, err)
			err = errors.New(message)
		}
		resp = &pb_gen.UnRegisterResponse{
			Message: message,
			Code:    code,
			Result:  result,
		}
		logger.Debugf("resp: %+v", resp)
		return
	}

	if len(req.Token) == 0 {
		return wrappedResp(contant.STATUS_BADREQUEST, "token can't be empty", false)
	}

	nm := sqlite.NetworkMgr{}
	item, err := nm.QueryByToken(req.Token)
	if item == nil {
		return wrappedResp(contant.STATUS_OK, err.Error(), false)
	} else if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), false)
	}

	result, err := nm.DeleteById(item.Id)
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), result)
	}

	service.unRegisterCh <- item.BindIp
	return wrappedResp(contant.STATUS_OK, contant.MESSAGE_OK, result)
}

func (service *NetworkControllerService) GetRegistedIp(ctx context.Context, req *pb_gen.GetRegistedIpRequest) (*pb_gen.GetRegistedIpResponse, error) {

	wrappedResp := func(code, message, ip string) (resp *pb_gen.GetRegistedIpResponse, err error) {
		if code != contant.STATUS_OK {
			logger.Errorf("query token: %+v err: %+v", req.Token, err)
			err = errors.New(message)
		}
		resp = &pb_gen.GetRegistedIpResponse{
			Message: message,
			Code:    code,
			Ip:      ip,
		}
		logger.Debugf("resp: %+v", resp)
		return
	}

	if len(req.Token) == 0 {
		return wrappedResp(contant.STATUS_BADREQUEST, "token can't be empty", "")
	}

	nm := sqlite.NetworkMgr{}
	item, err := nm.QueryByToken(req.Token)
	if item == nil {
		return wrappedResp(contant.STATUS_OK, err.Error(), "")
	} else if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}

	return wrappedResp(contant.STATUS_OK, contant.MESSAGE_OK, item.BindIp)
}

func (service *NetworkControllerService) GetToken(ctx context.Context, req *pb_gen.GetTokenRequest) (*pb_gen.GetTokenResponse, error) {

	wrappedResp := func(code, message, token string) (resp *pb_gen.GetTokenResponse, err error) {
		if code != contant.STATUS_OK {
			err = errors.New(message)
		}
		resp = &pb_gen.GetTokenResponse{
			Code:              code,
			Message:           message,
			Token:             token,
			GrpcCaCert:        "",
			GrpcClientKey:     "",
			GrpcClientCert:    "",
			NetworkCaCert:     "",
			NetworkClientKey:  "",
			NetworkClientCert: "",
		}
		logger.Debugf("resp: %+v", resp)
		return
	}

	if len(req.BootStrapToken) == 0 {
		return wrappedResp(contant.STATUS_BADREQUEST, "bootstrap token can't be empty", "")
	}

	tm := sqlite.TokenMgr{}
	item, err := tm.QueryByToken(req.BootStrapToken)
	if item == nil {
		return wrappedResp(contant.STATUS_OK, err.Error(), "")
	} else if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}

	token := utils.GetUniqueToken()
	nm := sqlite.NetworkMgr{}
	// no need
	//item, err := nm.QueryByToken(token)
	err = nm.Insert(sqlite.NetworkMgr{
		Token:  token,
		State:  contant.STATE_IDLE,
		BindIp: "",
	})
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}

	keyBytes, certBytes, _, err := certs.GenerateClientCertKey(true, "network-controller-grpc-client", []string{"network-controller-grpc"}, service.grpcTlsConfig.CAFile, service.grpcTlsConfig.CAKeyFile, service.grpcTlsConfig.ClientCertFile, service.grpcTlsConfig.ClientKeyFile)
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}

	resp, _ := wrappedResp(contant.STATUS_OK, contant.MESSAGE_OK, token)

	// load grpc ca.pem client.pem client-key.pem
	grpcCaCert, err := certs.LoadCertificate(service.grpcTlsConfig.CAFile)
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}
	resp.GrpcCaCert = base64.StdEncoding.EncodeToString(certutil.EncodeCertPEM(grpcCaCert))
	resp.GrpcClientKey = base64.StdEncoding.EncodeToString(keyBytes)
	resp.GrpcClientCert = base64.StdEncoding.EncodeToString(certBytes)

	keyBytes, certBytes, _, err = certs.GenerateClientCertKey(true, "network-controller-client", []string{"network-controller"}, service.networkTlsConfig.CAFile, service.networkTlsConfig.CAKeyFile, service.networkTlsConfig.ClientCertFile, service.networkTlsConfig.ClientKeyFile)
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}

	// load network ca.pem client.pem client-key.pem
	NetworkCaCert, err := certs.LoadCertificate(service.networkTlsConfig.CAFile)
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}
	resp.NetworkCaCert = base64.StdEncoding.EncodeToString(certutil.EncodeCertPEM(NetworkCaCert))
	resp.NetworkClientKey = base64.StdEncoding.EncodeToString(keyBytes)
	resp.NetworkClientCert = base64.StdEncoding.EncodeToString(certBytes)
	return resp, nil
}
