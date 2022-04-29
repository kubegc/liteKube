package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/litekube/LiteKube/pkg/certificate"
	"github.com/litekube/LiteKube/pkg/leader/authentication"
	"github.com/litekube/LiteKube/pkg/logger"
	options "github.com/litekube/LiteKube/pkg/options/leader"
	"github.com/litekube/LiteKube/pkg/options/leader/apiserver"
	"github.com/litekube/LiteKube/pkg/options/leader/controllermanager"
	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
	kineoptions "github.com/litekube/LiteKube/pkg/options/leader/kine"
	"github.com/litekube/LiteKube/pkg/options/leader/netmanager"
	"github.com/litekube/LiteKube/pkg/options/leader/scheduler"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

func (leaderRuntime *LeaderRuntime) SetFlags(opt *options.LeaderOptions) {
	leaderRuntime.FlagsOption = opt
}

// load all flags
func (leaderRuntime *LeaderRuntime) LoadFlags() error {
	if leaderRuntime.FlagsOption == nil {
		return fmt.Errorf("no flags input")
	}

	// init global flags
	if err := leaderRuntime.LoadGloabl(); err != nil {
		return err
	}

	// init kine flags
	if err := leaderRuntime.LoadKine(); err != nil {
		return err
	}

	// init network manager flags
	if err := leaderRuntime.LoadNetManager(); err != nil {
		return err
	}

	// run kine server, network manager server, network client to make environment for litekube
	if err := leaderRuntime.RunForward(); err != nil {
		return err
	}

	// init flags for kube-apiserver
	if err := leaderRuntime.LoadApiserver(); err != nil {
		return err
	}

	// init flags for controller-manager
	if err := leaderRuntime.LoadControllermanager(); err != nil {
		return err
	}

	// init kube-scheduler
	if err := leaderRuntime.LoadScheduler(); err != nil {
		return err
	}

	if config, err := yaml.Marshal(leaderRuntime.RuntimeOption); err != nil {
		return err
	} else {
		startupDir := filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "startup/")
		if err := os.MkdirAll(startupDir, os.ModePerm); err != nil {
			return err
		}
		if err := ioutil.WriteFile(filepath.Join(startupDir, "leader.yaml"), config, os.ModePerm); err != nil {
			return err
		}
	}

	// release unuseful data to save resource
	leaderRuntime.FlagsOption = nil
	return nil
}

// load or generate args for litekube-global
func (leaderRuntime *LeaderRuntime) LoadGloabl() error {
	defer func() {
		// set log
		// if leaderRuntime.RuntimeOption.GlobalOptions.LogToDir {
		// 	klog.MaxSize = 10240
		// 	logfile := filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.LogDir, "litekube.log")
		// 	flag.Set("log_file", logfile)
		// 	flag.Set("logtostderr", "false")

		// 	if leaderRuntime.RuntimeOption.GlobalOptions.LogToStd {
		// 		flag.Set("alsologtostderr", "true")
		// 	} else {
		// 		flag.Set("alsologtostderr", "false")
		// 	}

		// } else {
		// 	flag.Set("logtostderr", fmt.Sprintf("%t", leaderRuntime.RuntimeOption.GlobalOptions.LogToStd))
		// }
		flag.Set("log_file", "false")
		flag.Set("logtostderr", "false")
		klog.SetOutput(logger.NewDefaultLogger(leaderRuntime.RuntimeOption.GlobalOptions.LogToStd, leaderRuntime.RuntimeOption.GlobalOptions.LogToDir, filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.LogDir, "litekube.log")))
	}()

	defer func() {
		leaderRuntime.RuntimeAuthentication = NewRuntimeAuthentication(filepath.Join(leaderRuntime.RuntimeOption.GlobalOptions.WorkDir, "tls/"))
	}()

	raw := leaderRuntime.FlagsOption.GlobalOptions
	new := leaderRuntime.RuntimeOption.GlobalOptions

	// set default work-dir="~/litekube/"
	new.WorkDir = raw.WorkDir
	if new.WorkDir == "" {
		new.WorkDir = globaloptions.DefaultGO.WorkDir
	}

	// log
	new.LogDir = raw.LogDir
	if new.LogDir == "" {
		new.LogDir = filepath.Join(new.WorkDir, "logs/")
	}

	new.LogToDir = raw.LogToDir
	if new.LogToDir {
		if err := os.MkdirAll(new.LogDir, os.FileMode(0666)); err != nil {
			return err
		}
	}
	new.LogToStd = raw.LogToStd

	// kine
	new.RunKine = raw.RunKine
	// invalid etcd server will enable kine
	if leaderRuntime.RuntimeOption.ApiserverOptions.ProfessionalOptions.ECTDOptions.EtcdServers == "" {
		new.RunKine = true
	}

	// network-manager
	new.RunNetManager = raw.RunNetManager

	new.EnableWorker = raw.EnableWorker

	if !new.EnableWorker {
		new.WorkerConfig = ""
	}
	return nil
}

// load or generate args for kine server
// client certificate will be generate to path, too
func (leaderRuntime *LeaderRuntime) LoadKine() error {
	if !leaderRuntime.RuntimeOption.GlobalOptions.RunKine {
		leaderRuntime.RuntimeOption.KineOptions = nil
		return nil
	}

	raw := leaderRuntime.FlagsOption.KineOptions
	new := leaderRuntime.RuntimeOption.KineOptions

	// check bind-address
	if ip := net.ParseIP(raw.BindAddress); ip == nil {
		new.BindAddress = kineoptions.DefaultKO.BindAddress
	} else {
		new.BindAddress = raw.BindAddress
	}

	// check https port
	if raw.SecurePort < 1 || raw.SecurePort > 65535 {
		new.SecurePort = kineoptions.DefaultKO.SecurePort
	} else {
		new.SecurePort = raw.SecurePort
	}

	// check TLS certificate
	if certificate.NotExists(raw.CACert, raw.ServerCertFile, raw.ServerkeyFile) {
		klog.Info("built-in certificates for kine will be used")

		// invalid cert， generate kine certs
		leaderRuntime.OwnKineCert = true
		leaderRuntime.RuntimeAuthentication.Kine = authentication.NewKineAuthentication(leaderRuntime.RuntimeAuthentication.CertDir, new.BindAddress)
		if err := leaderRuntime.RuntimeAuthentication.Kine.GenerateOrSkip(); err != nil {
			return err
		}

		new.CACert = leaderRuntime.RuntimeAuthentication.Kine.CACert
		new.ServerCertFile = leaderRuntime.RuntimeAuthentication.Kine.ServerCert
		new.ServerkeyFile = leaderRuntime.RuntimeAuthentication.Kine.Serverkey
	} else {
		if !certificate.ValidateTLSPair(raw.ServerCertFile, raw.ServerkeyFile) || !certificate.ValidateCA(raw.ServerCertFile, raw.CACert) {
			klog.Errorf("You specified an unavailable certificate for kine")
			return fmt.Errorf("you specified an unavailable certificate for kine")
		}

		new.CACert = raw.CACert
		new.ServerCertFile = raw.ServerCertFile
		new.ServerkeyFile = raw.ServerkeyFile
		klog.Infof("kine certificate specified ok, skip generate")
	}

	return nil
}

// load network-manager client config
// if run-network-manager==true, runtime.RuntimeAuthentication.NetWorkManager will init
func (leaderRuntime *LeaderRuntime) LoadNetManager() error {
	raw := leaderRuntime.FlagsOption.NetmamagerOptions
	new := leaderRuntime.RuntimeOption.NetmamagerOptions

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

	// check Token
	if leaderRuntime.RuntimeOption.GlobalOptions.RunNetManager {
		// generate certificate for network manager
		klog.Info("certificates for built-in network manager server will be used")
		new.Token = "local"
		leaderRuntime.RuntimeAuthentication.NetWorkManager = authentication.NewNetworkAuthentication(leaderRuntime.RuntimeAuthentication.CertDir, new.RegisterOptions.Address, new.JoinOptions.Address)
		if err := leaderRuntime.RuntimeAuthentication.NetWorkManager.GenerateOrSkip(); err != nil {
			return err
		}
	} else {
		if raw.Token == "local" {
			return fmt.Errorf("bad token(local) to connect with network-manager, only enable when network manager run in leader node")
		}

		new.Token = raw.Token
	}

	// try to load certificate provide by user
	if certificate.NotExists(raw.RegisterOptions.CACert, raw.RegisterOptions.ClientCertFile, raw.RegisterOptions.ClientkeyFile) && certificate.NotExists(raw.JoinOptions.CACert, raw.JoinOptions.ClientCertFile, raw.JoinOptions.ClientkeyFile) && raw.NodeToken == "" {
		// check client certificate
		// into TLS bootstrap
		klog.Info("start load network manager client certificate and node-token by --token")
		leaderRuntime.RuntimeAuthentication.NetWorkManagerClient = authentication.NewNetworkManagerClient(leaderRuntime.RuntimeAuthentication.CertDir, new.Token, &new.RegisterOptions.Address, &new.RegisterOptions.SecurePort, &new.JoinOptions.Address, &new.JoinOptions.SecurePort)
		if err := leaderRuntime.RuntimeAuthentication.NetWorkManagerClient.GenerateOrSkip(); err != nil {
			return err
		}

		if !leaderRuntime.RuntimeAuthentication.NetWorkManagerClient.Check() {
			return fmt.Errorf("fail to load network-manager TLS args")
		}

		// node token
		if nodeToken, err := leaderRuntime.RuntimeAuthentication.NetWorkManagerClient.Nodetoken(); err != nil {
			return err
		} else {
			new.NodeToken = nodeToken
		}

		// cert
		// join
		new.JoinOptions.CACert = leaderRuntime.RuntimeAuthentication.NetWorkManagerClient.JoinCACert
		new.JoinOptions.ClientCertFile = leaderRuntime.RuntimeAuthentication.NetWorkManagerClient.JoinClientCert
		new.JoinOptions.ClientkeyFile = leaderRuntime.RuntimeAuthentication.NetWorkManagerClient.JoinClientkey

		// register
		new.RegisterOptions.CACert = leaderRuntime.RuntimeAuthentication.NetWorkManagerClient.RegisterCACert
		new.RegisterOptions.ClientCertFile = leaderRuntime.RuntimeAuthentication.NetWorkManagerClient.RegisterClientCert
		new.RegisterOptions.ClientkeyFile = leaderRuntime.RuntimeAuthentication.NetWorkManagerClient.RegisterClientkey

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
			leaderRuntime.RuntimeAuthentication.NetWorkManagerClient = nil
			klog.Infof("network manager client certificates specified ok, ignore --token")
		} else {
			return fmt.Errorf("you have provide bad certificates or host info for network manager")
		}
	}

	return nil
}

func (leaderRuntime *LeaderRuntime) LoadApiserver() error {
	raw := leaderRuntime.FlagsOption.ApiserverOptions
	new := leaderRuntime.RuntimeOption.ApiserverOptions

	new.ReservedOptions = raw.ReservedOptions
	new.IgnoreOptions = raw.IgnoreOptions

	// load *LitekubeOptions
	new.Options.AllowPrivileged = raw.Options.AllowPrivileged
	new.Options.AuthorizationMode = raw.Options.AuthorizationMode
	new.Options.AnonymousAuth = raw.Options.AnonymousAuth
	new.Options.EnableSwaggerUI = raw.Options.EnableSwaggerUI
	new.Options.EnableAdmissionPlugins = raw.Options.EnableAdmissionPlugins
	new.Options.EncryptionProviderConfig = raw.Options.EncryptionProviderConfig
	new.Options.Profiling = raw.Options.Profiling
	// check --service-cluster-ip-range
	if _, _, err := net.ParseCIDR(raw.Options.ServiceClusterIpRange); err != nil {
		new.Options.ServiceClusterIpRange = apiserver.DefaultALO.ServiceClusterIpRange
		new.IgnoreOptions["service-cluster-ip-range"] = raw.Options.ServiceClusterIpRange
	} else {
		new.Options.ServiceClusterIpRange = raw.Options.ServiceClusterIpRange
	}
	// check --service-node-port-range
	new.Options.ServiceNodePortRange = ""
	if ports := strings.Split(raw.Options.ServiceNodePortRange, "-"); len(ports) == 2 {
		port_min := 30000
		port_max := 65535
		if p, err := strconv.Atoi(strings.TrimSpace(ports[0])); err == nil && p > 0 {
			port_min = p
		}

		if p, err := strconv.Atoi(strings.TrimSpace(ports[1])); err == nil && p < 65536 {
			port_max = p
		}

		if port_max > port_min && (port_max-port_min) > 100 {
			new.Options.ServiceNodePortRange = fmt.Sprintf("%d-%d", port_min, port_max)
		}
	}

	if new.Options.ServiceNodePortRange == "" {
		// fail to parse port before
		new.Options.ServiceNodePortRange = apiserver.DefaultALO.ServiceNodePortRange
		//new.IgnoreOptions["service-node-port-range"] = raw.Options.ServiceNodePortRange
	}

	// secure-port
	if raw.Options.SecurePort < 1 || raw.Options.SecurePort > 65535 {
		new.Options.SecurePort = apiserver.DefaultALO.SecurePort
	} else {
		new.Options.SecurePort = raw.Options.SecurePort
	}

	// load *ProfessionalOptions
	// bind-address
	if ip := net.ParseIP(raw.ProfessionalOptions.BindAddress); ip == nil {
		new.ProfessionalOptions.BindAddress = apiserver.DefaultAPO.BindAddress
	} else {
		new.ProfessionalOptions.BindAddress = raw.ProfessionalOptions.BindAddress
	}

	// advertiseAddress
	if ip := net.ParseIP(raw.ProfessionalOptions.AdvertiseAddress); ip == nil {
		// no value util run network-manager client
		if remoteIp, err := leaderRuntime.NetworkRegisterClient.QueryIp(); err != nil {
			return err
		} else {
			new.ProfessionalOptions.AdvertiseAddress = remoteIp
		}
	} else {
		new.ProfessionalOptions.AdvertiseAddress = raw.ProfessionalOptions.AdvertiseAddress
	}

	// InsecurePort
	if raw.ProfessionalOptions.InsecurePort < 1 || raw.ProfessionalOptions.InsecurePort > 65535 {
		new.ProfessionalOptions.InsecurePort = apiserver.DefaultAPO.InsecurePort
	} else {
		new.ProfessionalOptions.InsecurePort = raw.ProfessionalOptions.InsecurePort
	}

	// FeatureGates
	new.ProfessionalOptions.FeatureGates = raw.ProfessionalOptions.FeatureGates
	if new.ProfessionalOptions.FeatureGates == "" {
		new.ProfessionalOptions.FeatureGates = apiserver.DefaultAPO.FeatureGates
	}

	// ECTDOptions
	new.ProfessionalOptions.StorageBackend = raw.ProfessionalOptions.StorageBackend
	if new.ProfessionalOptions.StorageBackend == "" {
		new.ProfessionalOptions.StorageBackend = apiserver.DefaultEO.StorageBackend
	}
	// address
	new.ProfessionalOptions.EtcdServers = raw.ProfessionalOptions.EtcdServers
	if new.ProfessionalOptions.EtcdServers == "" {
		new.ProfessionalOptions.EtcdServers = apiserver.DefaultEO.EtcdServers
	}
	// certificate
	if certificate.NotExists(raw.ProfessionalOptions.EtcdCertfile, raw.ProfessionalOptions.EtcdKeyfile, raw.ProfessionalOptions.EtcdCafile) {
		if leaderRuntime.RuntimeOption.KineOptions != nil && leaderRuntime.RuntimeOption.GlobalOptions.RunKine {
			new.ProfessionalOptions.EtcdCafile = leaderRuntime.RuntimeAuthentication.Kine.CACert
			new.ProfessionalOptions.EtcdCertfile = leaderRuntime.RuntimeAuthentication.Kine.ClientCert
			new.ProfessionalOptions.EtcdKeyfile = leaderRuntime.RuntimeAuthentication.Kine.Clientkey
		} else {
			return fmt.Errorf("no Etcd certificate")
		}
	} else {
		if !certificate.ValidateTLSPair(raw.ProfessionalOptions.EtcdCertfile, raw.ProfessionalOptions.EtcdKeyfile) {
			klog.Errorf("you specified an invaild certificate for ETCD Client")
			return fmt.Errorf("you specified an invaild certificate for ETCD Client")
		}

		new.ProfessionalOptions.EtcdCafile = raw.ProfessionalOptions.EtcdCafile
		new.ProfessionalOptions.EtcdCertfile = raw.ProfessionalOptions.EtcdCertfile
		new.ProfessionalOptions.EtcdKeyfile = raw.ProfessionalOptions.EtcdKeyfile
	}

	// server cert options
	new.ProfessionalOptions.ApiAudiences = raw.ProfessionalOptions.ApiAudiences
	if new.ProfessionalOptions.ApiAudiences == "" {
		new.ProfessionalOptions.ApiAudiences = apiserver.DefaultSCO.ApiAudiences
	}
	new.ProfessionalOptions.EnableBootstrapTokenAuth = raw.ProfessionalOptions.EnableBootstrapTokenAuth
	new.ProfessionalOptions.ServiceAccountIssuer = raw.ProfessionalOptions.ServiceAccountIssuer
	if new.ProfessionalOptions.ServiceAccountIssuer == "" {
		new.ProfessionalOptions.ServiceAccountIssuer = apiserver.DefaultSCO.ServiceAccountIssuer
	}
	new.ProfessionalOptions.RequestheaderAllowedNames = raw.ProfessionalOptions.RequestheaderAllowedNames
	if new.ProfessionalOptions.RequestheaderAllowedNames == "" {
		new.ProfessionalOptions.RequestheaderAllowedNames = apiserver.DefaultSCO.RequestheaderAllowedNames
	}
	new.ProfessionalOptions.RequestheaderExtraHeadersPrefix = raw.ProfessionalOptions.RequestheaderExtraHeadersPrefix
	if new.ProfessionalOptions.RequestheaderExtraHeadersPrefix == "" {
		new.ProfessionalOptions.RequestheaderExtraHeadersPrefix = apiserver.DefaultSCO.RequestheaderExtraHeadersPrefix
	}
	new.ProfessionalOptions.RequestheaderGroupHeaders = raw.ProfessionalOptions.RequestheaderGroupHeaders
	if new.ProfessionalOptions.RequestheaderGroupHeaders == "" {
		new.ProfessionalOptions.RequestheaderGroupHeaders = apiserver.DefaultSCO.RequestheaderGroupHeaders
	}
	new.ProfessionalOptions.RequestheaderUsernameHeaders = raw.ProfessionalOptions.RequestheaderUsernameHeaders
	if new.ProfessionalOptions.RequestheaderUsernameHeaders == "" {
		new.ProfessionalOptions.RequestheaderUsernameHeaders = apiserver.DefaultSCO.RequestheaderUsernameHeaders
	}
	new.ProfessionalOptions.RequestheaderAllowedNames = raw.ProfessionalOptions.RequestheaderAllowedNames
	if new.ProfessionalOptions.RequestheaderAllowedNames == "" {
		new.ProfessionalOptions.RequestheaderAllowedNames = apiserver.DefaultSCO.RequestheaderAllowedNames
	}
	new.ProfessionalOptions.EnableAggregatorRouting = raw.ProfessionalOptions.EnableAggregatorRouting

	// generate certificates
	ka := authentication.NewKubernetesAuthentication(leaderRuntime.RuntimeAuthentication.CertDir, new)
	leaderRuntime.RuntimeAuthentication.Kubernetes = ka
	if ka == nil {
		return fmt.Errorf("fail to generate certificates for kubernetes")
	}

	if err := ka.GenerateOrSkip(); err != nil {
		return err
	}

	if raw.ProfessionalOptions.CertDir != "" {
		ka.KubernetesTLSDir = raw.ProfessionalOptions.CertDir
	}
	if raw.ProfessionalOptions.TlsCertFile != "" {
		ka.ApiserverServerCert = raw.ProfessionalOptions.TlsCertFile
	}
	if raw.ProfessionalOptions.TlsPrivateKeyFile != "" {
		ka.ApiserverServerKey = raw.ProfessionalOptions.TlsPrivateKeyFile
	}
	if raw.ProfessionalOptions.TokenAuthFile != "" {
		ka.TokenAuthFile = raw.ProfessionalOptions.TokenAuthFile
	}
	if raw.ProfessionalOptions.ServiceAccountSigningKeyFile != "" {
		ka.ServiceKeyPair = raw.ProfessionalOptions.ServiceAccountSigningKeyFile
	}
	if raw.ProfessionalOptions.ServiceAccountKeyFile != "" {
		new.ProfessionalOptions.ServiceAccountKeyFile = raw.ProfessionalOptions.ServiceAccountKeyFile
	} else {
		new.ProfessionalOptions.ServiceAccountKeyFile = ka.ServiceKeyPair
	}
	if raw.ProfessionalOptions.ClientCAFile != "" {
		ka.ApiserverValidateClientsCA = raw.ProfessionalOptions.ClientCAFile
	}
	if raw.ProfessionalOptions.RequestheaderClientCAFile != "" {
		ka.ApiserverRequestHeaderCA = raw.ProfessionalOptions.RequestheaderClientCAFile
	}
	if raw.ProfessionalOptions.ProxyClientCertFile != "" {
		ka.ApiserverClientAuthProxyCert = raw.ProfessionalOptions.ProxyClientCertFile
	}
	if raw.ProfessionalOptions.ProxyClientKeyFile != "" {
		ka.ApiserverClientAuthProxyKey = raw.ProfessionalOptions.ProxyClientKeyFile
	}
	if raw.ProfessionalOptions.KubeletCertificateAuthority != "" {
		ka.ApiserverValidateKubeletServerCA = raw.ProfessionalOptions.KubeletCertificateAuthority
	}
	if raw.ProfessionalOptions.KubeletClientCertificate != "" {
		ka.ApiserverClientKubeletCert = raw.ProfessionalOptions.KubeletClientCertificate
	}
	if raw.ProfessionalOptions.KubeletClientKey != "" {
		ka.ApiserverClientKubeletKey = raw.ProfessionalOptions.KubeletClientKey
	}

	new.ProfessionalOptions.CertDir = ka.KubernetesTLSDir
	new.ProfessionalOptions.TlsCertFile = ka.ApiserverServerCert
	new.ProfessionalOptions.TlsPrivateKeyFile = ka.ApiserverServerKey
	new.ProfessionalOptions.TokenAuthFile = ka.TokenAuthFile
	new.ProfessionalOptions.ServiceAccountSigningKeyFile = ka.ServiceKeyPair
	new.ProfessionalOptions.ClientCAFile = ka.ApiserverValidateClientsCA
	new.ProfessionalOptions.RequestheaderClientCAFile = ka.ApiserverRequestHeaderCA
	new.ProfessionalOptions.ProxyClientCertFile = ka.ApiserverClientAuthProxyCert
	new.ProfessionalOptions.ProxyClientKeyFile = ka.ApiserverClientAuthProxyKey
	new.ProfessionalOptions.KubeletCertificateAuthority = ka.ApiserverValidateKubeletServerCA
	new.ProfessionalOptions.KubeletClientCertificate = ka.ApiserverClientKubeletCert
	new.ProfessionalOptions.KubeletClientKey = ka.ApiserverClientKubeletKey

	if err := os.MkdirAll(ka.KubernetesTLSDir, os.FileMode(0644)); err != nil {
		return err
	}

	return nil
}

func (leaderRuntime *LeaderRuntime) LoadControllermanager() error {
	raw := leaderRuntime.FlagsOption.ControllerManagerOptions
	new := leaderRuntime.RuntimeOption.ControllerManagerOptions
	ka := leaderRuntime.RuntimeAuthentication.Kubernetes

	new.ReservedOptions = raw.ReservedOptions
	new.IgnoreOptions = raw.IgnoreOptions

	// load *LitekubeOptions
	new.Options.AllocateNodeCidrs = raw.Options.AllocateNodeCidrs
	if _, _, err := net.ParseCIDR(raw.Options.ClusterCidr); err != nil {
		new.Options.ClusterCidr = controllermanager.DefaultCMLO.ClusterCidr
		//new.IgnoreOptions["service-cluster-ip-range"] = raw.Options.ServiceClusterIpRange
	} else {
		new.Options.ClusterCidr = raw.Options.ClusterCidr
	}
	new.Options.Profiling = raw.Options.Profiling
	new.Options.UseServiceAccountCredentials = raw.Options.UseServiceAccountCredentials

	// load *ProfessionalOptions
	// bind-address
	if ip := net.ParseIP(raw.ProfessionalOptions.BindAddress); ip == nil {
		new.ProfessionalOptions.BindAddress = controllermanager.DefaultCMPO.BindAddress
	} else {
		new.ProfessionalOptions.BindAddress = raw.ProfessionalOptions.BindAddress
	}

	// secure-port
	if raw.ProfessionalOptions.SecurePort < 1 || raw.ProfessionalOptions.SecurePort > 65535 {
		new.ProfessionalOptions.SecurePort = controllermanager.DefaultCMPO.SecurePort
	} else {
		new.ProfessionalOptions.SecurePort = raw.ProfessionalOptions.SecurePort
	}

	new.ProfessionalOptions.LeaderElect = raw.ProfessionalOptions.LeaderElect
	new.ProfessionalOptions.ConfigureCloudRoutes = raw.ProfessionalOptions.ConfigureCloudRoutes
	new.ProfessionalOptions.Controllers = raw.ProfessionalOptions.Controllers
	if new.ProfessionalOptions.Controllers == "" {
		new.ProfessionalOptions.Controllers = controllermanager.DefaultCMPO.Controllers
	}
	new.ProfessionalOptions.FeatureGates = raw.ProfessionalOptions.FeatureGates
	if new.ProfessionalOptions.FeatureGates == "" {
		new.ProfessionalOptions.FeatureGates = controllermanager.DefaultCMPO.FeatureGates
	}

	// cert
	new.ProfessionalOptions.Kubeconfig = ka.KubeConfigController
	new.ProfessionalOptions.AuthenticationKubeconfig = ka.KubeConfigController
	new.ProfessionalOptions.AuthorizationKubeconfig = ka.KubeConfigController
	new.ProfessionalOptions.ServiceAccountPrivateKeyFile = ka.ServiceKeyPair
	new.ProfessionalOptions.RootCaFile = ka.ClusterValidateServerCA
	new.ProfessionalOptions.ClusterSigningKubeApiserverClientCertFile = ka.KubeletValidateApiserverClientCA
	new.ProfessionalOptions.ClusterSigningKubeApiserverClientKeyFile = ka.KubeletValidateApiserverClientCAKey
	new.ProfessionalOptions.ClusterSigningKubeletClientCertFile = ka.ApiserverValidateClientsCA
	new.ProfessionalOptions.ClusterSigningKubeletClientKeyFile = ka.ApiserverValidateClientsCAKey
	new.ProfessionalOptions.ClusterSigningKubeletServingCertFile = ka.ApiserverValidateKubeletServerCA
	new.ProfessionalOptions.ClusterSigningKubeletServingKeyFile = ka.ApiserverValidateKubeletServerCAKey
	new.ProfessionalOptions.ClusterSigningLegacyUnknownCertFile = ka.ApiserverValidateClientsCA
	new.ProfessionalOptions.ClusterSigningLegacyUnknownKeyFile = ka.ApiserverValidateClientsCAKey

	if raw.ProfessionalOptions.Kubeconfig != "" {
		new.ProfessionalOptions.Kubeconfig = raw.ProfessionalOptions.Kubeconfig
	}
	if raw.ProfessionalOptions.AuthenticationKubeconfig != "" {
		new.ProfessionalOptions.AuthenticationKubeconfig = raw.ProfessionalOptions.AuthenticationKubeconfig
	}
	if raw.ProfessionalOptions.AuthorizationKubeconfig != "" {
		new.ProfessionalOptions.AuthorizationKubeconfig = raw.ProfessionalOptions.AuthorizationKubeconfig
	}
	if raw.ProfessionalOptions.ServiceAccountPrivateKeyFile != "" {
		new.ProfessionalOptions.ServiceAccountPrivateKeyFile = raw.ProfessionalOptions.ServiceAccountPrivateKeyFile
	}
	if raw.ProfessionalOptions.RootCaFile != "" {
		new.ProfessionalOptions.RootCaFile = raw.ProfessionalOptions.RootCaFile
	}
	if raw.ProfessionalOptions.ClusterSigningKubeApiserverClientCertFile != "" {
		new.ProfessionalOptions.ClusterSigningKubeApiserverClientCertFile = raw.ProfessionalOptions.ClusterSigningKubeApiserverClientCertFile
	}
	if raw.ProfessionalOptions.ClusterSigningKubeApiserverClientKeyFile != "" {
		new.ProfessionalOptions.ClusterSigningKubeApiserverClientKeyFile = raw.ProfessionalOptions.ClusterSigningKubeApiserverClientKeyFile
	}
	if raw.ProfessionalOptions.ClusterSigningKubeletClientCertFile != "" {
		new.ProfessionalOptions.ClusterSigningKubeletClientCertFile = raw.ProfessionalOptions.ClusterSigningKubeletClientCertFile
	}
	if raw.ProfessionalOptions.ClusterSigningKubeApiserverClientKeyFile != "" {
		new.ProfessionalOptions.ClusterSigningKubeApiserverClientKeyFile = raw.ProfessionalOptions.ClusterSigningKubeApiserverClientKeyFile
	}
	if raw.ProfessionalOptions.ClusterSigningKubeletServingCertFile != "" {
		new.ProfessionalOptions.ClusterSigningKubeletServingCertFile = raw.ProfessionalOptions.ClusterSigningKubeletServingCertFile
	}
	if raw.ProfessionalOptions.ClusterSigningKubeletServingKeyFile != "" {
		new.ProfessionalOptions.ClusterSigningKubeletServingKeyFile = raw.ProfessionalOptions.ClusterSigningKubeletServingKeyFile
	}
	if raw.ProfessionalOptions.ClusterSigningLegacyUnknownCertFile != "" {
		new.ProfessionalOptions.ClusterSigningLegacyUnknownCertFile = raw.ProfessionalOptions.ClusterSigningLegacyUnknownCertFile
	}
	if raw.ProfessionalOptions.ClusterSigningLegacyUnknownKeyFile != "" {
		new.ProfessionalOptions.ClusterSigningLegacyUnknownKeyFile = raw.ProfessionalOptions.ClusterSigningLegacyUnknownKeyFile
	}

	return nil
}

func (leaderRuntime *LeaderRuntime) LoadScheduler() error {
	raw := leaderRuntime.FlagsOption.SchedulerOptions
	new := leaderRuntime.RuntimeOption.SchedulerOptions
	ka := leaderRuntime.RuntimeAuthentication.Kubernetes

	new.ReservedOptions = raw.ReservedOptions
	new.IgnoreOptions = raw.IgnoreOptions

	// load * LitekubeOptions
	new.Options.Profiling = raw.Options.Profiling

	// load * ProfessionalOptions
	// bind-address
	if ip := net.ParseIP(raw.ProfessionalOptions.BindAddress); ip == nil {
		new.ProfessionalOptions.BindAddress = scheduler.DefaultSPO.BindAddress
	} else {
		new.ProfessionalOptions.BindAddress = raw.ProfessionalOptions.BindAddress
	}

	// secure-port
	if raw.ProfessionalOptions.SecurePort < 1 || raw.ProfessionalOptions.SecurePort > 65535 {
		new.ProfessionalOptions.SecurePort = scheduler.DefaultSPO.SecurePort
	} else {
		new.ProfessionalOptions.SecurePort = raw.ProfessionalOptions.SecurePort
	}

	new.ProfessionalOptions.LeaderElect = raw.ProfessionalOptions.LeaderElect
	new.ProfessionalOptions.KubeConfig = ka.KubeConfigScheduler
	new.ProfessionalOptions.AuthenticationKubeconfig = ka.KubeConfigScheduler
	new.ProfessionalOptions.AuthorizationKubeconfig = ka.KubeConfigScheduler

	if raw.ProfessionalOptions.KubeConfig != "" {
		new.ProfessionalOptions.KubeConfig = raw.ProfessionalOptions.KubeConfig
	}
	if raw.ProfessionalOptions.AuthenticationKubeconfig != "" {
		new.ProfessionalOptions.AuthenticationKubeconfig = raw.ProfessionalOptions.AuthenticationKubeconfig
	}
	if raw.ProfessionalOptions.AuthorizationKubeconfig != "" {
		new.ProfessionalOptions.AuthorizationKubeconfig = raw.ProfessionalOptions.AuthorizationKubeconfig
	}

	return nil
}

// func (leaderRuntime *LeaderRuntime) LoadWorker() error {
// 	if !leaderRuntime.FlagsOption.GlobalOptions.EnableWorker {
// 		return nil
// 	}

// 	return nil
// }
