package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"github.com/litekube/LiteKube/pkg/certificate"
	"github.com/litekube/LiteKube/pkg/global"
	leaderAuth "github.com/litekube/LiteKube/pkg/leader/authentication"
	"github.com/litekube/LiteKube/pkg/logger"
	"github.com/litekube/LiteKube/pkg/options/leader/netmanager"
	options "github.com/litekube/LiteKube/pkg/options/worker"
	globaloptions "github.com/litekube/LiteKube/pkg/options/worker/global"
	"github.com/litekube/LiteKube/pkg/options/worker/kubelet"
	"github.com/litekube/LiteKube/pkg/options/worker/kubeproxy"
	"github.com/litekube/LiteKube/pkg/worker/authentication"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

func (workerRuntime *WorkerRuntime) SetFlags(opt *options.WorkerOptions) {
	workerRuntime.FlagsOption = opt
}

// load all flags
func (workerRuntime *WorkerRuntime) LoadFlags() error {
	if workerRuntime.FlagsOption == nil {
		return fmt.Errorf("no flags input")
	}

	// init global flags
	if err := workerRuntime.LoadGloabl(); err != nil {
		return err
	}

	// init network manager flags
	if err := workerRuntime.LoadNetManager(); err != nil {
		return err
	}

	// run kine server, network manager server, network client to make environment for litekube
	if err := workerRuntime.RunForward(); err != nil {
		return err
	}

	// init flags for kube-apiserver
	if err := workerRuntime.LoadKubelet(); err != nil {
		return err
	}

	// init flags for controller-manager
	if err := workerRuntime.LoadKubeProxy(); err != nil {
		return err
	}

	if config, err := yaml.Marshal(workerRuntime.RuntimeOption); err != nil {
		return err
	} else {
		startupDir := filepath.Join(workerRuntime.RuntimeOption.GlobalOptions.WorkDir, "startup/")
		if err := os.MkdirAll(startupDir, os.ModePerm); err != nil {
			return err
		}
		if err := ioutil.WriteFile(filepath.Join(startupDir, "worker.yaml"), config, os.ModePerm); err != nil {
			return err
		}
	}

	// release unuseful data to save resource
	workerRuntime.FlagsOption = nil
	return nil
}

// load or generate args for litekube-global
func (workerRuntime *WorkerRuntime) LoadGloabl() error {
	if workerRuntime.FlagsOption.NetmamagerOptions.Token != "local" {
		defer func() {
			// set log
			// if workerRuntime.RuntimeOption.GlobalOptions.LogToDir {
			// 	klog.MaxSize = 10240
			// 	logfile := filepath.Join(workerRuntime.RuntimeOption.GlobalOptions.LogDir, "litekube.log")
			// 	flag.Set("log_file", logfile)
			// 	flag.Set("logtostderr", "false")

			// 	if workerRuntime.RuntimeOption.GlobalOptions.LogToStd {
			// 		flag.Set("alsologtostderr", "true")
			// 	} else {
			// 		flag.Set("alsologtostderr", "false")
			// 	}

			// } else {
			// 	flag.Set("logtostderr", fmt.Sprintf("%t", workerRuntime.RuntimeOption.GlobalOptions.LogToStd))
			// }
			flag.Set("log_file", "false")
			flag.Set("logtostderr", "false")
			klog.SetOutput(logger.NewLogWriter(workerRuntime.RuntimeOption.GlobalOptions.LogToStd, workerRuntime.RuntimeOption.GlobalOptions.LogToDir, filepath.Join(workerRuntime.RuntimeOption.GlobalOptions.LogDir, "litekube.log")).Logger())
		}()
	}

	defer func() {
		workerRuntime.RuntimeAuthentication = NewRuntimeAuthentication(filepath.Join(workerRuntime.RuntimeOption.GlobalOptions.WorkDir, "tls/"))
	}()

	raw := workerRuntime.FlagsOption.GlobalOptions
	new := workerRuntime.RuntimeOption.GlobalOptions

	// set default work-dir="~/litekube/"
	new.WorkDir = raw.WorkDir
	if new.WorkDir == "" {
		new.WorkDir = globaloptions.DefaultGO.WorkDir
	}

	// log
	new.LogDir = raw.LogDir
	if new.LogDir == "" {
		new.LogDir = filepath.Join(new.WorkDir, "logs/worker/")
	}

	new.LogToDir = raw.LogToDir
	if new.LogToDir {
		if err := os.MkdirAll(new.LogDir, os.FileMode(0666)); err != nil {
			return err
		}
	}
	new.LogToStd = raw.LogToStd

	if raw.LeaderToken != "" {
		new.LeaderToken = raw.LeaderToken
	} else {
		new.LeaderToken = globaloptions.DefaultGO.LeaderToken
	}
	return nil
}

func (workerRuntime *WorkerRuntime) LoadNetManager() error {
	raw := workerRuntime.FlagsOption.NetmamagerOptions
	new := workerRuntime.RuntimeOption.NetmamagerOptions

	if raw.Token != "" {
		new.Token = raw.Token
	} else {
		new.Token = netmanager.DefaultNMO.Token
	}

	// check bind-address
	if ip := net.ParseIP(raw.RegisterOptions.Address); ip == nil {
		new.RegisterOptions.Address = netmanager.DefaultRONO.Address
	} else {
		new.RegisterOptions.Address = raw.RegisterOptions.Address
	}

	if ip := net.ParseIP(raw.JoinOptions.Address); ip == nil {
		new.JoinOptions.Address = netmanager.DefaultJONO.Address
	} else {
		new.JoinOptions.Address = raw.JoinOptions.Address
	}

	// check https port
	if raw.RegisterOptions.SecurePort < 1 || raw.RegisterOptions.SecurePort > 65535 {
		new.RegisterOptions.SecurePort = netmanager.DefaultRONO.SecurePort
	} else {
		new.RegisterOptions.SecurePort = raw.RegisterOptions.SecurePort
	}

	// check https port
	if raw.JoinOptions.SecurePort < 1 || raw.JoinOptions.SecurePort > 65535 {
		new.JoinOptions.SecurePort = netmanager.DefaultJONO.SecurePort
	} else {
		new.JoinOptions.SecurePort = raw.JoinOptions.SecurePort
	}

	// try to load certificate provide by user
	if certificate.NotExists(raw.RegisterOptions.CACert, raw.RegisterOptions.ClientCertFile, raw.RegisterOptions.ClientkeyFile) && certificate.NotExists(raw.JoinOptions.CACert, raw.JoinOptions.ClientCertFile, raw.JoinOptions.ClientkeyFile) && raw.NodeToken == "" {
		// check client certificate
		// into TLS bootstrap
		klog.Info("start load network manager client certificate and node-token by --token")
		workerRuntime.RuntimeAuthentication.NetWorkManagerClient = leaderAuth.NewControllerClientAuthentication(workerRuntime.RuntimeAuthentication.CertDir, new.Token, &new.RegisterOptions.Address, &new.RegisterOptions.SecurePort, &new.JoinOptions.Address, &new.JoinOptions.SecurePort)
		if err := workerRuntime.RuntimeAuthentication.NetWorkManagerClient.GenerateOrSkip(); err != nil {
			return err
		}

		if !workerRuntime.RuntimeAuthentication.NetWorkManagerClient.Check() {
			return fmt.Errorf("fail to load network-manager TLS args")
		}

		// node token
		if nodeToken, err := workerRuntime.RuntimeAuthentication.NetWorkManagerClient.Nodetoken(); err != nil {
			return err
		} else {
			new.NodeToken = nodeToken
		}

		// cert
		// join
		new.JoinOptions.CACert = workerRuntime.RuntimeAuthentication.NetWorkManagerClient.JoinCACert
		new.JoinOptions.ClientCertFile = workerRuntime.RuntimeAuthentication.NetWorkManagerClient.JoinClientCert
		new.JoinOptions.ClientkeyFile = workerRuntime.RuntimeAuthentication.NetWorkManagerClient.JoinClientkey

		// register
		new.RegisterOptions.CACert = workerRuntime.RuntimeAuthentication.NetWorkManagerClient.RegisterCACert
		new.RegisterOptions.ClientCertFile = workerRuntime.RuntimeAuthentication.NetWorkManagerClient.RegisterClientCert
		new.RegisterOptions.ClientkeyFile = workerRuntime.RuntimeAuthentication.NetWorkManagerClient.RegisterClientkey

		klog.Info("success to load network manager client certificates node-token by --token")
		return nil

	} else {
		if certificate.ValidateTLSPair(raw.RegisterOptions.ClientCertFile, raw.RegisterOptions.ClientkeyFile) && certificate.ValidateCA(raw.RegisterOptions.ClientCertFile, raw.RegisterOptions.CACert) && certificate.ValidateTLSPair(raw.JoinOptions.ClientCertFile, raw.JoinOptions.ClientkeyFile) && certificate.ValidateCA(raw.JoinOptions.ClientCertFile, raw.JoinOptions.CACert) && len(raw.NodeToken) > 0 {
			// cert
			// join
			new.JoinOptions.CACert = raw.JoinOptions.CACert
			new.JoinOptions.ClientCertFile = raw.JoinOptions.ClientCertFile
			new.JoinOptions.ClientkeyFile = raw.JoinOptions.ClientkeyFile

			// register
			new.RegisterOptions.CACert = raw.RegisterOptions.CACert
			new.RegisterOptions.ClientCertFile = raw.RegisterOptions.ClientCertFile
			new.RegisterOptions.ClientkeyFile = raw.RegisterOptions.ClientkeyFile
			new.NodeToken = raw.NodeToken
			workerRuntime.RuntimeAuthentication.NetWorkManagerClient = nil
			klog.Infof("network manager client certificates specified ok, ignore --token")
		} else {
			raw.PrintFlags("error-tip", func() func(format string, a ...interface{}) error {
				return func(format string, a ...interface{}) error {
					klog.Errorf(format, a...)
					return nil
				}
			}())
			return fmt.Errorf("you have provide bad network manager client certificates or node-token for network manager")
		}
	}

	return nil
}

func (workerRuntime *WorkerRuntime) LoadKubelet() error {
	raw := workerRuntime.FlagsOption.KubeletOptions
	new := workerRuntime.RuntimeOption.KubeletOptions

	new.ReservedOptions = raw.ReservedOptions
	new.IgnoreOptions = raw.IgnoreOptions

	// load *KubeletLitekubeOptions
	// pod-infra-container-image
	if raw.Options.PodInfraContainerImage != "" {
		new.Options.PodInfraContainerImage = raw.Options.PodInfraContainerImage
	} else {
		new.Options.PodInfraContainerImage = kubelet.DefaultKLO.PodInfraContainerImage
	}
	// cert-dir
	if raw.Options.CertDir != "" {
		new.Options.CertDir = raw.Options.CertDir
	} else {
		new.Options.CertDir = filepath.Join(workerRuntime.RuntimeAuthentication.CertDir, workerRuntime.RuntimeOption.GlobalOptions.LeaderToken, "kubelet")
	}

	// load * ProfessionalOptions
	// node-ip
	// if ip := net.ParseIP(raw.ProfessionalOptions.NodeIp); ip == nil {
	// 	if localIp, err := workerRuntime.NetworkRegisterClient.QueryIp(); err != nil {
	// 		return err
	// 	} else {
	// 		new.ProfessionalOptions.NodeIp = localIp
	// 	}
	// } else {
	// 	new.ProfessionalOptions.NodeIp = raw.ProfessionalOptions.NodeIp
	// }
	// ips := []string{}
	// nodeIp := ""
	// addNodeIp := true
	// addLocalIp := true
	// if localIp, err := workerRuntime.NetworkRegisterClient.QueryIp(); err != nil {
	// 	return err
	// } else {
	// 	nodeIp = localIp
	// }
	// for _, ipStr := range strings.Split(raw.ProfessionalOptions.NodeIp, ",") {
	// 	ipStr = strings.TrimSpace(ipStr)
	// 	if ip := net.ParseIP(ipStr); ip == nil {
	// 		if ipStr == nodeIp {
	// 			addNodeIp = false
	// 		}

	// 		if ipStr == global.LocalhostIP.String() {
	// 			addLocalIp = false
	// 		}

	// 		ips = append(ips, ipStr)
	// 	}
	// }

	// if addLocalIp {
	// 	ips = append(ips, global.LocalhostIP.String())
	// }

	// if addNodeIp {
	// 	ips = append(ips, nodeIp)
	// }

	// new.ProfessionalOptions.NodeIp = strings.Join(ips, ",")

	// kubeconfig
	if raw.ProfessionalOptions.Kubeconfig != "" {
		new.ProfessionalOptions.Kubeconfig = raw.ProfessionalOptions.Kubeconfig
	} else {
		new.ProfessionalOptions.Kubeconfig = filepath.Join(workerRuntime.RuntimeAuthentication.CertDir, workerRuntime.RuntimeOption.GlobalOptions.LeaderToken, "kubelet.kubeconfig")
	}

	// runtime-cgroup
	if raw.ProfessionalOptions.RuntimeCgroups != "" {
		new.ProfessionalOptions.RuntimeCgroups = raw.ProfessionalOptions.RuntimeCgroups
	} else {
		new.ProfessionalOptions.RuntimeCgroups = kubelet.DefaultKPO.RuntimeCgroups
	}

	// cgroup-driver
	if raw.ProfessionalOptions.CgroupDriver != "" {
		new.ProfessionalOptions.CgroupDriver = raw.ProfessionalOptions.CgroupDriver
	} else {
		new.ProfessionalOptions.CgroupDriver = kubelet.DefaultKPO.CgroupDriver
	}

	// hostname-override
	if raw.ProfessionalOptions.HostnameOverride != "" {
		new.ProfessionalOptions.HostnameOverride = raw.ProfessionalOptions.HostnameOverride
	} else {
		if localIp, err := workerRuntime.NetworkRegisterClient.QueryIp(); err != nil {
			return err
		} else {
			new.ProfessionalOptions.HostnameOverride = localIp
		}
	}

	// container-runtime
	if raw.ProfessionalOptions.ContainerRuntime != "" {
		new.ProfessionalOptions.ContainerRuntime = raw.ProfessionalOptions.ContainerRuntime
	} else {
		new.ProfessionalOptions.ContainerRuntime = kubelet.DefaultKPO.ContainerRuntime
	}

	// container-runtime-endpoint
	if raw.ProfessionalOptions.ContainerRuntimeEndpoint != "" {
		new.ProfessionalOptions.ContainerRuntimeEndpoint = raw.ProfessionalOptions.ContainerRuntimeEndpoint
	} else {
		new.ProfessionalOptions.ContainerRuntimeEndpoint = kubelet.DefaultKPO.ContainerRuntimeEndpoint
	}

	if raw.ProfessionalOptions.BootstrapKubeconfig != "" && raw.ProfessionalOptions.Config != "" && global.Exists(raw.ProfessionalOptions.BootstrapKubeconfig, raw.ProfessionalOptions.Config) {
		new.ProfessionalOptions.BootstrapKubeconfig = raw.ProfessionalOptions.BootstrapKubeconfig
		new.ProfessionalOptions.Config = raw.ProfessionalOptions.Config
	} else {
		// into bootstrap
		if workerRuntime.RuntimeOption.GlobalOptions.LeaderToken == "" {
			return fmt.Errorf("leader token is need for join into cluster")
		}

		if workerRuntime.RuntimeAuthentication.KubernetesNode == nil {
			workerRuntime.RuntimeAuthentication.KubernetesNode = authentication.NewKubernetesNode(workerRuntime.RuntimeOption.GlobalOptions.WorkDir,
				workerRuntime.RuntimeOption.GlobalOptions.LeaderToken,
				workerRuntime.NetworkRegisterClient,
			)
			if workerRuntime.RuntimeAuthentication.KubernetesNode == nil {
				return fmt.Errorf("fail to run TLS-bootstrap for worker")
			}

			if err := workerRuntime.RuntimeAuthentication.KubernetesNode.GenerateOrSkip(); err != nil {
				return err
			}
		}

		new.ProfessionalOptions.BootstrapKubeconfig = workerRuntime.RuntimeAuthentication.KubernetesNode.BootStrapKubeConfig
		new.ProfessionalOptions.Config = workerRuntime.RuntimeAuthentication.KubernetesNode.KubeletConfig
	}

	return nil
}

func (workerRuntime *WorkerRuntime) LoadKubeProxy() error {
	raw := workerRuntime.FlagsOption.KubeProxyOptions
	new := workerRuntime.RuntimeOption.KubeProxyOptions

	new.ReservedOptions = raw.ReservedOptions
	new.IgnoreOptions = raw.IgnoreOptions

	// load *KubeletLitekubeOptions

	// load * ProfessionalOptions
	// hostname-override
	if raw.ProfessionalOptions.HostnameOverride != "" {
		new.ProfessionalOptions.HostnameOverride = raw.ProfessionalOptions.HostnameOverride
	} else {
		if localIp, err := workerRuntime.NetworkRegisterClient.QueryIp(); err != nil {
			return err
		} else {
			new.ProfessionalOptions.HostnameOverride = localIp
		}
	}

	// container-runtime
	if raw.ProfessionalOptions.ProxyMode != "" {
		new.ProfessionalOptions.ProxyMode = raw.ProfessionalOptions.ProxyMode
	} else {
		new.ProfessionalOptions.ProxyMode = kubeproxy.DefaultKPPO.ClusterCidr
	}

	if raw.ProfessionalOptions.Kubeconfig != "" && raw.ProfessionalOptions.ClusterCidr != "" && global.Exists(raw.ProfessionalOptions.Kubeconfig) {
		new.ProfessionalOptions.Kubeconfig = raw.ProfessionalOptions.Kubeconfig
		new.ProfessionalOptions.ClusterCidr = raw.ProfessionalOptions.ClusterCidr
	} else {
		// into bootstrap
		if workerRuntime.RuntimeOption.GlobalOptions.LeaderToken == "" {
			return fmt.Errorf("leader token is need for join into cluster")
		}

		if workerRuntime.RuntimeAuthentication.KubernetesNode == nil {
			workerRuntime.RuntimeAuthentication.KubernetesNode = authentication.NewKubernetesNode(workerRuntime.RuntimeOption.GlobalOptions.WorkDir,
				workerRuntime.RuntimeOption.GlobalOptions.LeaderToken,
				workerRuntime.NetworkRegisterClient,
			)
			if workerRuntime.RuntimeAuthentication.KubernetesNode == nil {
				return fmt.Errorf("fail to run TLS-bootstrap for worker")
			}

			if err := workerRuntime.RuntimeAuthentication.KubernetesNode.GenerateOrSkip(); err != nil {
				return err
			}
		}

		new.ProfessionalOptions.Kubeconfig = workerRuntime.RuntimeAuthentication.KubernetesNode.KubeProxyKubeConfig

		// cluster-cidr
		if raw.ProfessionalOptions.ClusterCidr != "" {
			if _, _, err := net.ParseCIDR(raw.ProfessionalOptions.ClusterCidr); err == nil {
				new.ProfessionalOptions.ClusterCidr = raw.ProfessionalOptions.ClusterCidr
			}
		}

		if new.ProfessionalOptions.ClusterCidr == "" {
			if clusterCIDR, err := workerRuntime.RuntimeAuthentication.KubernetesNode.ReadAddition("cluster-cidr"); err != nil {
				return err
			} else {
				if _, _, err := net.ParseCIDR(clusterCIDR); err != nil {
					return fmt.Errorf("bad cluster-cidr get from leader")
				} else {
					new.ProfessionalOptions.ClusterCidr = clusterCIDR
				}
			}
		}
	}

	return nil
}
