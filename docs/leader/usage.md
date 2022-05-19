# Catalogue

- [Catalogue](#catalogue)
- [Introduce](#introduce)
- [Simple start](#simple-start)
- [Command-Line parameters](#command-line-parameters)
- [`YAML` format](#yaml-format)
- [parameter specification](#parameter-specification)
  - [global](#global)
  - [`kube-apiserver`, `kube-scheduler` and `kube-controller-manager`](#kube-apiserver-kube-scheduler-and-kube-controller-manager)
  - [`kine`](#kine)
  - [`network-manager`](#network-manager)
- [install CNI](#install-cni)
# Introduce

`leader` is an important component for `LiteKube`. At its most basic, it contains `Kube-Apiserver`, `Kube-Scheduler`, and `Kube-Controller` for `k8s`, as well as `LiteKube`'s `network part` and a `control component`. For ease of use, `leader` allow `kine`(A lightweight similar to ETCD created by k3s), `network-Controller` and `worker` component to be automatically configured internally in an integrated manner.

# Simple start

Download leader, likuadm and kubectl, then you can start leader simply by:

```shell
mv ./leader* leader && chmod +x ./leader && mv ./leader /usr/bin/
mv ./kubectl* kubectl && chmod +x ./kubectl && mv ./kubectl /usr/bin/
mv ./likuadm* likuadm && chmod +x ./likuadm && mv ./likuadm /usr/bin/
mkdir -p /opt/litekube/

cat >/opt/litekube/leader.yaml <<EOF
global:
  log-to-dir: true
EOF

cat >/etc/systemd/system/leader.service <<EOF
[Unit]
Description=LiteKube leader

[Service]
ExecStart=/usr/bin/leader --config-file=/opt/litekube/leader.yaml
Restart=on-failure
KillMode=process
[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable leader
systemctl restart leader
```

goto [install CNI](#install-cni)

# Command-Line parameters

Because there are so many parameters that can be set, and some components have similar parameter meanings. We use completely file-based parameter input, and only a few necessary functions take command-line arguments. such as:

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

# `YAML` format

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

# parameter specification
> you can set only part or even none of yaml-config, there will be good default-value as usual.
>
> You can view the actual startup configuration in `global.work-dir/startup/leader.yaml` or even use it directly as a new startup configuration file
## global
- work-dir
  > we will try our best to set file-cache to this directory instead of `$HOME/.litekube`. Special files will still be saved in `$HOME/.litekube`, but they are  trivial.
- log-dir
  > if set `global.log-to-dir=true`, `log files` can redirect to independent path instead of `work-dir/logs`
- run-kine
  > LiteKube run `kine` as one lite-etcd to start for default. You can set to `false` of course. Instead you need to give `ETCD Args` for `kube-apiserver`
- run-network-manager
  > LiteKube run `network-controller server` for default. This results in `leader` still have to run at the top-level of the network. 
  >
  > You can set this value to `false` and set `network-manager.token`(manually configuring network-manager parameters is not recommended). With the help of network-controller-bootstrap, leader can config network-manager args automatically. `leader` will be able to run in any network-level once seperate `network-controller server` with `leader` like `worker`
- enable-worker
  > if you want to run `leader` and `worker` in same node, set this value to `true` and `false` is default. This need you prepare `worker` running environment for `leader`. if you want to give your own setting, you many need to change `global.work-dir/startup/worker.yaml` and restart `leader`.
## `kube-apiserver`, `kube-scheduler` and `kube-controller-manager`
> all args will finally be explain as `--<key>=<value>`
- professional
  > if you are not familiar with kubernetes startup parameters and architecture of LiteKube, ignoring these Settings is recommended.
- options
  > parameters that you may normally still want to set, although default values are given.
- reserve
  > parameters that Litekube doesn't mention can still be set by key-value pairs. They will also be explain as `--<key>=<value>` for `kubernetes component`, we do none check to these parameters.
  > 
  > parameters mentioned in `options` or `professional` will be ignored. Usually, we will print tips for you.
## `kine`
> you can discard these parameters if you set `global.run-kine=false`. Or you can partially customize kine Server parameters. Unless you set up the certificates manually, they will be generated automatically.
>
> Notice: `kube-Apiserver`'s ETCD parameters need to be configured manually once you manually set up the certificate, as we lack the necessary information to complete the automation.
## `network-manager`
> *(Due to the history of the program, we have retained the old name in our code and comments.It's essentially equivalent to `network-controller`)*
> 
> Take the complexity of the network components into consideration, we established the `bootstrap` mechanism for convenience. You only need to set the value for `token` here.
- token
  > By default, this value does not need to be configured for `leader`. But if you choose to separate the `leader` and `network controller server`, it will be necessary to run `ncadm` on node running `network-controller server` to get the token value. 

# install CNI
> Current `leader` version does not install network plug-ins such as `flannel` by default. You can install by:

```shell
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
```
