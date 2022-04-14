package apiserver

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
)

// options for Litekube to start kube-apiserver
type ApiserverLitekubeOptions struct {
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

var defaultALO ApiserverLitekubeOptions = ApiserverLitekubeOptions{}

func NewApiserverLitekubeOptions() *ApiserverLitekubeOptions {
	options := defaultALO
	return &options
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
}
