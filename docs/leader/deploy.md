# Catalogue

- [Catalogue](#catalogue)
- [Introduce](#introduce)
- [Configuration](#configuration)
  - [Indispensable](#indispensable)
  - [enable kine](#enable-kine)
  - [enable network-controller](#enable-network-controller)
  - [enable worker](#enable-worker)

# Introduce

`leader` is an important component for `LiteKube`. At its most basic, it contains `Kube-Apiserver`, `Kube-Scheduler`, and `Kube-Controller` for `k8s`, as well as `LiteKube`'s `network part` and a `control component`. For ease of use, `leader` allow `kine`(A lightweight similar to ETCD created by k3s), `network-Controller` and `worker` component to be automatically configured internally in an integrated manner.

# Configuration
By default, the leader does not require any additional running environment, just firewall rules .

You will need to allow the following ports to be access:

## Indispensable
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

## enable worker
> default: false
- please refer directly to [worker deployment](../worker/deploy.md) and configuration requirements, all in need.
