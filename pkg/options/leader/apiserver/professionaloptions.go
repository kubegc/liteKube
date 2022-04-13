package apiserver

import (
	"fmt"

	"github.com/litekube/LiteKube/pkg/help"
)

// Empirically assigned parameters are not recommended
type ApiserverProfessionalOptions struct {
	BindAddress                     string `yaml:"bind-address"`
	AdvertiseAddress                string `yaml:"advertise-address"`
	InsecurePort                    int16  `yaml:"insecure-port"`
	RequestheaderExtraHeadersPrefix string `yaml:"requestheader-extra-headers-prefix"`
	RequestheaderGroupHeaders       string `yaml:"requestheader-group-headers"`
	RequestheaderUsernameHeaders    string `yaml:"requestheader-username-headers"`
	FeatureGates                    string `yaml:"feature-gates"`
}

var defaultAPO ApiserverProfessionalOptions = ApiserverProfessionalOptions{
	BindAddress:                     "0.0.0.0",
	InsecurePort:                    0,
	RequestheaderExtraHeadersPrefix: "X-Remote-Extra-",
	RequestheaderGroupHeaders:       "X-Remote-Group",
	RequestheaderUsernameHeaders:    "X-Remote-User",
	FeatureGates:                    "JobTrackingWithFinalizers=true",
}

func NewApiserverProfessionalOptions() *ApiserverProfessionalOptions {
	options := defaultAPO
	return &options
}

func (opt *ApiserverProfessionalOptions) AddTips(section *help.Section) {
	section.AddTip("bind-address", "string", "The IP address on which to listen for the --secure-port port.", defaultAPO.BindAddress)
	section.AddTip("advertise-address", "string", "The IP address on which to advertise the apiserver to members of the cluster.", defaultAPO.AdvertiseAddress)
	section.AddTip("insecure-port", "int16", "Disabled, HTTP Apiserver port", fmt.Sprintf("%d", defaultAPO.InsecurePort))
	section.AddTip("requestheader-extra-headers-prefix", "string", "List of request header prefixes to inspect. X-Remote-Extra- is suggested.", defaultAPO.RequestheaderExtraHeadersPrefix)
	section.AddTip("requestheader-group-headers", "string", "List of request headers to inspect for groups. X-Remote-Group is suggested.", defaultAPO.RequestheaderGroupHeaders)
	section.AddTip("requestheader-username-headers", "string", "List of request headers to inspect for usernames. X-Remote-User is common.", defaultAPO.RequestheaderUsernameHeaders)
	section.AddTip("feature-gates", "string", "A set of key=value pairs that describe feature gates for alpha/experimental features.", defaultAPO.FeatureGates)
}
