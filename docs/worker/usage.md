# Catalogue

- [Catalogue](#catalogue)
- [Introduce](#introduce)
- [Command-Line parameters](#command-line-parameters)
- [`YAML` format](#yaml-format)
- [parameter specification](#parameter-specification)
  - [global](#global)
  - [`kubelet` and `kube-proxy`](#kubelet-and-kube-proxy)
  - [`network-manager`](#network-manager)
# Introduce

`worker` is an important component for `LiteKube`. At its most basic, it contains `kubelet` and `Kube-Proxy` for `k8s`, as well as `LiteKube`'s `network part`. 

# Command-Line parameters

Because there are so many parameters that can be set, and some components have similar parameter meanings. We use completely file-based parameter input, and only a few necessary functions take command-line arguments. such as:

- get help and view allowed args:

    ```shell
    ./worker --help
    ```

- view related verions

    ```shell
    ./worker --versions
    ```

- run `worker`

    ```shell
    ./worker --config-file=/path-to/worker.yaml
    ```

# `YAML` format

```yaml
global:
    # leader startup parameters and common args for kubernetes components
    leader-token: string
    log-dir:      string
    log-to-dir:   bool  
    log-to-std:   bool  
    work-dir:     string
kubelet:
    # kubelet's startup parameters
    options:
        # Litekube normal options
        cert-dir:                  string
        pod-infra-container-image: string
    professional:
        # parameters are not recommended to set by users
        bootstrap-kubeconfig:       string
        cgroup-driver:              string
        config:                     string
        container-runtime:          string
        container-runtime-endpoint: string
        hostname-override:          string
        kubeconfig:                 string
        runtime-cgroups:            string
    reserve:
        # reserve parameters
        <name-1>: <value-1>
        <name-n>: <value-n>
kube-proxy:
    # kube-proxy's startup parameters
    options:
        # Litekube normal options
    professional:
        # parameters are not recommended to set by users
        cluster-cidr:      string
        hostname-override: string
        kubeconfig:        string
        proxy-mode:        string
    reserve:
        # reserve parameters
        <name-1>: <value-1>
        <name-n>: <value-n>
network-manager:
    # network register and manager component for litekube
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
```

You can go straight to a real YAML startup [config-file](../examples/worker.yaml)

# parameter specification
> you can set only part or even none of yaml-config unless `global.leader-token` and `network-manager.token` (Simply, you can get these two value by run `likuadm create-token` in `leader` node), there will be good default-value as usual. 
>
> You can view the actual startup configuration in `global.work-dir/startup/worker.yaml` or even use it directly as a new startup configuration file
## global
- leader-token
  > important args for worker. LiteKube via this parameters
   to run `bootstrap` and help worker ready.
- work-dir
  > we will try our best to set file-cache to this directory instead of `$HOME/.litekube`. Special files will still be saved in `$HOME/.litekube`, but they are  trivial.
- log-dir
  > if set `global.log-to-dir=true`, `log files` can redirect to independent path instead of `work-dir/logs`
## `kubelet` and `kube-proxy`
> all args will finally be explain as `--<key>=<value>`
- professional
  > if you are not familiar with kubernetes startup parameters and architecture of LiteKube, ignoring these Settings is recommended.
- options
  > parameters that you may normally still want to set, although default values are given.
- reserve
  > parameters that Litekube doesn't mention can still be set by key-value pairs. They will also be explain as `--<key>=<value>` for `kubernetes component`, we do none check to these parameters.
  > 
  > parameters mentioned in `options` or `professional` will be ignored. Usually, we will print tips for you.

## `network-manager`
> *(Due to the history of the program, we have retained the old name in our code and comments.It's essentially equivalent to `network-controller`)*
> 
> Take the complexity of the network components into consideration, we established the `bootstrap` mechanism for convenience. You only need to set the value for `token` here.
- token
  > Simply, you can get this value by run `likuadm create-token` in `leader` node
