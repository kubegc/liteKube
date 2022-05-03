package likutemplate

import "html/template"

var (
	Kubelet_config_template = template.Must(template.New("kubelet_config").Parse(`kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
authentication:
  x509:
    clientCAFile: "{{.CaPath}}"
  webhook:
    enabled: true
    cacheTTL: 2m0s
  anonymous:
    enabled: false
authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: 5m0s
    cacheUnauthorizedTTL: 30s
address: "0.0.0.0"
port: 10250
readOnlyPort: 10255
cgroupDriver: systemd
hairpinMode: promiscuous-bridge
serializeImagePulls: false
clusterDomain: cluster.local.
clusterDNS:
- "{{.CluserDNS}}"
`))

	Kubelet_bootstrap_kubeconfig_template = template.Must(template.New("kubelet_bootstrap_kubeconfig").Parse(`apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: {{.CACert}}
    server: {{.URL}}
  name: litekube
contexts:
- context:
    cluster: litekube
    user: kubelet-bootstrap
  name: default
current-context: default
kind: Config
preferences: {}
users:
- name: kubelet-bootstrap
  user:
    token: {{.Token}}
`))

	Kubeproxy_bootstrap_kubeconfig_template = template.Must(template.New("kubeproxy_kubeconfig").Parse(`apiVersion: v1
clusters:
- cluster:
    server: {{.URL}}
    certificate-authority-data: {{.CACert}}
  name: litekube
contexts:
- context:
    cluster: litekube
    user: kube-proxy
  name: default
current-context: default
kind: Config
preferences: {}
users:
- name: kube-proxy
  user:
    client-certificate-data: {{.ClientCert}}
    client-key-data: {{.ClientKey}}
`))

	Kubeconfig_template = template.Must(template.New("kubeconfig").Parse(`apiVersion: v1
clusters:
- cluster:
    server: {{.URL}}
    certificate-authority: {{.CACert}}
  name: local
contexts:
- context:
    cluster: local
    namespace: default
    user: user
  name: Default
current-context: Default
kind: Config
preferences: {}
users:
- name: user
  user:
    client-certificate: {{.ClientCert}}
    client-key: {{.ClientKey}}
`))
)
