**Catalogues:**

[TOC]

# deploy leader

leader is one important component for `LiteKube`. At its most basic, it contains `Kube-Apiserver`, Kube-Scheduler, and Kube-Controller for `k8s`, as well as `LiteKube`'s network plug-in and a control component. For ease of use, the leader allows the `Kine(A lightweight similar to ETCD created by k3s)`, `network-Controller` and `worker` component to be automatically configured internally in an integrated manner.

## 1. Command-Line parameters

Because there are so many parameters that can be set, and some components have similar parameter meanings. We use completely file-based parameter input, and only a few necessary functions take command-line arguments. such as ：

* get help and view allowed args:

    ```shell
    ./leader --help
    ```

* view related verions

    ```shell
    ./leader --versions
    ```

* run `leader`

    ```shell
    ./leader --config-file=/path-to/config.yaml
    ```

`YAML` format:

```yaml
global:
    # leader startup parameters and common args for kubernetes components
    enable-worker:       bool  
    log-dir:             string
    log-to-dir:          bool  
    log-to-std:          bool  
    run-kine:            bool  
    run-network-manager: bool  
    work-dir:            string
network-manager:
    # network-controller component for litekube
    node-token: string
    token:      string
    join:
        # to be joined and managered. certificates need to be given together with --node-token
        ca-cert:          string
        client-cert-file: string
        client-key-file:  string
        network-address:  string
        secure-port:      uint16
    register:
        # to register and query from manager. certificates need to be given together with --node-token. Or you can only 
        ca-cert:          string
        client-cert-file: string
        client-key-file:  string
        network-address:  string
        secure-port:      uint16
kube-apiserver:
    # kube-Apiserver's startup parameters
    options:
        # Litekube normal options
        allow-privileged:           bool  
        anonymous-auth:             bool  
        authorization-mode:         string
        enable-admission-plugins:   string
        encryption-provider-config: string
        profiling:                  bool  
        secure-port:                uint16
        service-cluster-ip-range:   string
        service-node-port-range:    string
    professional:
        # parameters are not recommended to set by users
        advertise-address:                  string
        api-audiences:                      string
        bind-address:                       string
        cert-dir:                           string
        client-ca-file:                     string
        enable-aggregator-routing:          bool  
        enable-bootstrap-token-auth:        bool  
        etcd-cafile:                        string
        etcd-certfile:                      string
        etcd-keyfile:                       string
        etcd-servers:                       string
        feature-gates:                      string
        kubelet-certificate-authority:      string
        kubelet-client-certificate:         string
        kubelet-client-key:                 string
        proxy-client-cert-file:             string
        proxy-client-key-file:              string
        requestheader-allowed-names:        string
        requestheader-client-ca-file:       string
        requestheader-extra-headers-prefix: string
        requestheader-group-headers:        string
        requestheader-username-headers:     string
        service-account-issuer:             string
        service-account-key-file:           string
        service-account-signing-key-file:   string
        storage-backend:                    string
        tls-cert-file:                      string
        tls-private-key-file:               string
        token-auth-file:                    string
    reserve:
        # reserve parameters
        <name-1>: <value-1>
        <name-n>: <value-n>
kube-controller-manager:
    # kube-controller-manager's startup parameters
    options:
        # Litekube normal options
        allocate-node-cidrs:             bool  
        cluster-cidr:                    string
        profiling:                       bool  
        use-service-account-credentials: bool  
    professional:
        # parameters are not recommended to set by users
        authentication-kubeconfig:                       string
        authorization-kubeconfig:                        string
        bind-address:                                    string
        cluster-signing-kube-apiserver-client-cert-file: string
        cluster-signing-kube-apiserver-client-key-file:  string
        cluster-signing-kubelet-client-cert-file:        string
        cluster-signing-kubelet-client-key-file:         string
        cluster-signing-kubelet-serving-cert-file:       string
        cluster-signing-kubelet-serving-key-file:        string
        cluster-signing-legacy-unknown-cert-file:        string
        cluster-signing-legacy-unknown-key-file:         string
        configure-cloud-routes:                          bool  
        controllers:                                     string
        feature-gates:                                   string
        kubeconfig:                                      string
        leader-elect:                                    bool  
        root-ca-file:                                    string
        secure-port:                                     uint16
        service-account-private-key-file:                string
    reserve:
        # reserve parameters
        <name-1>: <value-1>
        <name-n>: <value-n>
kube-scheduler:
    # kube-scheduler's startup parameters

    options:
        # Litekube normal options
        profiling: bool
    professional:
        # parameters are not recommended to set by users
        authentication-kubeconfig: string
        authorization-kubeconfig:  string
        bind-address:              string
        kubeconfig:                string
        leader-elect:              bool  
        secure-port:               uint16
    reserve:
        # reserve parameters
        <name-1>: <value-1>
        <name-n>: <value-n>
kine:
    # lite-Database for litekube
    bind-address:     string
    ca-cert:          string
    secure-port:      uint16
    server-cert-file: string
    server-key-file:  string,  SSL key file used to secure etcd communication.
```









