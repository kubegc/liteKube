package controllermanager

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
)

// options for Litekube to start kube-controller-manager
type ControllerManagerLitekubeOptions struct {
	FeatureGates                              string `yaml:"feature-gates"`
	Kubeconfig                                string `yaml:"kubeconfig"`
	AuthorizationKubeconfig                   string `yaml:"authorization-kubeconfig"`
	AuthenticationKubeconfig                  string `yaml:"authentication-kubeconfig"`
	ServiceAccountPrivateKeyFile              string `yaml:"service-account-private-key-file"`
	AllocateNodeCidrs                         bool   `yaml:"allocate-node-cidrs"`
	ClusterCidr                               string `yaml:"cluster-cidr"`
	RootCaFile                                string `yaml:"root-ca-file"`
	Profiling                                 bool   `yaml:"profiling"`
	UseServiceAccountCredentials              bool   `yaml:"use-service-account-credentials"`
	ClusterSigningKubeApiserverClientCertFile string `yaml:"cluster-signing-kube-apiserver-client-cert-file"`
	ClusterSigningKubeApiserverClientKeyFile  string `yaml:"cluster-signing-kube-apiserver-client-key-file"`
	ClusterSigningKubeletClientCertFile       string `yaml:"cluster-signing-kubelet-client-cert-file"`
	ClusterSigningKubeletClientKeyFile        string `yaml:"cluster-signing-kubelet-client-key-file"`
	ClusterSigningKubeletServingCertFile      string `yaml:"cluster-signing-kubelet-serving-cert-file"`
	ClusterSigningKubeletServingKeyFile       string `yaml:"cluster-signing-kubelet-serving-key-file"`
	ClusterSigningLegacyUnknownCertFile       string `yaml:"cluster-signing-legacy-unknown-cert-file"`
	ClusterSigningLegacyUnknownKeyFile        string `yaml:"cluster-signing-legacy-unknown-key-file"`
}

var defaultCMLO ControllerManagerLitekubeOptions = ControllerManagerLitekubeOptions{
	FeatureGates:      "JobTrackingWithFinalizers=true",
	AllocateNodeCidrs: false,
	Profiling:         false,
}

func NewControllerManagerLitekubeOptions() *ControllerManagerLitekubeOptions {
	options := defaultCMLO
	return &options
}

func (opt *ControllerManagerLitekubeOptions) AddTips(section *help.Section) {
	section.AddTip("feature-gates", "string", "A set of key=value pairs that describe feature gates for alpha/experimental features. ", defaultCMLO.FeatureGates)
	section.AddTip("kubeconfig", "string", "Path to kubeconfig file with authorization and master location information.", defaultCMLO.Kubeconfig)
	section.AddTip("authorization-kubeconfig", "string", "kubeconfig file pointing at the 'core' kubernetes server with enough rights to create subjectaccessreviews.authorization.k8s.io. ", defaultCMLO.AuthorizationKubeconfig)
	section.AddTip("authentication-kubeconfig", "string", "kubeconfig file pointing at the 'core' kubernetes server with enough rights to create tokenreviews.authentication.k8s.io. ", defaultCMLO.AuthenticationKubeconfig)
	section.AddTip("service-account-private-key-file", "string", "Filename containing a PEM-encoded private RSA or ECDSA key used to sign service account tokens.", defaultCMLO.ServiceAccountPrivateKeyFile)
	section.AddTip("allocate-node-cidrs", "bool", "Should CIDRs for Pods be allocated and set on the cloud provider.", fmt.Sprintf("%t", defaultCMLO.AllocateNodeCidrs))
	section.AddTip("cluster-cidr", "string", "CIDR Range for Pods in cluster. Requires --allocate-node-cidrs to be true", defaultCMLO.ClusterCidr)
	section.AddTip("root-ca-file", "string", "If set, this root certificate authority will be included in service account's token secret. This must be a valid PEM-encoded CA bundle.", defaultCMLO.RootCaFile)
	section.AddTip("profiling", "bool", "Enable profiling via web interface host:port/debug/pprof/", fmt.Sprintf("%t", defaultCMLO.Profiling))
	section.AddTip("use-service-account-credentials", "bool", "If true, use individual service account credentials for each controller.", fmt.Sprintf("%t", defaultCMLO.UseServiceAccountCredentials))
	section.AddTip("cluster-signing-kube-apiserver-client-cert-file", "string", "Filename containing a PEM-encoded X509 CA certificate used to issue certificates for the kubernetes.io/kube-apiserver-client signer.", defaultCMLO.ClusterSigningKubeApiserverClientCertFile)
	section.AddTip("cluster-signing-kube-apiserver-client-key-file", "string", "Filename containing a PEM-encoded RSA or ECDSA private key used to sign certificates for the kubernetes.io/kube-apiserver-client signer. ", defaultCMLO.ClusterSigningKubeApiserverClientKeyFile)
	section.AddTip("cluster-signing-kubelet-client-cert-file", "string", "Filename containing a PEM-encoded X509 CA certificate used to issue certificates for the kubernetes.io/kube-apiserver-client-kubelet signer. ", defaultCMLO.ClusterSigningKubeletClientCertFile)
	section.AddTip("cluster-signing-kubelet-client-key-file", "string", "Filename containing a PEM-encoded RSA or ECDSA private key used to sign certificates for the kubernetes.io/kube-apiserver-client-kubelet signer.", defaultCMLO.ClusterSigningKubeletClientKeyFile)
	section.AddTip("cluster-signing-kubelet-serving-cert-file", "string", "Filename containing a PEM-encoded X509 CA certificate used to issue certificates for the kubernetes.io/kubelet-serving signer. ", defaultCMLO.ClusterSigningKubeletServingCertFile)
	section.AddTip("cluster-signing-kubelet-serving-key-file", "string", "Filename containing a PEM-encoded RSA or ECDSA private key used to sign certificates for the kubernetes.io/kubelet-serving signer.", defaultCMLO.ClusterSigningKubeletServingKeyFile)
	section.AddTip("cluster-signing-legacy-unknown-cert-file", "string", "Filename containing a PEM-encoded X509 CA certificate used to issue certificates for the kubernetes.io/legacy-unknown signer.", defaultCMLO.ClusterSigningLegacyUnknownCertFile)
	section.AddTip("cluster-signing-legacy-unknown-key-file", "string", "Filename containing a PEM-encoded RSA or ECDSA private key used to sign certificates for the kubernetes.io/legacy-unknown signer. ", defaultCMLO.ClusterSigningLegacyUnknownKeyFile)

}
