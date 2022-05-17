# Catalogue

- [Catalogue](#catalogue)
- [Introduce](#introduce)
- [Command-Line parameters](#command-line-parameters)
- [`YAML` format](#yaml-format)
- [parameter specification](#parameter-specification)
  - [global](#global)
  - [`kube-apiserver`, `kube-scheduler` and `kube-controller-manager`](#kube-apiserver-kube-scheduler-and-kube-controller-manager)
  - [`kine`](#kine)
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
