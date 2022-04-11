package apiserver

var UnreservedArgs []string = []string{
	"bind-address",
	"advertise-address",
	"insecure-port",
	"requestheader-extra-headers-prefix",
	"requestheader-group-headers",
	"requestheader-username-headers",
	"feature-gates",
	"allow-privileged",
	"authorization-mode",
	"anonymous-auth",
	"enable-swagger-ui",
	"enable-admission-plugins",
	"encryption-provider-config",
	"profiling",
	"service-cluster-ip-range",
	"service-node-port-range",
	"secure-port",
	"cert-dir",
	"tls-cert-file",
	"tls-private-key-file",
	"api-audiences",
	"token-auth-file",
	"enable-bootstrap-token-auth",
	"service-account-key-file",
	"service-account-issuer",
	"kubelet-certificate-authority",
	"kubelet-client-certificate",
	"kubelet-client-key",
	"client-ca-file",
	"requestheader-client-ca-file",
	"requestheader-allowed-names",
	"proxy-client-cert-file",
	"proxy-client-key-file",
	"storage-backend",
	"etcd-servers",
	"etcd-cafile",
	"etcd-certfile",
	"etcd-keyfile",
}
