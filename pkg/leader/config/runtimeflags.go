package config

import (
	"fmt"
	"net"
	"path/filepath"

	"github.com/litekube/LiteKube/pkg/certificate"
	"github.com/litekube/LiteKube/pkg/leader/authentication"
	options "github.com/litekube/LiteKube/pkg/options/leader"
	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
	kineoptions "github.com/litekube/LiteKube/pkg/options/leader/kine"
)

func (runtime *LeaderRuntime) SetFlags(opt *options.LeaderOptions) {
	runtime.FlagsOption = opt
}

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

	//runtime.runtimeAuthentication = NewRuntimeAuthentication(opt.GlobalOptions.WorkDir)
	return nil
}

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

	new.EnableWorker = raw.EnableWorker

	if !new.EnableWorker {
		new.WorkerConfig = ""
	}
	return nil
}

func (runtime *LeaderRuntime) LoadKine() error {
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

	// check file path
	if !certificate.Exists(raw.CACert, raw.ServerCertFile, raw.ServerkeyFile) || !certificate.ValidateTLSPair(raw.ServerCertFile, raw.ServerkeyFile) {
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
		new.CACert = raw.CACert
		new.ServerCertFile = raw.ServerCertFile
		new.ServerkeyFile = raw.ServerkeyFile
	}

	return nil
}
