package runtime

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/litekube/LiteKube/pkg/certificate"
	"github.com/litekube/LiteKube/pkg/global"
	"github.com/litekube/LiteKube/pkg/leader/runtime/control"
	"github.com/litekube/LiteKube/pkg/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

var (
	globalFold                            = filepath.Join(global.HomePath, ".litekube")
	globalToken                           = filepath.Join(globalFold, "token")
	kubelet_bootstrap_kubeconfig_template = template.Must(template.New("kubelet_bootstrap_kubeconfig").Parse(`apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: {{.CACert}}
    server: {{.URL}}
  name: litekube
contexts:
- context:
    cluster: litekube
    user: kubelet-bootstrap
  name: default
current-context: default
kind: Config
preferences: {}
users:
- name: kubelet-bootstrap
  user:
    token: {{.Token}}
`))

	kubeproxy_bootstrap_kubeconfig_template = template.Must(template.New("kubeproxy_kubeconfig").Parse(`apiVersion: v1
clusters:
- cluster:
	server: {{.URL}}
	certificate-authority: {{.CACert}}
	name: litekube
contexts:
- context:
	cluster: litekube
	namespace: default
	user: kubeproxy
	name: Default
current-context: Default
kind: Config
preferences: {}
users:
- name: kubeproxy
	user:
	client-certificate: {{.ClientCert}}
	client-key: {{.ClientKey}}
`))
)

type LiteKubeControl struct {
	ctx                context.Context
	BindAddress        string
	BindPort           uint16
	server             *grpc.Server
	NetworkClient      *NetWorkRegisterClient
	AuthFile           string
	LocalHostNodeToken string

	BufferPath                string
	EndPoint                  string
	ValidateApiserverServerCA string
	ValidateApiserverClientCA string
	SignalApiserverClientCert string
	SignalApiserverClientKey  string
	AuthTokenFile             string
	ClusterCIDR               string
	ServiceClusterIpRange     string
	ClusterDNS                string

	ValidateApiserverServerCABase64 string
	ValidateApiserverClientCABase64 string
	CacheCertPath                   string
	CacheKeyPath                    string
	Token                           string
}

type TokenDesc struct {
	Token      string
	CreateBy   string
	CreateTime string
	Life       int64
	IsAdmin    bool
	IsValid    bool
}

func NewLiteKubeControl(ctx context.Context, networkClient *NetWorkRegisterClient, bufferPath string, nodeToken string, endpoint string, validateApiserverServerCA string, validateApiserverClientCA string, signalApiserverClientCert string, signalApiserverClientKey string, authTokenFile string, clusterCIDR string, serviceClusterIpRange string) *LiteKubeControl {
	tmp := &LiteKubeControl{
		ctx:                       ctx,
		BindAddress:               "0.0.0.0",
		BindPort:                  6442,
		BufferPath:                bufferPath,
		LocalHostNodeToken:        nodeToken,
		EndPoint:                  endpoint,
		NetworkClient:             networkClient,
		ValidateApiserverServerCA: validateApiserverServerCA,
		ValidateApiserverClientCA: validateApiserverClientCA,
		SignalApiserverClientCert: signalApiserverClientCert,
		SignalApiserverClientKey:  signalApiserverClientKey,
		AuthTokenFile:             authTokenFile,
		ClusterCIDR:               clusterCIDR,
		ServiceClusterIpRange:     serviceClusterIpRange,
		ClusterDNS:                "",
	}

	if err := tmp.Init(); err != nil {
		klog.Errorf(err.Error())
		return nil
	} else {
		return tmp
	}
}

func (s *LiteKubeControl) Init() error {
	if s == nil || !global.Exists(s.ValidateApiserverServerCA, s.ValidateApiserverClientCA, s.SignalApiserverClientCert, s.SignalApiserverClientKey, s.AuthTokenFile) {
		return fmt.Errorf("loss args")
	}

	if bytedatas, err := ioutil.ReadFile(s.AuthTokenFile); err != nil {
		return fmt.Errorf("fail to read token Auth File")
	} else {
		token := strings.SplitN(string(bytedatas), ",", 2)[0]
		if len(token) < 1 {
			return fmt.Errorf("bad token auth file")
		} else {
			s.Token = token
		}
	}

	_, serviceClusterIpRange, err := net.ParseCIDR(s.ServiceClusterIpRange)
	if err != nil {
		return nil
	}

	s.ClusterDNS = global.GetDefaultClusterDNSIP(serviceClusterIpRange).String()

	if data, err := certificate.LoadCertificateAsBase64(s.ValidateApiserverServerCA); err != nil {
		return err
	} else {
		s.ValidateApiserverServerCABase64 = data
	}

	if data, err := certificate.LoadCertificateAsBase64(s.ValidateApiserverClientCA); err != nil {
		return err
	} else {
		s.ValidateApiserverClientCABase64 = data
	}

	if err := os.MkdirAll(s.BufferPath, os.FileMode(0644)); err != nil {
		klog.Errorf("fail to create cache fold: %s", s.BufferPath)
		return err
	}

	s.AuthFile = filepath.Join(globalFold, ".auth.yaml")
	s.CacheCertPath = filepath.Join(s.BufferPath, "kube-proxy.crt")
	s.CacheKeyPath = filepath.Join(s.BufferPath, "kube-proxy.key")

	return ensureToken(s.AuthFile)
}

func (s *LiteKubeControl) Run() error {
	if s == nil {
		return fmt.Errorf("nil control")
	}

	if s.server != nil {
		return nil
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", s.BindPort))
	if err != nil {
		klog.Errorf("fail to listen %s", fmt.Sprintf("0.0.0.0:%d", s.BindPort))
		return fmt.Errorf("fail to listen %s", fmt.Sprintf(":%d", s.BindPort))
	}

	s.server = grpc.NewServer(grpc.UnaryInterceptor(s.TokenInterceptor()))
	control.RegisterLeaderControlServer(s.server, s)
	reflection.Register(s.server)

	signal := make(chan struct{})
	go func() {
		klog.Info("start litekube control")
		if err := s.server.Serve(listen); err != nil {
			klog.Errorf("==>control exit, error: %s", err.Error())
			close(signal)
		}

		select {
		case <-s.ctx.Done():
			return
		case <-signal:
			return
		}
	}()

	return nil
}

func (s *LiteKubeControl) TokenInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, exist := metadata.FromIncomingContext(ctx)
		if !exist {
			return nil, status.Errorf(codes.Unauthenticated, "need authentication info")
		}

		if tokens, ok := md["token"]; ok {
			tokenDesc := queryToken(tokens[0], s.AuthFile)
			if tokenDesc == nil {
				return nil, status.Errorf(codes.Unauthenticated, "Unauthorized info")
			}

			if !tokenDesc.IsValid {
				return nil, status.Errorf(codes.Unauthenticated, "Your authentication information has expired")
			}

			if !tokenDesc.IsAdmin && info.FullMethod != "/control.LeaderControl/BootStrapKubeProxy" && info.FullMethod != "/control.LeaderControl/BootStrapKubelet" && info.FullMethod != "/control.LeaderControl/CheckHealth" {
				return nil, status.Errorf(codes.Unauthenticated, "Insufficient permission, please contact the administrator")
			}

			return handler(ctx, req)
		} else {
			return nil, status.Errorf(codes.Unauthenticated, "need authentication info")
		}
	}
}

func (s *LiteKubeControl) BootstrapValidateKubeApiserverClient(ctx context.Context, in *control.NoneValue) (*control.BootstrapValidateKubeApiserverClientResponse, error) {
	klog.Info("bootstrap ca info to validata kube-apiserver client")
	if s.ValidateApiserverClientCABase64 != "" {
		return &control.BootstrapValidateKubeApiserverClientResponse{StatusCode: int32(http.StatusOK), Message: "success", Certificate: s.ValidateApiserverClientCABase64}, nil
	} else {
		return &control.BootstrapValidateKubeApiserverClientResponse{StatusCode: int32(http.StatusInternalServerError), Message: "fail to get info"}, nil
	}
}

func (s *LiteKubeControl) CheckHealth(ctx context.Context, in *control.NoneValue) (*control.HealthDescription, error) {
	return &control.HealthDescription{Message: "ok"}, nil
}

func (s *LiteKubeControl) NodeToken(ctx context.Context, in *control.NoneValue) (*control.TokenString, error) {
	return &control.TokenString{Token: s.LocalHostNodeToken}, nil
}

func (s *LiteKubeControl) BootStrapNetwork(ctx context.Context, in *control.BootStrapNetworkRequest) (*control.BootStrapNetworkResponse, error) {
	klog.Info("bootstrap network info")
	if token, err := s.NetworkClient.CreateBootStrapToken(in.GetLife()); err != nil {
		return &control.BootStrapNetworkResponse{StatusCode: int32(http.StatusInternalServerError), Message: "fail to create network bootstrap-token"}, nil
	} else {
		if address, err := s.NetworkClient.GetBootStrapAddress(); err != nil {
			return &control.BootStrapNetworkResponse{StatusCode: int32(http.StatusInternalServerError), Message: "fail to query network bootstrap address"}, nil
		} else {
			if port, err := s.NetworkClient.GetBootStrapPort(); err != nil {
				return &control.BootStrapNetworkResponse{StatusCode: int32(http.StatusInternalServerError), Message: "fail to query network bootstrap port"}, nil
			} else {
				return &control.BootStrapNetworkResponse{StatusCode: int32(http.StatusInternalServerError), Message: "success", Ip: address, Port: uint32(port), Token: token}, nil
			}
		}
	}
}

func (s *LiteKubeControl) CreateToken(ctx context.Context, in *control.CreateTokenRequest) (*control.TokenValue, error) {
	klog.Info("create litekube service token")
	md, exist := metadata.FromIncomingContext(ctx)
	createBy := "unknown"
	if exist && md["token"] != nil {
		createBy = md["token"][0]
	}

	tokenDesc, err := createToken(in.GetLife(), createBy, s.AuthFile, in.GetIsAdmin())
	if err != nil {
		return &control.TokenValue{StatusCode: int32(http.StatusInternalServerError), Message: "fail to create token"}, nil
	}

	return &control.TokenValue{StatusCode: int32(http.StatusCreated), Message: "success", Token: &control.TokenDescription{
		Token:      tokenDesc.Token,
		CreateTime: tokenDesc.CreateTime,
		Life:       tokenDesc.Life,
		IsAdmin:    tokenDesc.IsAdmin,
		Valid:      tokenDesc.IsValid,
	}}, nil
}

func (s *LiteKubeControl) QueryTokens(ctx context.Context, in *control.NoneValue) (*control.TokenValueList, error) {
	tokensDesc, err := readTokens(s.AuthFile)
	if err != nil {
		return &control.TokenValueList{StatusCode: int32(http.StatusInternalServerError), Message: "fail to query tokens"}, nil
	}

	tokens := make([]*control.TokenDescription, 0, len(tokensDesc))
	for _, tokenDesc := range tokensDesc {
		tokens = append(tokens, &control.TokenDescription{
			Token:      tokenDesc.Token,
			CreateTime: tokenDesc.CreateTime,
			Life:       tokenDesc.Life,
			IsAdmin:    tokenDesc.IsAdmin,
			Valid:      tokenDesc.IsValid,
		})
	}

	return &control.TokenValueList{StatusCode: int32(http.StatusCreated), Message: "success", TokenList: tokens}, nil
}

func (s *LiteKubeControl) DeleteToken(ctx context.Context, in *control.TokenString) (*control.NoneResponse, error) {
	klog.Info("delete litekube service token")
	md, exist := metadata.FromIncomingContext(ctx)
	if !exist {
		return &control.NoneResponse{StatusCode: int32(http.StatusCreated), Message: "need token"}, nil
	}

	if in.GetToken() == "" {
		return &control.NoneResponse{StatusCode: int32(http.StatusExpectationFailed), Message: "need value to mark token to delete"}, nil
	}

	tokenDesc, err := deleteToken(s.AuthFile, in.GetToken(), md["token"][0])
	if err != nil {
		return &control.NoneResponse{StatusCode: int32(http.StatusInternalServerError), Message: "fail to delete token"}, nil
	}

	if tokenDesc != nil {
		return &control.NoneResponse{StatusCode: int32(http.StatusCreated), Message: fmt.Sprintf("now %s is remove", tokenDesc.Token)}, nil
	} else {
		return &control.NoneResponse{StatusCode: int32(http.StatusCreated), Message: "token to delete is not exist"}, nil
	}

}

func (s *LiteKubeControl) BootStrapKubelet(ctx context.Context, request *control.BootStrapKubeletRequest) (*control.BootStrapKubeletResponse, error) {
	klog.Info("bootstrap for kubelet")
	var StatusCode int = http.StatusOK
	var returnErr error = nil
	buf := &bytes.Buffer{}
	var data interface{}

	// if !global.Exists(s.AuthTokenFile, s.ValidateApiserverCA) {
	// 	return nil, fmt.Errorf("loss info")
	// }
	if s == nil || !global.Exists(s.ValidateApiserverServerCA, s.ValidateApiserverClientCA, s.AuthTokenFile) {
		klog.Errorf("loss args")
		returnErr = fmt.Errorf("loss args to generate kubeconfig")
		StatusCode = http.StatusInternalServerError
		goto ERROR
	}

	klog.Infof("=>Get request for kubelet certificate bootstrap, latest will be recode to %s", s.BufferPath)

	data = struct {
		URL    string
		CACert string
		Token  string
	}{
		URL:    s.EndPoint,
		CACert: s.ValidateApiserverServerCABase64,
		Token:  s.Token,
	}

	kubelet_bootstrap_kubeconfig_template.Execute(buf, &data)

	if s.ClusterDNS == "" || s.ValidateApiserverClientCABase64 == "" {
		klog.Errorf("loss args cluster-dns or validata apiserver ca for kubelet")
		returnErr = fmt.Errorf("loss args to finish bootstrap for kubelet-litekube")
		StatusCode = http.StatusInternalServerError
		goto ERROR
	}

	return &control.BootStrapKubeletResponse{StatusCode: int32(StatusCode), Message: "success", Kubeconfig: base64.StdEncoding.EncodeToString(buf.Bytes()), ValidataCaCert: s.ValidateApiserverClientCABase64, ClusterDNS: s.ClusterDNS}, nil

ERROR:
	return &control.BootStrapKubeletResponse{StatusCode: int32(StatusCode), Message: returnErr.Error()}, nil
}

func (s *LiteKubeControl) BootStrapKubeProxy(ctx context.Context, request *control.BootStrapKubeProxyRequest) (*control.BootStrapKubeProxyResponse, error) {
	klog.Info("bootstrap for kube-proxy")
	var StatusCode int = http.StatusOK
	var returnErr error = nil
	buf := &bytes.Buffer{}
	var data interface{}
	var cert, key string

	// if !global.Exists(s.AuthTokenFile, s.ValidateApiserverCA) {
	// 	return nil, fmt.Errorf("loss info")
	// }
	if s == nil || !global.Exists(s.ValidateApiserverServerCA, s.ValidateApiserverClientCA, s.AuthTokenFile) {
		klog.Errorf("loss args")
		returnErr = fmt.Errorf("loss args to generate kubeconfig")
		StatusCode = http.StatusInternalServerError
		goto ERROR
	}

	klog.Infof("=>Get request for kube-proxy certificate bootstrap, latest will be recode to %s", s.BufferPath)

	if _, err := certificate.GenerateClientCertKey(true, "system:kube-proxy", nil, s.SignalApiserverClientCert, s.SignalApiserverClientKey, s.CacheCertPath, s.CacheKeyPath); err != nil {
		returnErr = fmt.Errorf("fail to write certificates cache")
		StatusCode = http.StatusInternalServerError
		goto ERROR
	}

	if certStr, err := certificate.LoadCertificateAsBase64(s.CacheCertPath); err != nil {
		returnErr = fmt.Errorf("fail to read certificates cache")
		StatusCode = http.StatusInternalServerError
		goto ERROR
	} else {
		cert = certStr
	}

	if keyStr, err := certificate.LoadFileAsBase64(s.CacheKeyPath); err != nil {
		returnErr = fmt.Errorf("fail to read certificates key cache")
		StatusCode = http.StatusInternalServerError
		goto ERROR
	} else {
		key = keyStr
	}

	data = struct {
		URL        string
		CACert     string
		ClientCert string
		ClientKey  string
	}{
		URL:        s.EndPoint,
		CACert:     s.ValidateApiserverServerCABase64,
		ClientCert: cert,
		ClientKey:  key,
	}

	kubeproxy_bootstrap_kubeconfig_template.Execute(buf, &data)

	return &control.BootStrapKubeProxyResponse{StatusCode: int32(StatusCode), Message: "success", Kubeconfig: base64.StdEncoding.EncodeToString(buf.Bytes()), ClusterCIDR: s.ClusterCIDR}, nil

ERROR:
	return &control.BootStrapKubeProxyResponse{StatusCode: int32(StatusCode), Message: returnErr.Error()}, nil
}

func readTokens(authFile string) (map[string]*TokenDesc, error) {
	if err := os.MkdirAll(globalFold, os.FileMode(0644)); err != nil {
		return nil, err
	}

	datas := make(map[string]*TokenDesc)

	if global.Exists(authFile) {
		bytes, err := ioutil.ReadFile(authFile)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(bytes, &datas); err != nil {
			return nil, err
		}

		for key, tokenDesc := range datas {
			if len(tokenDesc.Token) != 32 {
				delete(datas, key)
				continue
			}

			// permanent account
			if tokenDesc.Life < 0 {
				continue
			}

			createTime, err := time.Parse("2006-01-02 15:04:05", tokenDesc.CreateTime)
			if err != nil {
				return nil, err
			}

			if time.Now().UTC().Unix()-createTime.UTC().Unix() > tokenDesc.Life*60 {
				if !tokenDesc.IsAdmin {
					delete(datas, key)
				} else {
					tokenDesc.IsValid = false
				}
			}
		}

		if bytes, err := yaml.Marshal(datas); err != nil {
			return nil, err
		} else {
			if err := ioutil.WriteFile(authFile, bytes, fs.FileMode(0644)); err != nil {
				return nil, err
			}
		}

		return datas, nil
	} else {
		return nil, nil
	}
}

func createToken(leaveTime int64, createBy string, authFile string, isAdmin bool) (*TokenDesc, error) {
	tokens, err := readTokens(authFile)
	if err != nil {
		return nil, err
	}

	tokenStr, err := token.Random(32)
	if err != nil {
		return nil, err
	}

	if leaveTime == 0 {
		leaveTime = 10
	}

	t := &TokenDesc{
		Token:      tokenStr,
		CreateBy:   createBy,
		CreateTime: time.Now().UTC().Format("2006-01-02 15:04:05"),
		Life:       leaveTime,
		IsAdmin:    isAdmin,
		IsValid:    true,
	}

	tokens[tokenStr] = t

	bytes, err := yaml.Marshal(tokens)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(globalFold, os.FileMode(0644)); err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(authFile, bytes, fs.FileMode(0644)); err != nil {
		return nil, err
	} else {
		return t, nil
	}
}

func queryToken(token string, authFile string) *TokenDesc {
	tokenDescs, err := readTokens(authFile)
	if err != nil {
		return nil
	} else {
		if tokenDesc, ok := tokenDescs[token]; ok {
			return tokenDesc
		} else {
			return nil
		}
	}
}

func ensureToken(authFile string) error {
	if err := os.MkdirAll(globalFold, os.FileMode(0644)); err != nil {
		return err
	}

	bytes, err := ioutil.ReadFile(globalToken)

	if err != nil || len(string(bytes)) != 32 {
		tokenDesc, err := createToken(-1, "litekube", authFile, true)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(globalToken, []byte(tokenDesc.Token), fs.FileMode(0644))
	}

	return nil
}

func deleteToken(authFile string, token string, deleteBy string) (*TokenDesc, error) {
	bytes, err := ioutil.ReadFile(authFile)
	if err != nil {
		return nil, err
	}

	datas := make(map[string]TokenDesc)
	if err := yaml.Unmarshal(bytes, &datas); err != nil {
		return nil, err
	}

	if tokenDesc, ok := datas[token]; ok {
		delete(datas, token)

		tokenbytes, err := yaml.Marshal(datas)
		if err != nil {
			return nil, err
		}

		if err := os.MkdirAll(globalFold, os.FileMode(0644)); err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(authFile, tokenbytes, fs.FileMode(0644)); err != nil {
			return nil, err
		}

		return &tokenDesc, nil
	}

	return nil, nil
}
