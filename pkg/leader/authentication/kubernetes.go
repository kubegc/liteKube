package authentication

import (
	"path/filepath"

	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
)

type KubernetesAuthentication struct {
	KubernetesCertDir string
}

func NewKubernetesAuthentication(rootCertPath string) *KubernetesAuthentication {
	if rootCertPath == "" {
		rootCertPath = filepath.Join(globaloptions.DefaultGO.WorkDir, "tls/")
	}
	return &KubernetesAuthentication{
		KubernetesCertDir: filepath.Join(rootCertPath, "kubernetes/"),
	}
}
