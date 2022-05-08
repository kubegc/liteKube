package runtime

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	// link to github.com/Litekube/kine, we have make some addition

	"github.com/litekube/LiteKube/pkg/options/leader/apiserver"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kube-apiserver/app"
)

type Apiserver struct {
	ctx     context.Context
	LogPath string
	Options *apiserver.ApiserverOptions
	// Handler       *http.Handler
	// Authenticator *authenticator.Request
}

func NewApiserver(ctx context.Context, opt *apiserver.ApiserverOptions, logPath string) *Apiserver {
	return &Apiserver{
		ctx:     ctx,
		Options: opt,
		LogPath: logPath,
		// Handler:       nil,
		// Authenticator: nil,
	}
}

// start run in routine and no wait
func (s *Apiserver) Run() error {
	klog.Info("run kube-apiserver")

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

	command := app.NewAPIServerCommand(s.ctx.Done())
	command.SetArgs(argsValue)

	go func() {
		for i := 0; i < 10; i++ {
			etcdServer := strings.Split(apiserver.DefaultAPO.EtcdServers, ",")
			if len(etcdServer) < 1 || etcdServer[0] == "" {
				klog.Errorf("bad etcd servers.")
			}
			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}

			if resp, err := client.Get(fmt.Sprintf("%s/health", etcdServer[0])); err != nil {
				klog.Infof("waiting ETCD ready")
				time.Sleep(1 * time.Second)
				if i == 9 {
					klog.Errorf("error to start check ETCD health: error: %s", err.Error())
					return
				}
				continue
			} else {
				if 200 <= resp.StatusCode && resp.StatusCode < 300 {
					klog.Infof("check ETCD ok.")
					break
				} else {
					klog.Error("ETCD meet some error, error code: %d", resp.StatusCode)
					return
				}

			}
		}

		klog.Infof("==>kube-apiserver: %s\n", argsValue)

		err := command.ExecuteContext(s.ctx)
		if err != nil {
			klog.Fatalf("kube-apiserver exited: %v", err)
		}
	}()

	// startupConfig := <-app.StartupConfig

	// s.Handler = &startupConfig.Handler
	// s.Authenticator = &startupConfig.Authenticator

	return nil
}

// func (s *Apiserver) StartUpConfig() (*http.Handler, *authenticator.Request) {
// 	return s.Handler, s.Authenticator
// }
