package runtime

import (
	"context"
	goruntime "runtime"

	// link to github.com/Litekube/kine, we have make some addition
	"github.com/litekube/LiteKube/pkg/leader/authentication"
	"github.com/litekube/LiteKube/pkg/logger"
	"k8s.io/klog/v2"
)

type NetWorkManager struct {
	ctx     context.Context
	LogPath string
}

func NewNetWorkManager(ctx context.Context, opt *authentication.NetworkAuthentication, logPath string) *NetWorkManager {

	return &NetWorkManager{
		ctx:     ctx,
		LogPath: logPath,
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
