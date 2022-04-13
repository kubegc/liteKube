package apiserver

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
)

// options for Litekube to start kube-apiserver
type ApiserverLitekubeOptions struct {
	ECTDOptions              `yaml:",inline"`
	ServerCertOptions        `yaml:",inline"`
	KubeletClientCertOptions `yaml:",inline"`

	AllowPrivileged          bool   `yaml:"allow-privileged"`
	AuthorizationMode        string `yaml:"authorization-mode"`
	AnonymousAuth            bool   `yaml:"anonymous-auth"`
	EnableSwaggerUI          bool   `yaml:"enable-swagger-ui"`
	EnableAdmissionPlugins   string `yaml:"enable-admission-plugins"`
	EncryptionProviderConfig string `yaml:"encryption-provider-config"`
	Profiling                bool   `yaml:"profiling"`
	ServiceClusterIpRange    string `yaml:"service-cluster-ip-range"`
	ServiceNodePortRange     string `yaml:"service-node-port-range"`
	SecurePort               int16  `yaml:"secure-port"`
}

// server security
type ServerCertOptions struct {
	CertDir                  string `yaml:"cert-dir"`
	TlsCertFile              string `yaml:"tls-cert-file"`
	TlsPrivateKeyFile        string `yaml:"tls-private-key-file"`
	ApiAudiencesr            string `yaml:"api-audiences"`
	TokenAuthFile            string `yaml:"token-auth-file"`
	EnableBootstrapTokenAuth bool   `yaml:"enable-bootstrap-token-auth"`
	ServiceAccountKeyFile    string `yaml:"service-account-key-file"`
	ServiceAccountIssuer     string `yaml:"service-account-issuer"`
	//ServiceAccountSigningKeyFile string `yaml:""`
}

// security information for apiserver-kubelet-client-config
type KubeletClientCertOptions struct {
	KubeletCertificateAuthority string `yaml:"kubelet-certificate-authority"`
	KubeletClientCertificate    string `yaml:"kubelet-client-certificate"`
	KubeletClientKey            string `yaml:"kubelet-client-key"`
	ClientCAFile                string `yaml:"client-ca-file"`
	RequestheaderClientCAFile   string `yaml:"requestheader-client-ca-file"`
	RequestheaderAllowedNames   string `yaml:"requestheader-allowed-names"`
	ProxyClientCertFile         string `yaml:"proxy-client-cert-file"`
	ProxyClientKeyFile          string `yaml:"proxy-client-key-file"`
}

// etcd options
type ECTDOptions struct {
	StorageBackend string `yaml:"storage-backend"`
	EtcdServers    string `yaml:"etcd-servers"`
	EtcdCafile     string `yaml:"etcd-cafile"`
	EtcdCertfile   string `yaml:"etcd-certfile"`
	EtcdKeyfile    string `yaml:"etcd-keyfile"`
}

var defaultEO ECTDOptions = ECTDOptions{
	StorageBackend: "etcd3",
	EtcdServers:    "https://127.0.0.1:2379",
}
var defaultKCCO KubeletClientCertOptions = KubeletClientCertOptions{}
var defaultSCO ServerCertOptions = ServerCertOptions{}
var defaultALO ApiserverLitekubeOptions = ApiserverLitekubeOptions{
	ECTDOptions:              ECTDOptions{},
	ServerCertOptions:        ServerCertOptions{},
	KubeletClientCertOptions: KubeletClientCertOptions{},
}

func NewKubeletClientCertOptions() *KubeletClientCertOptions {
	options := defaultKCCO
	return &options
}

func NewServerCertOptions() *ServerCertOptions {
	options := defaultSCO
	return &options
}

func NewECTDOptions() *ECTDOptions {
	options := defaultEO
	return &options
}

func NewApiserverLitekubeOptions() *ApiserverLitekubeOptions {
	options := defaultALO
	return &options
}

func (opt *ServerCertOptions) AddTips(section *help.Section) {
	section.AddTip("cert-dir", "string", "The directory where the TLS certs are located. If --tls-cert-file and --tls-private-key-file are provided, this flag will be ignored.", defaultSCO.CertDir)
	section.AddTip("tls-cert-file", "string", "File containing the default x509 Certificate for HTTPS.", defaultSCO.TlsCertFile)
	section.AddTip("tls-private-key-file", "string", "File containing the default x509 private key matching --tls-cert-file.", defaultSCO.TlsPrivateKeyFile)
	section.AddTip("api-audiences", "string", "Identifiers of the API.", defaultSCO.ApiAudiencesr)
	section.AddTip("token-auth-file", "string", "If set, the file that will be used to secure the secure port of the API server via token authentication.", defaultSCO.TokenAuthFile)
	section.AddTip("enable-bootstrap-token-auth", "bool", "Enable to allow secrets of type 'bootstrap.kubernetes.io/token' in the 'kube-system' namespace to be used for TLS bootstrapping authentication.", fmt.Sprintf("%t", defaultSCO.EnableBootstrapTokenAuth))
	section.AddTip("service-account-key-file", "string", "File containing PEM-encoded x509 RSA or ECDSA private or public keys, used to verify ServiceAccount tokens.", defaultSCO.ServiceAccountKeyFile)
	section.AddTip("service-account-issuer", "string", "Identifier of the service account token issuer.", defaultSCO.ServiceAccountIssuer)
}

func (opt *ECTDOptions) AddTips(section *help.Section) {
	section.AddTip("storage-backend", "string", "The storage backend for persistence. (default).", defaultEO.StorageBackend)
	section.AddTip("etcd-servers", "string", "List of etcd servers to connect with (scheme://ip:port), comma separated.", defaultEO.EtcdServers)
	section.AddTip("etcd-cafile", "string", "SSL Certificate Authority file used to secure etcd communication.", defaultEO.EtcdCafile)
	section.AddTip("etcd-certfile", "string", "SSL certification file used to secure etcd communication.", defaultEO.EtcdCertfile)
	section.AddTip("etcd-keyfile", "string", "SSL key file used to secure etcd communication.", defaultEO.EtcdKeyfile)
}

func (opt *KubeletClientCertOptions) AddTips(section *help.Section) {
	section.AddTip("kubelet-certificate-authority", "string", "Path to a cert file for the certificate authority.", defaultKCCO.KubeletCertificateAuthority)
	section.AddTip("kubelet-client-certificate", "string", "Path to a client cert file for TLS.", defaultKCCO.KubeletClientCertificate)
	section.AddTip("kubelet-client-key", "string", "Path to a client key file for TLS.", defaultKCCO.KubeletClientKey)
	section.AddTip("client-ca-file", "string", "If set, any request presenting a client certificate signed by one of the authorities in the client-ca-file is authenticated with an identity corresponding to the CommonName of the client certificate.", defaultKCCO.ClientCAFile)
	section.AddTip("requestheader-client-ca-file", "string", "Root certificate bundle to use to verify client certificates on incoming requests before trusting usernames in headers specified by --requestheader-username-headers.", defaultKCCO.RequestheaderClientCAFile)
	section.AddTip("requestheader-allowed-names", "string", "List of client certificate common names to allow to provide usernames in headers specified by --requestheader-username-headers.", defaultKCCO.RequestheaderAllowedNames)
	section.AddTip("proxy-client-cert-file", "string", "Client certificate used to prove the identity of the aggregator or kube-apiserver when it must call out during a request. ", defaultKCCO.ProxyClientCertFile)
	section.AddTip("proxy-client-key-file", "string", "Private key for the client certificate used to prove the identity of the aggregator or kube-apiserver when it must call out during a request. ", defaultKCCO.ProxyClientKeyFile)
}

func (opt *ApiserverLitekubeOptions) AddTips(section *help.Section) {
	section.AddTip("allow-privileged", "bool", "If true, allow privileged containers. ", fmt.Sprintf("%t", defaultALO.AllowPrivileged))
	section.AddTip("authorization-mode", "string", "File with authorization policy in json line by line format", defaultALO.AuthorizationMode)
	section.AddTip("anonymous-auth", "bool", "Enables anonymous requests to the secure port of the API server.", fmt.Sprintf("%t", defaultALO.AnonymousAuth))
	section.AddTip("enable-swagger-ui", "bool", "Disabled, enable swagger ui.", fmt.Sprintf("%t", defaultALO.EnableSwaggerUI))
	section.AddTip("enable-admission-plugins", "string", "admission plugins that should be enabled in addition to default enabled ones", defaultALO.EnableAdmissionPlugins)
	section.AddTip("encryption-provider-config", "string", "The file containing configuration for encryption providers to be used for storing secrets in etcd", defaultALO.EncryptionProviderConfig)
	section.AddTip("profiling", "bool", "Enable profiling via web interface host:port/debug/pprof/", fmt.Sprintf("%t", defaultALO.Profiling))
	section.AddTip("service-cluster-ip-range", "string", "A CIDR notation IP range from which to assign service cluster IPs. This must not overlap with any IP ranges assigned to nodes or pods. Max of two dual-stack CIDRs is allowed.", defaultALO.ServiceClusterIpRange)
	section.AddTip("service-node-port-range", "string", "A port range to reserve for services with NodePort visibility. Example: '30000-32767'. Inclusive at both ends of the range.", defaultALO.ServiceNodePortRange)
	section.AddTip("secure-port", "int16", "The port on which to serve HTTPS with authentication and authorization. It cannot be switched off with 0.", fmt.Sprintf("%d", defaultALO.SecurePort))

	opt.ECTDOptions.AddTips(section)
	opt.KubeletClientCertOptions.AddTips(section)
	opt.ServerCertOptions.AddTips(section)
}
