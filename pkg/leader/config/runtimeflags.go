package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"github.com/litekube/LiteKube/pkg/certificate"
	"github.com/litekube/LiteKube/pkg/leader/authentication"
	options "github.com/litekube/LiteKube/pkg/options/leader"
	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
	kineoptions "github.com/litekube/LiteKube/pkg/options/leader/kine"
	"github.com/litekube/LiteKube/pkg/options/leader/netmanager"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

func (runtime *LeaderRuntime) SetFlags(opt *options.LeaderOptions) {
	runtime.FlagsOption = opt
}

// load all flags
func (runtime *LeaderRuntime) LoadFlags() error {
	if runtime.FlagsOption == nil {
		return fmt.Errorf("no flags input")
	}

	// init global flags
	if err := runtime.LoadGloabl(); err != nil {
		return err
	}

	// init kine flags
	if err := runtime.LoadKine(); err != nil {
		return err
	}

	// init network manager flags
	if err := runtime.LoadNetManager(); err != nil {
		return err
	}

	if config, err := yaml.Marshal(runtime.RuntimeOption.LeaderOptions); err != nil {
		return err
	} else {
		if e := ioutil.WriteFile(filepath.Join(runtime.RuntimeOption.GlobalOptions.WorkDir, "config.yaml"), config, os.ModePerm); e != nil {
			return e
		}
	}
	return nil
}

// load or generate args for litekube-global
func (runtime *LeaderRuntime) LoadGloabl() error {
	defer func() {
		runtime.RuntimeAuthentication = NewRuntimeAuthentication(filepath.Join(runtime.RuntimeOption.GlobalOptions.WorkDir, "tls/"))
	}()

	raw := runtime.FlagsOption.GlobalOptions
	new := runtime.RuntimeOption.GlobalOptions

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
	new.LogToStd = raw.LogToStd

	// kine
	new.RunKine = raw.RunKine
	// invalid etcd server will enable kine
	if runtime.RuntimeOption.ApiserverOptions.ProfessionalOptions.ECTDOptions.EtcdServers == "" {
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

// load or generate args for kine
func (runtime *LeaderRuntime) LoadKine() error {
	if !runtime.RuntimeOption.GlobalOptions.RunKine {
		runtime.RuntimeOption.KineOptions = nil
		return nil
	}

	raw := runtime.FlagsOption.KineOptions
	new := runtime.RuntimeOption.KineOptions

	// check bind-address
	if ip := net.ParseIP(raw.BindAddress); ip == nil {
		new.BindAddress = kineoptions.DefaultKO.BindAddress
	}

	// check https port
	if raw.SecurePort < 1 {
		new.SecurePort = kineoptions.DefaultKO.SecurePort
	}

	// check TLS certificate
	if certificate.NotExists(raw.CACert, raw.ServerCertFile, raw.ServerkeyFile) {
		klog.Info("built-in certificates for kine will be used")

		// invalid cert， generate kine certs
		runtime.RuntimeOption.OwnKineCert = true
		runtime.RuntimeAuthentication.Kine = authentication.NewKineAuthentication(runtime.RuntimeAuthentication.CertDir, new.BindAddress)
		if err := runtime.RuntimeAuthentication.Kine.GenerateOrSkip(); err != nil {
			return err
		}

		new.CACert = runtime.RuntimeAuthentication.Kine.CACert
		new.ServerCertFile = runtime.RuntimeAuthentication.Kine.ServerCert
		new.ServerkeyFile = runtime.RuntimeAuthentication.Kine.Serverkey
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
func (runtime *LeaderRuntime) LoadNetManager() error {
	raw := runtime.FlagsOption.NetmamagerOptions
	new := runtime.RuntimeOption.NetmamagerOptions

	// check bind-address
	if ip := net.ParseIP(raw.RegisterOptions.Address); ip == nil {
		new.RegisterOptions.Address = netmanager.DefaultRONO.Address
	}

	if ip := net.ParseIP(raw.JoinOptions.Address); ip == nil {
		new.JoinOptions.Address = netmanager.DefaultJONO.Address
	}

	// check https port
	if raw.RegisterOptions.SecurePort < 1 {
		new.RegisterOptions.SecurePort = netmanager.DefaultRONO.SecurePort
	}

	// check https port
	if raw.JoinOptions.SecurePort < 1 {
		new.JoinOptions.SecurePort = netmanager.DefaultJONO.SecurePort
	}

	// check Token
	if runtime.RuntimeOption.GlobalOptions.RunNetManager {
		// generate certificate for network manager
		new.Token = "local"
		runtime.RuntimeAuthentication.NetWorkManager = authentication.NewNetworkAuthentication(runtime.RuntimeAuthentication.CertDir, new.RegisterOptions.Address, new.JoinOptions.Address)
		if err := runtime.RuntimeAuthentication.NetWorkManager.GenerateOrSkip(); err != nil {
			return err
		}

	} else {
		if raw.Token == "local" {
			return fmt.Errorf("bad token to connect with network-manager")
		}
		new.Token = raw.Token
	}

	// check client certificate
	runtime.RuntimeAuthentication.NetWorkManagerClient = authentication.NewNetworkManagerClient(runtime.RuntimeAuthentication.CertDir, new.Token)
	if err := runtime.RuntimeAuthentication.NetWorkManagerClient.GenerateOrSkip(new.RegisterOptions.Address, int(new.RegisterOptions.SecurePort)); err != nil {
		return err
	}

	if !runtime.RuntimeAuthentication.NetWorkManagerClient.Check() {
		return fmt.Errorf("fail to load network-manager TLS args")
	}

	// node token
	if nodeToken, err := runtime.RuntimeAuthentication.NetWorkManagerClient.Nodetoken(); err != nil {
		return err
	} else {
		new.NodeToken = nodeToken
	}

	// cert
	// join
	new.JoinOptions.CACert = runtime.RuntimeAuthentication.NetWorkManagerClient.JoinCACert
	new.JoinOptions.ClientCertFile = runtime.RuntimeAuthentication.NetWorkManagerClient.JoinClientCert
	new.JoinOptions.ClientkeyFile = runtime.RuntimeAuthentication.NetWorkManagerClient.JoinClientkey

	// register
	new.RegisterOptions.CACert = runtime.RuntimeAuthentication.NetWorkManagerClient.RegisterCACert
	new.RegisterOptions.ClientCertFile = runtime.RuntimeAuthentication.NetWorkManagerClient.RegisterClientCert
	new.RegisterOptions.ClientkeyFile = runtime.RuntimeAuthentication.NetWorkManagerClient.RegisterClientkey

	return nil
}
