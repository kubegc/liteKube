package apiserver

import "github.com/litekube/LiteKube/pkg/help"

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

func NewApiserverProfessionalOptions() *ApiserverProfessionalOptions {
	return &ApiserverProfessionalOptions{}
}

func (opt *ApiserverProfessionalOptions) AddTips(section *help.Section) {
	section.AddTip("bind-address", "string", "The IP address on which to listen for the --secure-port port.", "0.0.0.0")
	section.AddTip("advertise-address", "string", "The IP address on which to advertise the apiserver to members of the cluster.", "<auto by register>")
	section.AddTip("insecure-port", "int16", "Disabled, HTTP Apiserver port", "0")
	section.AddTip("requestheader-extra-headers-prefix", "string", "List of request header prefixes to inspect. X-Remote-Extra- is suggested.", "X-Remote-Extra-")
	section.AddTip("requestheader-group-headers", "string", "List of request headers to inspect for groups. X-Remote-Group is suggested.", "X-Remote-Group")
	section.AddTip("requestheader-username-headers", "string", "List of request headers to inspect for usernames. X-Remote-User is common.", "X-Remote-User")
	section.AddTip("feature-gates", "string", "A set of key=value pairs that describe feature gates for alpha/experimental features.", "JobTrackingWithFinalizers=true")
}

// c
