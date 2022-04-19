package runtime

import (
	"context"
	goruntime "runtime"

	// link to github.com/Litekube/kine, we have make some addition
	"github.com/litekube/LiteKube/pkg/leader/authentication"
	"github.com/litekube/LiteKube/pkg/logger"
	"github.com/litekube/LiteKube/pkg/options/leader/netmanager"
	"k8s.io/klog/v2"
)

type NetWorkManager struct {
	ctx     context.Context
	LogPath string

	// register
	RegisterBindAddress string
	RegisterPort        uint16
	RegisterCACert      string
	RegisterCAKey       string
	RegisterServerCert  string
	RegisterServerkey   string

	// join
	JoinBindAddress string
	JoinPort        uint16
	JoinCACert      string
	JoinCAKey       string
	JoinServerCert  string
	JoinServerkey   string
}

func NewNetWorkManager(ctx context.Context, opt *authentication.NetworkAuthentication, clientOpt *netmanager.NetManagerOptions, logPath string) *NetWorkManager {

	return &NetWorkManager{
		ctx:     ctx,
		LogPath: logPath,

		// register
		RegisterBindAddress: opt.RegisterBindAddress,
		RegisterPort:        clientOpt.RegisterOptions.SecurePort,
		RegisterCACert:      opt.RegisterCACert,
		RegisterCAKey:       opt.RegisterCAKey,
		RegisterServerCert:  opt.RegisterServerCert,
		RegisterServerkey:   opt.RegisterServerkey,

		// join
		JoinBindAddress: opt.JoinBindAddress,
		JoinPort:        clientOpt.JoinOptions.SecurePort,
		JoinCACert:      opt.JoinCACert,
		JoinCAKey:       opt.JoinCAKey,
		JoinServerCert:  opt.JoinServerCert,
		JoinServerkey:   opt.JoinServerkey,
	}
}

// start run in routine and no wait
func (s *NetWorkManager) Run() error {
	ptr, _, _, ok := goruntime.Caller(0)
	if ok {
		logger.DefaultLogger.SetLog(goruntime.FuncForPC(ptr).Name(), s.LogPath)
	} else {
		klog.Errorf("fail to init kine log")
	}

	klog.Info("run network manager")

	return nil
}
