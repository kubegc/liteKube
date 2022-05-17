# Catalogue

- [Catalogue](#catalogue)
- [Introduce](#introduce)
- [Configuration](#configuration)
  - [indispensable](#indispensable)
  - [enable kine](#enable-kine)
  - [enable network-controller](#enable-network-controller)
  - [enable worker](#enable-worker)

# Introduce

`leader` is an important component for `LiteKube`. At its most basic, it contains `Kube-Apiserver`, `Kube-Scheduler`, and `Kube-Controller` for `k8s`, as well as `LiteKube`'s `network part` and a `control component`. For ease of use, `leader` allow `kine`(A lightweight similar to ETCD created by k3s), `network-Controller` and `worker` component to be automatically configured internally in an integrated manner.

# Configuration
The operation of leader has a certain complexity. According to its own operation, it has different requirements for the environment, which will be discussed in modules below:

## indispensable
- Exposing firewall Ports
  - 6440-6443
  - 10257
  - 10259

## enable kine
> default: true
- Exposing firewall Ports
  - 2379

## enable network-controller
> default: true
- Exposing firewall Ports
  - 6439
- refer to https://github.com/Litekube/network-controller

## enable worker
> default: false
- please refer directly to [worker deployment](../worker/deploy.md) and configuration requirements, all in need.
