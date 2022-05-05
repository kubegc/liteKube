package runtime

import (
	"context"
	"fmt"
	"net/http"
	"time"

	// link to github.com/Litekube/kine, we have make some addition

	"github.com/litekube/LiteKube/pkg/options/worker/kubeproxy"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kube-proxy/app"
)

type KubeProxy struct {
	ctx     context.Context
	LogPath string
	Options *kubeproxy.KubeProxyOptions
}

func NewKubeProxy(ctx context.Context, opt *kubeproxy.KubeProxyOptions, logPath string) *KubeProxy {
	return &KubeProxy{
		ctx:     ctx,
		Options: opt,
		LogPath: logPath,
	}
}

// start run in routine and no wait
func (s *KubeProxy) Run() error {
	args, err := s.Options.ToMap()
	if err != nil {
		return err
	}

	argsValue := make([]string, 0, len(args))
	for k, v := range args {
		if v == "-" || v == "" {
			continue
		}
		argsValue = append(argsValue, fmt.Sprintf("--%s=%s", k, v))
	}

	command := app.NewProxyCommand()
	command.SetArgs(argsValue)

	klog.Infof("==>kube-proxy: %s\n", argsValue)

	go func() {
		err := command.ExecuteContext(s.ctx)
		if err != nil {
			klog.Fatalf("kube-proxy exited: %v", err)
		}
	}()

	return nil
}

func (s *KubeProxy) WaitForAPIServer(ctx context.Context, kubeAdminPath string) error {
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeAdminPath)
	if err != nil {
		return err
	}

	k8sClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	for {

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-promise(func() error { return waitForAPIServerReady(ctx, k8sClient, 30*time.Second) }):
			if err != nil {
				klog.Infof("Waiting for API server to become available")
				continue
			}
			return err
		}
	}
}

func waitForAPIServerReady(ctx context.Context, client kubernetes.Interface, timeout time.Duration) error {
	var lastErr error
	restClient := client.Discovery().RESTClient()

	err := wait.PollImmediateWithContext(ctx, time.Second, timeout, func(ctx context.Context) (bool, error) {
		healthStatus := 0
		result := restClient.Get().AbsPath("/readyz").Do(ctx).StatusCode(&healthStatus)
		if rerr := result.Error(); rerr != nil {
			lastErr = errors.Wrap(rerr, "failed to get apiserver /readyz status")
			return false, nil
		}
		if healthStatus != http.StatusOK {
			content, _ := result.Raw()
			lastErr = fmt.Errorf("APIServer isn't ready: %v", string(content))
			logrus.Warnf("APIServer isn't ready yet: %v. Waiting a little while.", string(content))
			return false, nil
		}

		return true, nil
	})

	if err != nil {
		return fmt.Errorf("Error: %s and %s", err.Error(), lastErr.Error())
	}

	return nil
}

func promise(f func() error) <-chan error {
	c := make(chan error, 1)
	go func() {
		c <- f()
		close(c)
	}()
	return c
}
