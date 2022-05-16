# Catalogue

- [Catalogue](#catalogue)
- [Deployment](#deployment)
  - [Introduce](#introduce)
  - [Configuration](#configuration)
    - [Indispensable](#indispensable)
    - [Enable Kine in leader (default: true)](#enable-kine-in-leader-default-true)
    - [Enable Network-Controller (default: true)](#enable-network-controller-default-true)
    - [Enable Worker (default: false)](#enable-worker-default-false)
# Deployment

## Introduce

`leader` is an important component for `LiteKube`. At its most basic, it contains `Kube-Apiserver`, `Kube-Scheduler`, and `Kube-Controller` for `k8s`, as well as `LiteKube`'s `network part` and a `control component`. For ease of use, `leader` allow `kine`(A lightweight similar to ETCD created by k3s), `network-Controller` and `worker` component to be automatically configured internally in an integrated manner.

## Configuration
By default, the leader does not require any additional running environment, just firewall rules .

You will need to allow the following ports to be access:

### Indispensable
- Exposing firewall Ports
  - 6440-6443
  - 10257
  - 10259

### Enable Kine in leader (default: true)
- Exposing firewall Ports
  - 2379

### Enable Network-Controller (default: true)
- Exposing firewall Ports
  - 6439

### Enable Worker (default: false)
- please refer directly to [worker deployment](../worker/deploy.md) and configuration requirements, all in need.
