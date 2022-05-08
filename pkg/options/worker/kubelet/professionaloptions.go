package kubelet

import (
	"github.com/litekube/LiteKube/pkg/help"
)

// Empirically assigned parameters are not recommended
type KubeletProfessionalOptions struct {
	//NodeIp                   string `yaml:"node-ip"`
	Kubeconfig               string `yaml:"kubeconfig"`
	Config                   string `yaml:"config"`
	CgroupDriver             string `yaml:"cgroup-driver"`
	RuntimeCgroups           string `yaml:"runtime-cgroups"`
	HostnameOverride         string `yaml:"hostname-override"`
	ContainerRuntime         string `yaml:"container-runtime"`
	BootstrapKubeconfig      string `yaml:"bootstrap-kubeconfig"`
	ContainerRuntimeEndpoint string `yaml:"container-runtime-endpoint"`
}

func NewKubeletProfessionalOptions() *KubeletProfessionalOptions {
	options := DefaultKPO
	return &options
}

var DefaultKPO KubeletProfessionalOptions = KubeletProfessionalOptions{
	//NodeIp:                   "127.0.0.1",
	CgroupDriver:             "systemd",
	ContainerRuntime:         "remote",
	ContainerRuntimeEndpoint: "unix:///run/containerd/containerd.sock",
	RuntimeCgroups:           "/system.slice/containerd.service",
}

func (opt *KubeletProfessionalOptions) AddTips(section *help.Section) {
	//section.AddTip("node-ip", "string", "IP address (or comma-separated dual-stack IP addresses) of the node.", DefaultKPO.NodeIp)
	section.AddTip("kubeconfig", "string", "Path to a kubeconfig file, specifying how to connect to the API server.", DefaultKPO.Kubeconfig)
	section.AddTip("config", "string", "The Kubelet will load its initial configuration from this file. ", DefaultKPO.Config)
	section.AddTip("cgroup-driver", "string", "Driver that the kubelet uses to manipulate cgroups on the host. ", DefaultKPO.CgroupDriver)
	section.AddTip("hostname-override", "string", "	If non-empty, will use this string as identification instead of the actual hostname.", DefaultKPO.HostnameOverride)
	section.AddTip("container-runtime", "string", "The container runtime to use. Possible values: docker, remote.", DefaultKPO.ContainerRuntime)
	section.AddTip("runtime-cgroups", "string", "Optional absolute name of cgroups to create and run the runtime in.", DefaultKPO.RuntimeCgroups)
	section.AddTip("bootstrap-kubeconfig", "string", "Path to a kubeconfig file that will be used to get client certificate for kubelet. If the file specified by --kubeconfig does not exist, the bootstrap kubeconfig is used to request a client certificate from the API server. ", DefaultKPO.BootstrapKubeconfig)
	section.AddTip("container-runtime-endpoint", "string", "[Experimental] The endpoint of remote runtime service. Currently unix socket endpoint is supported on Linux, while npipe and tcp endpoints are supported on windows. Examples: unix:///var/run/dockershim.sock, npipe:////./pipe/dockershim.", DefaultKPO.ContainerRuntimeEndpoint)
}
