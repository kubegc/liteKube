package controllermanager

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
)

// Empirically assigned parameters are not recommended
type ControllerManagerProfessionalOptions struct {
	BindAddress                               string `yaml:"bind-address"`
	SecurePort                                uint16 `yaml:"secure-port"`
	LeaderElect                               bool   `yaml:"leader-elect"`
	ConfigureCloudRoutes                      bool   `yaml:"configure-cloud-routes"`
	Controllers                               string `yaml:"controllers"`
	FeatureGates                              string `yaml:"feature-gates"`
	Kubeconfig                                string `yaml:"kubeconfig"`
	AuthorizationKubeconfig                   string `yaml:"authorization-kubeconfig"`
	AuthenticationKubeconfig                  string `yaml:"authentication-kubeconfig"`
	ServiceAccountPrivateKeyFile              string `yaml:"service-account-private-key-file"`
	RootCaFile                                string `yaml:"root-ca-file"`
	ClusterSigningKubeApiserverClientCertFile string `yaml:"cluster-signing-kube-apiserver-client-cert-file"`
	ClusterSigningKubeApiserverClientKeyFile  string `yaml:"cluster-signing-kube-apiserver-client-key-file"`
	ClusterSigningKubeletClientCertFile       string `yaml:"cluster-signing-kubelet-client-cert-file"`
	ClusterSigningKubeletClientKeyFile        string `yaml:"cluster-signing-kubelet-client-key-file"`
	ClusterSigningKubeletServingCertFile      string `yaml:"cluster-signing-kubelet-serving-cert-file"`
	ClusterSigningKubeletServingKeyFile       string `yaml:"cluster-signing-kubelet-serving-key-file"`
	ClusterSigningLegacyUnknownCertFile       string `yaml:"cluster-signing-legacy-unknown-cert-file"`
	ClusterSigningLegacyUnknownKeyFile        string `yaml:"cluster-signing-legacy-unknown-key-file"`
}

func NewControllerManagerProfessionalOptions() *ControllerManagerProfessionalOptions {
	options := defaultCMPO
	return &options
}

var defaultCMPO ControllerManagerProfessionalOptions = ControllerManagerProfessionalOptions{
	BindAddress:          "0.0.0.0",
	SecurePort:           10257,
	LeaderElect:          false,
	ConfigureCloudRoutes: false,
	Controllers:          "*,-service,-route,-cloud-node-lifecycle",
	FeatureGates:         "JobTrackingWithFinalizers=true",
}

func (opt *ControllerManagerProfessionalOptions) AddTips(section *help.Section) {
	section.AddTip("root-ca-file", "string", "If set, this root certificate authority will be included in service account's token secret. This must be a valid PEM-encoded CA bundle.", defaultCMPO.RootCaFile)
	section.AddTip("feature-gates", "string", "A set of key=value pairs that describe feature gates for alpha/experimental features. ", defaultCMPO.FeatureGates)
	section.AddTip("kubeconfig", "string", "Path to kubeconfig file with authorization and master location information.", defaultCMPO.Kubeconfig)
	section.AddTip("authorization-kubeconfig", "string", "kubeconfig file pointing at the 'core' kubernetes server with enough rights to create subjectaccessreviews.authorization.k8s.io. ", defaultCMPO.AuthorizationKubeconfig)
	section.AddTip("authentication-kubeconfig", "string", "kubeconfig file pointing at the 'core' kubernetes server with enough rights to create tokenreviews.authentication.k8s.io. ", defaultCMPO.AuthenticationKubeconfig)
	section.AddTip("service-account-private-key-file", "string", "Filename containing a PEM-encoded private RSA or ECDSA key used to sign service account tokens.", defaultCMPO.ServiceAccountPrivateKeyFile)
	section.AddTip("bind-address", "string", "The IP address on which to listen for the --secure-port port. ", defaultCMPO.BindAddress)
	section.AddTip("secure-port", "uint16", "The port on which to serve HTTPS with authentication and authorization. If 0, don't serve HTTPS at all.", fmt.Sprintf("%d", defaultCMPO.SecurePort))
	section.AddTip("leader-elect", "bool", "Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability.", fmt.Sprintf("%t", defaultCMPO.LeaderElect))
	section.AddTip("configure-cloud-routes", "bool", "Should CIDRs allocated by allocate-node-cidrs be configured on the cloud provider.", fmt.Sprintf("%t", defaultCMPO.ConfigureCloudRoutes))
	section.AddTip("controllers", "string", "A list of controllers to enable. ", defaultCMPO.Controllers)
	section.AddTip("cluster-signing-kube-apiserver-client-cert-file", "string", "Filename containing a PEM-encoded X509 CA certificate used to issue certificates for the kubernetes.io/kube-apiserver-client signer.", defaultCMPO.ClusterSigningKubeApiserverClientCertFile)
	section.AddTip("cluster-signing-kube-apiserver-client-key-file", "string", "Filename containing a PEM-encoded RSA or ECDSA private key used to sign certificates for the kubernetes.io/kube-apiserver-client signer. ", defaultCMPO.ClusterSigningKubeApiserverClientKeyFile)
	section.AddTip("cluster-signing-kubelet-client-cert-file", "string", "Filename containing a PEM-encoded X509 CA certificate used to issue certificates for the kubernetes.io/kube-apiserver-client-kubelet signer. ", defaultCMPO.ClusterSigningKubeletClientCertFile)
	section.AddTip("cluster-signing-kubelet-client-key-file", "string", "Filename containing a PEM-encoded RSA or ECDSA private key used to sign certificates for the kubernetes.io/kube-apiserver-client-kubelet signer.", defaultCMPO.ClusterSigningKubeletClientKeyFile)
	section.AddTip("cluster-signing-kubelet-serving-cert-file", "string", "Filename containing a PEM-encoded X509 CA certificate used to issue certificates for the kubernetes.io/kubelet-serving signer. ", defaultCMPO.ClusterSigningKubeletServingCertFile)
	section.AddTip("cluster-signing-kubelet-serving-key-file", "string", "Filename containing a PEM-encoded RSA or ECDSA private key used to sign certificates for the kubernetes.io/kubelet-serving signer.", defaultCMPO.ClusterSigningKubeletServingKeyFile)
	section.AddTip("cluster-signing-legacy-unknown-cert-file", "string", "Filename containing a PEM-encoded X509 CA certificate used to issue certificates for the kubernetes.io/legacy-unknown signer.", defaultCMPO.ClusterSigningLegacyUnknownCertFile)
	section.AddTip("cluster-signing-legacy-unknown-key-file", "string", "Filename containing a PEM-encoded RSA or ECDSA private key used to sign certificates for the kubernetes.io/legacy-unknown signer. ", defaultCMPO.ClusterSigningLegacyUnknownKeyFile)
}
