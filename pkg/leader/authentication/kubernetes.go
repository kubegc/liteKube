package authentication

import (
	"path"
	"path/filepath"
	"text/template"

	"github.com/litekube/LiteKube/pkg/options/leader/apiserver"
	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
)

var (
	KubeconfigTemplate = template.Must(template.New("kubeconfig").Parse(`apiVersion: v1
clusters:
- cluster:
    server: {{.URL}}
    certificate-authority: {{.CACert}}
  name: local
contexts:
- context:
    cluster: local
    namespace: default
    user: user
  name: Default
current-context: Default
kind: Config
preferences: {}
users:
- name: user
  user:
    client-certificate: {{.ClientCert}}
    client-key: {{.ClientKey}}
`))
)

type KubernetesAuthentication struct {
	KubernetesRootDir         string
	KubernetesCertDir         string
	KubernetesKubeDir         string
	CheckFile                 string
	RequestheaderAllowedNames string

	ApiserverValidateClientsCA          string
	ApiserverValidateClientsCAKey       string
	ClusterValidateServerCA             string
	ClusterValidateServerCAKey          string
	ApiserverRequestHeaderCA            string
	ApiserverRequestHeaderCAKey         string
	IPSECKey                            string
	ServiceKey                          string
	PasswdFile                          string
	KubeletValidateApiserverClientCA    string
	KubeletValidateApiserverClientCAKey string
	ApiserverValidateKubeletServerCA    string
	ApiserverValidateKubeletServerCAKey string

	NodePasswdFile string

	KubeConfigAdmin              string
	KubeConfigController         string
	KubeConfigScheduler          string
	KubeConfigApiserverToKubelet string
	KubeConfigCloudController    string

	ApiserverClientKubeletCA     string
	ClientAdminCert              string
	ClientAdminKey               string
	ClientControllerCert         string
	ClientControllerKey          string
	ClientCloudControllerCert    string
	ClientCloudControllerKey     string
	ClientSchedulerCert          string
	ClientSchedulerKey           string
	ApiserverClientKubeletCert   string
	ApiserverClientKubeletKey    string
	ClientKubeProxyCert          string
	ClientKubeProxyKey           string
	ClientLitekubeControllerCert string
	ClientLitekubeControllerKey  string

	ApiserverServerCert string
	ApiserverServerKey  string

	ClientKubeletKey  string
	ServingKubeletKey string

	ApiserverClientAuthProxyCert string
	ApiserverClientAuthProxyKey  string
}

func NewKubernetesAuthentication(rootCertPath string, requestheaderAllowedNames string) *KubernetesAuthentication {
	if rootCertPath == "" {
		rootCertPath = filepath.Join(globaloptions.DefaultGO.WorkDir, "tls/")
	}

	kubernetesRootDir := filepath.Join(rootCertPath, "kubernetes/")
	kubernetesCertDir := filepath.Join(kubernetesRootDir, "cert/")
	kubernetesKubeDir := filepath.Join(kubernetesRootDir, "kube/")

	if requestheaderAllowedNames == "" {
		requestheaderAllowedNames = apiserver.DefaultSCO.RequestheaderAllowedNames
	}

	return &KubernetesAuthentication{
		KubernetesRootDir:         kubernetesRootDir,
		KubernetesCertDir:         kubernetesCertDir,
		KubernetesKubeDir:         kubernetesKubeDir,
		CheckFile:                 filepath.Join(kubernetesRootDir, ".valid"),
		RequestheaderAllowedNames: requestheaderAllowedNames,

		// kube-apiserver certificates

		// key or certificate file for kube-apiserver to validate service-account-token
		ServiceKey: path.Join(kubernetesCertDir, "kube-apiserver", "service", "service-account.key"), // --service-account-key-file

		// kube-apiserver CA certificate for validate others
		ApiserverValidateClientsCA:    path.Join(kubernetesCertDir, "ca", "apiserver-client.crt"), // --client-ca-file
		ApiserverValidateClientsCAKey: path.Join(kubernetesCertDir, "ca", "apiserver-client.key"),
		ClusterValidateServerCA:       path.Join(kubernetesCertDir, "ca", "cluster-server.crt"), // --kubelet-certificate-authority
		ClusterValidateServerCAKey:    path.Join(kubernetesCertDir, "ca", "cluster-server.key"),
		// request-header CA
		ApiserverRequestHeaderCA:            path.Join(kubernetesCertDir, "ca", "kube-apiserver-auth-proxy.crt"), // --requestheader-client-ca-file
		ApiserverRequestHeaderCAKey:         path.Join(kubernetesCertDir, "ca", "kube-apiserver-auth-proxy.key"),
		KubeletValidateApiserverClientCA:    path.Join(kubernetesCertDir, "ca", "kubelet-apiserver-client.crt"),
		KubeletValidateApiserverClientCAKey: path.Join(kubernetesCertDir, "ca", "kubelet-apiserver-client.key"),
		ApiserverValidateKubeletServerCA:    path.Join(kubernetesCertDir, "ca", "apiserver-kubelet-server.crt"),
		ApiserverValidateKubeletServerCAKey: path.Join(kubernetesCertDir, "ca", "apiserver-kubelet-server.key"),

		// kube-apiserver certificate for service others
		ApiserverServerCert: path.Join(kubernetesCertDir, "kube-apiserver", "server", "server.crt"), // --tls-cert-file
		ApiserverServerKey:  path.Join(kubernetesCertDir, "kube-apiserver", "server", "server.key"), // --tls-private-key-file
		// kube-apiserver certificate for access kubelet
		ApiserverClientKubeletCert:   path.Join(kubernetesCertDir, "kube-apiserver", "client", "client.crt"), // --kubelet-client-certificate
		ApiserverClientKubeletKey:    path.Join(kubernetesCertDir, "kube-apiserver", "client", "client.key"), // --kubelet-client-key
		KubeConfigApiserverToKubelet: path.Join(kubernetesCertDir, "kube-apiserver", "client", "api-server.kubeconfig"),

		// another way to access kube-apiserver: by proxy.
		// client can access proxy-server by http, proxy-server add certificate information before access kube-apiserver automatically
		// the follow are need:
		// --requestheader-client-ca-file=<path to aggregator CA cert>
		// --requestheader-allowed-names=front-proxy-client
		// --requestheader-extra-headers-prefix=X-Remote-Extra-
		// --requestheader-group-headers=X-Remote-Group
		// --requestheader-username-headers=X-Remote-User
		// --proxy-client-cert-file=<path to aggregator proxy cert>
		// --proxy-client-key-file=<path to aggregator proxy key>
		// --enable-aggregator-routing=true // while proxy is not in same machine

		ApiserverClientAuthProxyCert: path.Join(kubernetesCertDir, "kube-apiserver", "auth-proxy", "auth-proxy.crt"), // --proxy-client-cert-file
		ApiserverClientAuthProxyKey:  path.Join(kubernetesCertDir, "kube-apiserver", "auth-proxy", "auth-proxy.key"), // --proxy-client-key-file

		// admin
		ClientAdminCert: path.Join(kubernetesCertDir, "admin", "admin.crt"),
		ClientAdminKey:  path.Join(kubernetesCertDir, "admin", "admin.key"),
		KubeConfigAdmin: path.Join(kubernetesCertDir, "admin", "admin.kubeconfig"),

		// controller
		ClientControllerCert: path.Join(kubernetesCertDir, "controller", "controller.crt"),
		ClientControllerKey:  path.Join(kubernetesCertDir, "controller", "controller.key"),
		KubeConfigController: path.Join(kubernetesCertDir, "controller", "controller.kubeconfig"),

		// cloud-controller
		// ClientCloudControllerCert: path.Join(kubernetesCertDir, "cloud-controller", "client-cloud-controller.crt"),
		// ClientCloudControllerKey:  path.Join(kubernetesCertDir, "cloud-controller", "client-cloud-controller.key"),
		// KubeConfigCloudController: path.Join(kubernetesCertDir, "cloud-controller", "cloud-controller.kubeconfig"),

		// scheduler
		ClientSchedulerCert: path.Join(kubernetesCertDir, "scheduler", "scheduler.crt"),
		ClientSchedulerKey:  path.Join(kubernetesCertDir, "scheduler", "scheduler.key"),
		KubeConfigScheduler: path.Join(kubernetesCertDir, "scheduler", "scheduler.kubeconfig"),

		ClientKubeProxyCert: path.Join(kubernetesCertDir, "kube-proxy", "kube-proxy.crt"),
		ClientKubeProxyKey:  path.Join(kubernetesCertDir, "kube-proxy", "kube-proxy.key"),

		ClientLitekubeControllerCert: path.Join(kubernetesCertDir, "litekube-controller", "litekube-controller.crt"),
		ClientLitekubeControllerKey:  path.Join(kubernetesCertDir, "litekube-controller", "litekube-controller.key"),

		ClientKubeletKey:  path.Join(kubernetesCertDir, "kubelet", "client-kubelet.key"),
		ServingKubeletKey: path.Join(kubernetesCertDir, "kubelet", "serving-kubelet.key"),

		// cluster CA certificate for server like kubelet at port:10250 or kube-controller

		IPSECKey:   path.Join(kubernetesCertDir, "other", "ipsec.psk"),
		PasswdFile: path.Join(kubernetesCertDir, "other", "passwd"),

		NodePasswdFile: path.Join(kubernetesCertDir, "other", "node-passwd"),
	}
}

func (na *KubernetesAuthentication) GenerateOrSkip() error {

	return nil
}
