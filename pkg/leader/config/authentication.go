package config

import (
	"path/filepath"

	"github.com/litekube/LiteKube/pkg/leader/authentication"
	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
)

type RuntimeAuthentications struct {
	CertDir    string
	Kine       *authentication.KineAuthentication
	Kubernetes *authentication.KubernetesAuthentication
}

func NewRuntimeAuthentication(rootCertPath string) *RuntimeAuthentications {
	if rootCertPath == "" {
		rootCertPath = filepath.Join(globaloptions.DefaultGO.WorkDir, "tls/")
	}
	return &RuntimeAuthentications{
		CertDir:    rootCertPath,
		Kine:       nil,
		Kubernetes: nil,
	}
}
