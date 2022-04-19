package runtime

import (
	"context"
	goruntime "runtime"

	// link to github.com/Litekube/kine, we have make some addition

	"github.com/litekube/LiteKube/pkg/logger"
	"github.com/litekube/LiteKube/pkg/options/leader/netmanager"
	"k8s.io/klog/v2"
)

type NetWorkClient struct {
	ctx     context.Context
	LogPath string
}

func NewNetWorkClient(ctx context.Context, opt *netmanager.NetManagerOptions, logPath string) *NetWorkClient {

	return &NetWorkClient{
		ctx:     ctx,
		LogPath: logPath,
	}
}

// start run in routine and no wait
func (s *NetWorkClient) Run() error {
	ptr, _, _, ok := goruntime.Caller(0)
	if ok {
		logger.DefaultLogger.SetLog(goruntime.FuncForPC(ptr).Name(), s.LogPath)
	} else {
		klog.Errorf("fail to init kine log")
	}

	klog.Info("run network manager client")
	return nil
}
