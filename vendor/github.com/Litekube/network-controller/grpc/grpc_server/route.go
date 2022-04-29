package grpc_server

import (
	"context"
	"github.com/Litekube/network-controller/grpc/pb_gen"
)

func (s *GrpcServer) HelloWorld(ctx context.Context, req *pb_gen.HelloWorldRequest) (*pb_gen.HelloWorldResponse, error) {
	logger.Infof("get HelloWorld request: %+v", req)
	reply := &pb_gen.HelloWorldResponse{ThanksText: "hello,this wanna"}
	return reply, nil
}

func (s *GrpcServer) CheckConnState(ctx context.Context, req *pb_gen.CheckConnStateRequest) (*pb_gen.CheckConnResponse, error) {
	logger.Infof("get CheckConnState request: %+v", req)
	resp, _ := s.service.CheckConnState(ctx, req)
	return resp, nil
}

func (s *GrpcServer) UnRegister(ctx context.Context, req *pb_gen.UnRegisterRequest) (*pb_gen.UnRegisterResponse, error) {
	logger.Infof("get UnRegister request: %+v", req)
	resp, _ := s.service.UnRegister(ctx, req)
	return resp, nil
}

func (s *GrpcServer) GetRegistedIp(ctx context.Context, req *pb_gen.GetRegistedIpRequest) (*pb_gen.GetRegistedIpResponse, error) {
	logger.Infof("get GetRegistedIp request: %+v", req)
	resp, _ := s.service.GetRegistedIp(ctx, req)
	return resp, nil
}

func (s *GrpcServer) GetToken(ctx context.Context, req *pb_gen.GetTokenRequest) (*pb_gen.GetTokenResponse, error) {
	logger.Infof("get GetToken request: %+v", req)
	resp, _ := s.service.GetToken(ctx, req)
	return resp, nil
}

func (s *GrpcServer) GetBootStrapToken(ctx context.Context, req *pb_gen.GetBootStrapTokenRequest) (*pb_gen.GetBootStrapTokenResponse, error) {
	logger.Infof("get GetBootStrapToken request: %+v", req)
	resp, _ := s.service.GetBootStrapToken(ctx, req)
	return resp, nil
}
