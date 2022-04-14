package apiserver

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
)

// Empirically assigned parameters are not recommended
type ApiserverProfessionalOptions struct {
	ECTDOptions              `yaml:",inline"`
	ServerCertOptions        `yaml:",inline"`
	KubeletClientCertOptions `yaml:",inline"`

	BindAddress                     string `yaml:"bind-address"`
	AdvertiseAddress                string `yaml:"advertise-address"`
	InsecurePort                    uint16 `yaml:"insecure-port"`
	RequestheaderExtraHeadersPrefix string `yaml:"requestheader-extra-headers-prefix"`
	RequestheaderGroupHeaders       string `yaml:"requestheader-group-headers"`
	RequestheaderUsernameHeaders    string `yaml:"requestheader-username-headers"`
	FeatureGates                    string `yaml:"feature-gates"`
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

var defaultAPO ApiserverProfessionalOptions = ApiserverProfessionalOptions{
	ECTDOptions:                     *NewECTDOptions(),
	ServerCertOptions:               *NewServerCertOptions(),
	KubeletClientCertOptions:        *NewKubeletClientCertOptions(),
	BindAddress:                     "0.0.0.0",
	InsecurePort:                    0,
	RequestheaderExtraHeadersPrefix: "X-Remote-Extra-",
	RequestheaderGroupHeaders:       "X-Remote-Group",
	RequestheaderUsernameHeaders:    "X-Remote-User",
	FeatureGates:                    "JobTrackingWithFinalizers=true",
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

func NewApiserverProfessionalOptions() *ApiserverProfessionalOptions {
	options := defaultAPO
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

func (opt *ApiserverProfessionalOptions) AddTips(section *help.Section) {
	section.AddTip("bind-address", "string", "The IP address on which to listen for the --secure-port port.", defaultAPO.BindAddress)
	section.AddTip("advertise-address", "string", "The IP address on which to advertise the apiserver to members of the cluster.", defaultAPO.AdvertiseAddress)
	section.AddTip("insecure-port", "uint16", "Disabled, HTTP Apiserver port", fmt.Sprintf("%d", defaultAPO.InsecurePort))
	section.AddTip("requestheader-extra-headers-prefix", "string", "List of request header prefixes to inspect. X-Remote-Extra- is suggested.", defaultAPO.RequestheaderExtraHeadersPrefix)
	section.AddTip("requestheader-group-headers", "string", "List of request headers to inspect for groups. X-Remote-Group is suggested.", defaultAPO.RequestheaderGroupHeaders)
	section.AddTip("requestheader-username-headers", "string", "List of request headers to inspect for usernames. X-Remote-User is common.", defaultAPO.RequestheaderUsernameHeaders)
	section.AddTip("feature-gates", "string", "A set of key=value pairs that describe feature gates for alpha/experimental features.", defaultAPO.FeatureGates)

	opt.ECTDOptions.AddTips(section)
	opt.KubeletClientCertOptions.AddTips(section)
	opt.ServerCertOptions.AddTips(section)
}
