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

	BindAddress      string `yaml:"bind-address"`
	AdvertiseAddress string `yaml:"advertise-address"`
	InsecurePort     uint16 `yaml:"insecure-port"`
	FeatureGates     string `yaml:"feature-gates"`
}

// server security
type ServerCertOptions struct {
	CertDir                      string `yaml:"cert-dir"`
	TlsCertFile                  string `yaml:"tls-cert-file"`
	TlsPrivateKeyFile            string `yaml:"tls-private-key-file"`
	ApiAudiences                 string `yaml:"api-audiences"`
	TokenAuthFile                string `yaml:"token-auth-file"`
	EnableBootstrapTokenAuth     bool   `yaml:"enable-bootstrap-token-auth"`
	ServiceAccountSigningKeyFile string `yaml:"service-account-signing-key-file"`
	ServiceAccountKeyFile        string `yaml:"service-account-key-file"`
	ServiceAccountIssuer         string `yaml:"service-account-issuer"`
	ClientCAFile                 string `yaml:"client-ca-file"`

	// for access-proxy to kube-apiserver
	RequestheaderExtraHeadersPrefix string `yaml:"requestheader-extra-headers-prefix"`
	RequestheaderGroupHeaders       string `yaml:"requestheader-group-headers"`
	RequestheaderUsernameHeaders    string `yaml:"requestheader-username-headers"`
	RequestheaderClientCAFile       string `yaml:"requestheader-client-ca-file"`
	RequestheaderAllowedNames       string `yaml:"requestheader-allowed-names"`
	ProxyClientCertFile             string `yaml:"proxy-client-cert-file"`
	ProxyClientKeyFile              string `yaml:"proxy-client-key-file"`
	EnableAggregatorRouting         bool   `yaml:"enable-aggregator-routing"`
	//ServiceAccountSigningKeyFile string `yaml:""`
}

// security information for apiserver-kubelet-client-config
type KubeletClientCertOptions struct {
	KubeletCertificateAuthority string `yaml:"kubelet-certificate-authority"`
	KubeletClientCertificate    string `yaml:"kubelet-client-certificate"`
	KubeletClientKey            string `yaml:"kubelet-client-key"`
}

// etcd options
type ECTDOptions struct {
	StorageBackend string `yaml:"storage-backend"`
	EtcdServers    string `yaml:"etcd-servers"`
	EtcdCafile     string `yaml:"etcd-cafile"`
	EtcdCertfile   string `yaml:"etcd-certfile"`
	EtcdKeyfile    string `yaml:"etcd-keyfile"`
}

var DefaultEO ECTDOptions = ECTDOptions{
	StorageBackend: "etcd3",
	EtcdServers:    "https://127.0.0.1:2379",
}
var DefaultKCCO KubeletClientCertOptions = KubeletClientCertOptions{}
var DefaultSCO ServerCertOptions = ServerCertOptions{
	ApiAudiences:             "unknown",
	EnableBootstrapTokenAuth: true,
	ServiceAccountIssuer:     "litekube",

	RequestheaderExtraHeadersPrefix: "X-Remote-Extra-",
	RequestheaderGroupHeaders:       "X-Remote-Group",
	RequestheaderUsernameHeaders:    "X-Remote-User",
	RequestheaderAllowedNames:       "system:auth-proxy",
	EnableAggregatorRouting:         true,
}

var DefaultAPO ApiserverProfessionalOptions = ApiserverProfessionalOptions{
	ECTDOptions:              *NewECTDOptions(),
	ServerCertOptions:        *NewServerCertOptions(),
	KubeletClientCertOptions: *NewKubeletClientCertOptions(),
	BindAddress:              "0.0.0.0",
	InsecurePort:             0,
	FeatureGates:             "JobTrackingWithFinalizers=true",
}

func NewKubeletClientCertOptions() *KubeletClientCertOptions {
	options := DefaultKCCO
	return &options
}

func NewServerCertOptions() *ServerCertOptions {
	options := DefaultSCO
	return &options
}

func NewECTDOptions() *ECTDOptions {
	options := DefaultEO
	return &options
}

func NewApiserverProfessionalOptions() *ApiserverProfessionalOptions {
	options := DefaultAPO
	return &options
}

func (opt *ServerCertOptions) AddTips(section *help.Section) {
	section.AddTip("cert-dir", "string", "The directory where the TLS certs are located. If --tls-cert-file and --tls-private-key-file are provided, this flag will be ignored.", DefaultSCO.CertDir)
	section.AddTip("tls-cert-file", "string", "File containing the default x509 Certificate for HTTPS.", DefaultSCO.TlsCertFile)
	section.AddTip("tls-private-key-file", "string", "File containing the default x509 private key matching --tls-cert-file.", DefaultSCO.TlsPrivateKeyFile)
	section.AddTip("api-audiences", "string", "Identifiers of the API.", DefaultSCO.ApiAudiences)
	section.AddTip("token-auth-file", "string", "If set, the file that will be used to secure the secure port of the API server via token authentication.", DefaultSCO.TokenAuthFile)
	section.AddTip("enable-bootstrap-token-auth", "bool", "Enable to allow secrets of type 'bootstrap.kubernetes.io/token' in the 'kube-system' namespace to be used for TLS bootstrapping authentication.", fmt.Sprintf("%t", DefaultSCO.EnableBootstrapTokenAuth))
	section.AddTip("service-account-signing-key-file", "string", "File containing PEM-encoded x509 RSA or ECDSA private or public keys, used to verify ServiceAccount tokens.", DefaultSCO.ServiceAccountSigningKeyFile)
	section.AddTip("service-account-key-file", "string", "File containing PEM-encoded x509 RSA or ECDSA private or public keys, used to verify ServiceAccount tokens.", DefaultSCO.ServiceAccountKeyFile)
	section.AddTip("service-account-issuer", "string", "Identifier of the service account token issuer.", DefaultSCO.ServiceAccountIssuer)
	section.AddTip("requestheader-extra-headers-prefix", "string", "List of request header prefixes to inspect. X-Remote-Extra- is suggested.", DefaultSCO.RequestheaderExtraHeadersPrefix)
	section.AddTip("requestheader-group-headers", "string", "List of request headers to inspect for groups. X-Remote-Group is suggested.", DefaultSCO.RequestheaderGroupHeaders)
	section.AddTip("requestheader-username-headers", "string", "List of request headers to inspect for usernames. X-Remote-User is common.", DefaultSCO.RequestheaderUsernameHeaders)
	section.AddTip("requestheader-client-ca-file", "string", "Root certificate bundle to use to verify client certificates on incoming requests before trusting usernames in headers specified by --requestheader-username-headers.", DefaultSCO.RequestheaderClientCAFile)
	section.AddTip("requestheader-allowed-names", "string", "List of client certificate common names to allow to provide usernames in headers specified by --requestheader-username-headers.", DefaultSCO.RequestheaderAllowedNames)
	section.AddTip("proxy-client-cert-file", "string", "Client certificate used to prove the identity of the aggregator or kube-apiserver when it must call out during a request. ", DefaultSCO.ProxyClientCertFile)
	section.AddTip("proxy-client-key-file", "string", "Private key for the client certificate used to prove the identity of the aggregator or kube-apiserver when it must call out during a request. ", DefaultSCO.ProxyClientKeyFile)
	section.AddTip("enable-aggregator-routing", "bool", "set true is suggested.", fmt.Sprintf("%t", DefaultSCO.EnableAggregatorRouting))
	section.AddTip("client-ca-file", "string", "If set, any request presenting a client certificate signed by one of the authorities in the client-ca-file is authenticated with an identity corresponding to the CommonName of the client certificate.", DefaultSCO.ClientCAFile)
}

func (opt *ECTDOptions) AddTips(section *help.Section) {
	section.AddTip("storage-backend", "string", "The storage backend for persistence. (default).", DefaultEO.StorageBackend)
	section.AddTip("etcd-servers", "string", "List of etcd servers to connect with (scheme://ip:port), comma separated.", DefaultEO.EtcdServers)
	section.AddTip("etcd-cafile", "string", "SSL Certificate Authority file used to secure etcd communication.", DefaultEO.EtcdCafile)
	section.AddTip("etcd-certfile", "string", "SSL certification file used to secure etcd communication.", DefaultEO.EtcdCertfile)
	section.AddTip("etcd-keyfile", "string", "SSL key file used to secure etcd communication.", DefaultEO.EtcdKeyfile)
}

func (opt *KubeletClientCertOptions) AddTips(section *help.Section) {
	section.AddTip("kubelet-certificate-authority", "string", "Path to a cert file for the certificate authority.", DefaultKCCO.KubeletCertificateAuthority)
	section.AddTip("kubelet-client-certificate", "string", "Path to a client cert file for TLS.", DefaultKCCO.KubeletClientCertificate)
	section.AddTip("kubelet-client-key", "string", "Path to a client key file for TLS.", DefaultKCCO.KubeletClientKey)
}

func (opt *ApiserverProfessionalOptions) AddTips(section *help.Section) {
	section.AddTip("bind-address", "string", "The IP address on which to listen for the --secure-port port.", DefaultAPO.BindAddress)
	section.AddTip("advertise-address", "string", "The IP address on which to advertise the apiserver to members of the cluster.", DefaultAPO.AdvertiseAddress)
	section.AddTip("insecure-port", "uint16", "Disabled, HTTP Apiserver port", fmt.Sprintf("%d", DefaultAPO.InsecurePort))
	section.AddTip("feature-gates", "string", "A set of key=value pairs that describe feature gates for alpha/experimental features.", DefaultAPO.FeatureGates)

	opt.ECTDOptions.AddTips(section)
	opt.KubeletClientCertOptions.AddTips(section)
	opt.ServerCertOptions.AddTips(section)
}
