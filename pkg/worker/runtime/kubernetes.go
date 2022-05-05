package runtime

import (
	"context"
	"path/filepath"

	// link to github.com/Litekube/kine, we have make some addition

	"github.com/litekube/LiteKube/pkg/options/worker/kubelet"
	"github.com/litekube/LiteKube/pkg/options/worker/kubeproxy"
	"k8s.io/klog/v2"
)

type KubernatesClient struct {
	ctx              context.Context
	logPath          string
	kubelet          *Kubelet
	kubeProxy        *KubeProxy
	KubeletOptions   *kubelet.KubeletOptions
	KubeProxyOptions *kubeproxy.KubeProxyOptions
}

func NewKubernatesClient(ctx context.Context, kubeletOptions *kubelet.KubeletOptions, kubeProxyOptions *kubeproxy.KubeProxyOptions, logPath string) *KubernatesClient {
	return &KubernatesClient{
		ctx:              ctx,
		logPath:          logPath,
		kubelet:          NewKubelet(ctx, kubeletOptions, filepath.Join(logPath, "kubelet.log")),
		kubeProxy:        NewKubeProxy(ctx, kubeProxyOptions, filepath.Join(logPath, "kube-proxy.log")),
		KubeletOptions:   kubeletOptions,
		KubeProxyOptions: kubeProxyOptions,
	}
}

// start run in routine and no wait
func (s *KubernatesClient) Run() error {
	klog.Info("start to run kubernates node")

	if err := s.kubelet.Run(); err != nil {
		return err
	}

	if err := s.kubeProxy.Run(); err != nil {
		return err
	}

	return nil
}
