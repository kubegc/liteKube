package config

import (
	"path/filepath"

	leaderAuth "github.com/litekube/LiteKube/pkg/leader/authentication"
	globaloptions "github.com/litekube/LiteKube/pkg/options/worker/global"
	"github.com/litekube/LiteKube/pkg/worker/authentication"
)

type RuntimeAuthentications struct {
	CertDir              string
	NetWorkManagerClient *leaderAuth.NetworkControllerClientAuthentication // nil if user provide certificate
	KubernetesNode       *authentication.KubernetesNode
}

func NewRuntimeAuthentication(rootCertPath string) *RuntimeAuthentications {
	if rootCertPath == "" {
		rootCertPath = filepath.Join(globaloptions.DefaultGO.WorkDir, "tls/")
	}
	return &RuntimeAuthentications{
		CertDir:              rootCertPath,
		NetWorkManagerClient: nil,
		KubernetesNode:       nil,
	}
}
