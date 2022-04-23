package authentication

import (
	"fmt"

	"io/fs"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/litekube/LiteKube/pkg/certificate"
	"github.com/litekube/LiteKube/pkg/global"
	"github.com/litekube/LiteKube/pkg/options/leader/apiserver"
	globaloptions "github.com/litekube/LiteKube/pkg/options/leader/global"
	token "github.com/litekube/LiteKube/pkg/token"
	certutil "github.com/rancher/dynamiclistener/cert"
)

type signedCertFactory func(commonName string, organization []string, certFile, keyFile string) (bool, error)

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
	KubectlPath               string
	ApiserverEndpoint         string
	ApiserverEndpointIp       string
	ApiserverEndpointPort     uint16
	ApiServerServiceIP        string
	ServiceClusterIpRange     string
	KubernetesRootDir         string
	KubernetesCertDir         string
	KubernetesTLSDir          string
	KubernetesKubeDir         string
	RequestheaderAllowedNames string

	ServiceKeyPair string

	// apiserver as server:
	// --------------------start-------------------
	ApiserverValidateClientsCA    string
	ApiserverValidateClientsCAKey string
	ClusterValidateServerCA       string
	ClusterValidateServerCAKey    string
	ApiserverServerCert           string
	ApiserverServerKey            string
	// clients to access apiserver
	// admin
	AdminClientCert string
	AdminClientKey  string
	KubeConfigAdmin string
	// controller
	ControllerClientCert string
	ControllerClientKey  string
	KubeConfigController string
	// cloud-controller
	// ClientCloudControllerCert string
	// ClientCloudControllerKey string
	// KubeConfigCloudController string
	// scheduler
	SchedulerClientCert string
	SchedulerClientKey  string
	KubeConfigScheduler string
	// kube-proxy
	KubeProxyClientCert string
	KubeProxyClientKey  string
	// litekube
	LitekubeControllerClientCert string
	LitekubeControllerClientKey  string
	// ---------------------end--------------------

	// apiserver as client:
	// --------------------start-------------------
	// apiserver as client to external-aggregation-apiserver
	ApiserverRequestHeaderCA     string
	ApiserverRequestHeaderCAKey  string
	ApiserverClientAuthProxyCert string
	ApiserverClientAuthProxyKey  string
	// apiserver as client to kubelet
	KubeletValidateApiserverClientCA    string
	KubeletValidateApiserverClientCAKey string
	ApiserverValidateKubeletServerCA    string
	ApiserverValidateKubeletServerCAKey string
	ApiserverClientKubeletCert          string
	ApiserverClientKubeletKey           string
	KubeConfigApiserverToKubelet        string
	// ---------------------end--------------------

	TokenAuthFile string

	ClientKubeletKey  string
	ServingKubeletKey string

	IPSECKey   string
	PasswdFile string

	NodePasswdFile string
}

func NewKubernetesAuthentication(rootCertPath string, opt *apiserver.ApiserverOptions) *KubernetesAuthentication {
	if rootCertPath == "" {
		rootCertPath = filepath.Join(globaloptions.DefaultGO.WorkDir, "tls/")
	}

	kubernetesRootDir := filepath.Join(rootCertPath, "kubernetes/")
	kubernetesCertDir := filepath.Join(kubernetesRootDir, "cert/")

	requestheaderAllowedNames := opt.ProfessionalOptions.ServerCertOptions.RequestheaderAllowedNames
	if requestheaderAllowedNames == "" {
		requestheaderAllowedNames = apiserver.DefaultSCO.RequestheaderAllowedNames
	}

	_, clusterIpRange, err := net.ParseCIDR(opt.Options.ServiceClusterIpRange)
	if err != nil {
		return nil
	}

	serviceIp := global.GetDefaultServiceIp(clusterIpRange)
	if serviceIp == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Join(global.HomePath, ".kube/"), os.FileMode(0644)); err != nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Join(kubernetesRootDir, "tls/"), os.FileMode(0644)); err != nil {
		return nil
	}

	kubeDir := filepath.Join(global.HomePath, ".kube/")
	return &KubernetesAuthentication{
		KubectlPath:           filepath.Join(kubeDir, "config"),
		ApiserverEndpoint:     fmt.Sprintf("https://%s:%d", opt.ProfessionalOptions.AdvertiseAddress, opt.Options.SecurePort),
		ApiserverEndpointIp:   opt.ProfessionalOptions.AdvertiseAddress,
		ApiserverEndpointPort: opt.Options.SecurePort,
		ServiceClusterIpRange: opt.Options.ServiceClusterIpRange,
		ApiServerServiceIP:    serviceIp.To4().String(),
		KubernetesRootDir:     kubernetesRootDir,
		KubernetesCertDir:     kubernetesCertDir,
		KubernetesKubeDir:     kubeDir,
		KubernetesTLSDir:      filepath.Join(kubernetesRootDir, "tls/"),
		//CheckFile:                 filepath.Join(kubernetesRootDir, ".valid"),
		RequestheaderAllowedNames: requestheaderAllowedNames,

		// kube-apiserver certificates

		// key or certificate file for kube-apiserver to validate service-account-token
		ServiceKeyPair: path.Join(kubernetesCertDir, "kube-apiserver", "service", "service-account.key"),

		// kube-apiserver CA certificate for validate others
		ApiserverValidateClientsCA:    path.Join(kubernetesCertDir, "ca", "apiserver-client.crt"),
		ApiserverValidateClientsCAKey: path.Join(kubernetesCertDir, "ca", "apiserver-client.key"),
		ClusterValidateServerCA:       path.Join(kubernetesCertDir, "ca", "cluster-server.crt"),
		ClusterValidateServerCAKey:    path.Join(kubernetesCertDir, "ca", "cluster-server.key"),
		// request-header CA
		ApiserverRequestHeaderCA:            path.Join(kubernetesCertDir, "ca", "kube-apiserver-auth-proxy.crt"),
		ApiserverRequestHeaderCAKey:         path.Join(kubernetesCertDir, "ca", "kube-apiserver-auth-proxy.key"),
		KubeletValidateApiserverClientCA:    path.Join(kubernetesCertDir, "ca", "kubelet-apiserver-client.crt"),
		KubeletValidateApiserverClientCAKey: path.Join(kubernetesCertDir, "ca", "kubelet-apiserver-client.key"),
		ApiserverValidateKubeletServerCA:    path.Join(kubernetesCertDir, "ca", "apiserver-kubelet-server.crt"),
		ApiserverValidateKubeletServerCAKey: path.Join(kubernetesCertDir, "ca", "apiserver-kubelet-server.key"),

		// kube-apiserver certificate for service others
		ApiserverServerCert: path.Join(kubernetesCertDir, "kube-apiserver", "server", "server.crt"),
		ApiserverServerKey:  path.Join(kubernetesCertDir, "kube-apiserver", "server", "server.key"),
		// kube-apiserver certificate for access kubelet
		ApiserverClientKubeletCert:   path.Join(kubernetesCertDir, "kube-apiserver", "client", "client.crt"),
		ApiserverClientKubeletKey:    path.Join(kubernetesCertDir, "kube-apiserver", "client", "client.key"),
		KubeConfigApiserverToKubelet: path.Join(kubernetesCertDir, "kube-apiserver", "client", "api-server.kubeconfig"),
		TokenAuthFile:                path.Join(kubernetesCertDir, "kubelet", "token.csv"),
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
		AdminClientCert: path.Join(kubernetesCertDir, "admin", "admin.crt"),
		AdminClientKey:  path.Join(kubernetesCertDir, "admin", "admin.key"),
		KubeConfigAdmin: path.Join(kubernetesCertDir, "admin", "admin.kubeconfig"),

		// controller
		ControllerClientCert: path.Join(kubernetesCertDir, "controller", "controller.crt"),
		ControllerClientKey:  path.Join(kubernetesCertDir, "controller", "controller.key"),
		KubeConfigController: path.Join(kubernetesCertDir, "controller", "controller.kubeconfig"),

		// cloud-controller
		// ClientCloudControllerCert: path.Join(kubernetesCertDir, "cloud-controller", "client-cloud-controller.crt"),
		// ClientCloudControllerKey:  path.Join(kubernetesCertDir, "cloud-controller", "client-cloud-controller.key"),
		// KubeConfigCloudController: path.Join(kubernetesCertDir, "cloud-controller", "cloud-controller.kubeconfig"),

		// scheduler
		SchedulerClientCert: path.Join(kubernetesCertDir, "scheduler", "scheduler.crt"),
		SchedulerClientKey:  path.Join(kubernetesCertDir, "scheduler", "scheduler.key"),
		KubeConfigScheduler: path.Join(kubernetesCertDir, "scheduler", "scheduler.kubeconfig"),

		// kube-proxy
		KubeProxyClientCert: path.Join(kubernetesCertDir, "kube-proxy", "kube-proxy.crt"),
		KubeProxyClientKey:  path.Join(kubernetesCertDir, "kube-proxy", "kube-proxy.key"),

		// litekube
		LitekubeControllerClientCert: path.Join(kubernetesCertDir, "litekube-controller", "litekube-controller.crt"),
		LitekubeControllerClientKey:  path.Join(kubernetesCertDir, "litekube-controller", "litekube-controller.key"),

		ClientKubeletKey:  path.Join(kubernetesCertDir, "kubelet", "client-kubelet.key"),
		ServingKubeletKey: path.Join(kubernetesCertDir, "kubelet", "serving-kubelet.key"),

		IPSECKey:   path.Join(kubernetesCertDir, "other", "ipsec.psk"),
		PasswdFile: path.Join(kubernetesCertDir, "other", "passwd"),

		NodePasswdFile: path.Join(kubernetesCertDir, "other", "node-passwd"),
	}
}

// generate all certificates for kubernetes to start
func (na *KubernetesAuthentication) GenerateOrSkip() error {
	if err := na.generateApiserverServingCerts(); err != nil {
		return err
	}

	if err := na.generateApiserverClientKubeletCerts(); err != nil {
		return err
	}

	if err := na.generateAggregationApiserverProxyCerts(); err != nil {
		return err
	}

	if err := na.generateServiceAccountKeys(); err != nil {
		return err
	}

	if err := na.generateTokenAuthFile(); err != nil {
		return err
	}

	if global.Exists(na.KubectlPath) {
		os.Remove(na.KubectlPath)
	}
	if err := os.Symlink(na.KubeConfigAdmin, na.KubectlPath); err != nil {
		return err
	}

	return nil
}

func (na *KubernetesAuthentication) generateTokenAuthFile() error {
	if global.Exists(na.TokenAuthFile) {
		return nil
	}

	if err := os.MkdirAll(path.Join(na.KubernetesCertDir, "kubelet"), fs.FileMode(0644)); err != nil {
		return err
	}

	token, err := token.Random(16)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(na.TokenAuthFile, []byte(fmt.Sprintf("%s,kubelet-bootstrap,10001,\"system:kubelet-bootstrap\"", token)), os.FileMode(0644)); err != nil {
		return fmt.Errorf("fail to create bootstrap token")
	}

	return nil
}

func (na *KubernetesAuthentication) generateServiceAccountKeys() error {
	if _, err := certificate.GenerateRSAKeyPair(false, na.ServiceKeyPair); err != nil {
		return err
	}
	return nil
}

// generate for kube-apiserver proxy to aggregation-apiserver
func (na *KubernetesAuthentication) generateAggregationApiserverProxyCerts() error {
	regen, err := certificate.GenerateSigningCertKey(false, "litekube-request-header", na.ApiserverRequestHeaderCA, na.ApiserverRequestHeaderCAKey)
	if err != nil {
		return err
	}

	if _, err := certificate.GenerateClientCertKey(regen, na.RequestheaderAllowedNames, nil, na.ApiserverRequestHeaderCA, na.ApiserverRequestHeaderCAKey, na.ApiserverClientAuthProxyCert, na.ApiserverClientAuthProxyKey); err != nil {
		return err
	}
	return nil
}

// generate for kube-apiserver to communicate with cluster as server
func (na *KubernetesAuthentication) generateApiserverServingCerts() error {
	regenForClients, err := certificate.GenerateSigningCertKey(false, "apiserver-validate-clients", na.ApiserverValidateClientsCA, na.ApiserverValidateClientsCAKey)
	if err != nil {
		return err
	}

	regenForServer, err := certificate.GenerateSigningCertKey(false, "cluster-validate-server", na.ClusterValidateServerCA, na.ClusterValidateServerCAKey)
	if err != nil {
		return err
	}

	// kube-apiserver sign for clients
	apiseverSignFactory := getSignedClientFactory(regenForClients, na.ApiserverValidateClientsCA, na.ApiserverValidateClientsCAKey)

	// admin
	newCert, err := apiseverSignFactory("system:admin", []string{"system:masters"}, na.AdminClientCert, na.AdminClientKey)
	if err != nil {
		return err
	}
	if newCert || regenForServer {
		if err := GenKubeConfig(na.KubeConfigAdmin, na.ApiserverEndpoint, na.ClusterValidateServerCA, na.AdminClientCert, na.AdminClientKey); err != nil {
			return err
		}
	}

	// kube-controller-manager
	newCert, err = apiseverSignFactory("system:kube-controller-manager", nil, na.ControllerClientCert, na.ControllerClientKey)
	if err != nil {
		return err
	}
	if newCert || regenForServer {
		if err := GenKubeConfig(na.KubeConfigController, na.ApiserverEndpoint, na.ClusterValidateServerCA, na.ControllerClientCert, na.ControllerClientKey); err != nil {
			return err
		}
	}

	// kube-scheduler
	newCert, err = apiseverSignFactory("system:kube-scheduler", nil, na.SchedulerClientCert, na.SchedulerClientKey)
	if err != nil {
		return err
	}
	if newCert || regenForServer {
		if err := GenKubeConfig(na.KubeConfigScheduler, na.ApiserverEndpoint, na.ClusterValidateServerCA, na.SchedulerClientCert, na.SchedulerClientKey); err != nil {
			return err
		}
	}

	// kube-proxy
	if _, err = apiseverSignFactory("system:kube-proxy", nil, na.KubeProxyClientCert, na.KubeProxyClientKey); err != nil {
		return err
	}

	// litekube
	if _, err = apiseverSignFactory("system:litekube-controller", nil, na.LitekubeControllerClientCert, na.LitekubeControllerClientKey); err != nil {
		return err
	}

	// sign for  kube-apiserver server
	clusterSignFactory := getSignedServerFactory(regenForServer, &certutil.AltNames{
		DNSNames: []string{"kubernetes.default.svc", "kubernetes.default", "kubernetes", "localhost"},
		IPs:      global.RemoveRepeatIps(append(global.LocalIPs, []net.IP{net.ParseIP(na.ApiServerServiceIP), global.LocalhostIP, net.ParseIP(na.ApiserverEndpointIp)}...)),
	}, na.ClusterValidateServerCA, na.ClusterValidateServerCAKey)

	// kube-apiserver server
	if _, err = clusterSignFactory("system:kube-apiserver", nil, na.ApiserverServerCert, na.ApiserverServerKey); err != nil {
		return err
	}
	return nil
}

// generate for kube-apiserver to communicate with kubelet as client
func (na *KubernetesAuthentication) generateApiserverClientKubeletCerts() error {
	if _, err := certificate.GenerateSigningCertKey(false, "apiserver-validate-kubelet", na.ApiserverValidateKubeletServerCA, na.ApiserverValidateKubeletServerCAKey); err != nil {
		return err
	}

	regenForClients, err := certificate.GenerateSigningCertKey(false, "kubelet-validate-apiserver", na.KubeletValidateApiserverClientCA, na.KubeletValidateApiserverClientCAKey)
	if err != nil {
		return err
	}

	if _, err := certificate.GenerateClientCertKey(regenForClients, "system:kube-apiserver", nil, na.KubeletValidateApiserverClientCA, na.KubeletValidateApiserverClientCAKey, na.ApiserverClientKubeletCert, na.ApiserverClientKubeletKey); err != nil {
		return err
	}

	return nil
}

func getSignedClientFactory(regen bool, caCertPath, caKeyPath string) signedCertFactory {
	return func(commonName string, organization []string, certPath, keyPath string) (bool, error) {
		return certificate.GenerateClientCertKey(regen, commonName, organization, caCertPath, caKeyPath, certPath, keyPath)
	}
}

func getSignedServerFactory(regen bool, altNames *certutil.AltNames, caCertPath, caKeyPath string) signedCertFactory {
	return func(commonName string, organization []string, certPath, keyPath string) (bool, error) {
		return certificate.GenerateServerCertKey(regen, commonName, organization, altNames, caCertPath, caKeyPath, certPath, keyPath)
	}
}

func GenKubeConfig(filePath, url, caCert, clientCert, clientKey string) error {
	data := struct {
		URL        string
		CACert     string
		ClientCert string
		ClientKey  string
	}{
		URL:        url,
		CACert:     caCert,
		ClientCert: clientCert,
		ClientKey:  clientKey,
	}

	output, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer output.Close()

	return KubeconfigTemplate.Execute(output, &data)
}
