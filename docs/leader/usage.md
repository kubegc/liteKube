# Catalogue

- [Catalogue](#catalogue)
- [Usage](#usage)
  - [Introduce](#introduce)
  - [Command-Line parameters](#command-line-parameters)
  - [`YAML` format](#yaml-format)
  - [parameter specification](#parameter-specification)
# Usage

## Introduce

`leader` is an important component for `LiteKube`. At its most basic, it contains `Kube-Apiserver`, `Kube-Scheduler`, and `Kube-Controller` for `k8s`, as well as `LiteKube`'s `network part` and a `control component`. For ease of use, `leader` allow `kine`(A lightweight similar to ETCD created by k3s), `network-Controller` and `worker` component to be automatically configured internally in an integrated manner.

## Command-Line parameters

Because there are so many parameters that can be set, and some components have similar parameter meanings. We use completely file-based parameter input, and only a few necessary functions take command-line arguments. such as ：

- get help and view allowed args:

    ```shell
    ./leader --help
    ```

- view related verions

    ```shell
    ./leader --versions
    ```

- run `leader`

    ```shell
    ./leader --config-file=/path-to/config.yaml
    ```

## `YAML` format

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

You can go straight to a real YAML startup [config-file](../examples/leader.yaml)

## parameter specification
> you can set only part or even none of yaml-config, there will be good default-value as usual.
>
> You can view the actual startup configuration in `global.work-dir/startup/leader.yaml` or even use it directly as a new startup configuration file
- global
  - work-dir
    > we will try our best to set file-cache to this directory instead of `$HOME/.litekube`. Special files will still be saved in `$HOME/.litekube`, but they are  trivial.
  - log-dir
    > if set `global.log-to-dir=true`, `log files` can redirect to independent path instead of `work-dir/logs`
  - run-kine
    > LiteKube run `kine` as one lite-etcd to start for default. You can set to `false` of course. Instead you need to give `ETCD Args` for `kube-apiserver`
  - run-network-manager
    > LiteKube run `network-controller server` for default. This results in Leader still have to run at the top-level of the network. 
    >
    > You can set this value to `false` and set `network-manager.token`(manually configuring network-manager parameters is not recommended). With the help of network-controller-bootstrap, leader can config network-manager args automatically. `leader` will be able to run in any network-level once Seperate `network-controller server` with `leader` like `worker`
  - enable-worker
    > if you want to run `leader` and `worker` in same node, set this value to `true` and `false` is for default. This need you prepare worker running environment for `leader`. if you want to give your own setting, you many need to change `global.work-dir/startup/worker.yaml` and restart `leader`.
- `kube-apiserver`, `kube-scheduler` and `kube-controller-manager`
  - professional
    > if you are not familiar with kubernetes startup parameters and architecture of LiteKube, ignoring these Settings is recommended.
  - options
    > parameters that you may normally still want to set, although default values are given.
  - reserve
    > parameters that Litekube doesn't mention can still be set by key-value pairs. They will be explain as `--<key>=<value>` for `kubernetes component`
