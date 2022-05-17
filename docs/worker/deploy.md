# Catalogue

- [Catalogue](#catalogue)
- [Introduce](#introduce)
- [Configuration](#configuration)
  - [set up os](#set-up-os)
  - [install containerd](#install-containerd)
  - [adapt containerd to `LiteKube`](#adapt-containerd-to-litekube)

# Introduce

`worker` is an important component for `LiteKube`. At its most basic, it contains `kubelet` and `Kube-Proxy` for `k8s`, as well as `LiteKube`'s `network part`. 

# Configuration

We will provide the recommended reference configuration, if you have some knowledge of kubernetes configuration parameters, you can also follow your own custom settings.

## set up os
> The following commands are only tested in `Centos 7.9`. You can replace them with the local settings. [script](../../scripts/prepare-container-runtime/prepare-os.md) for centos is ready.

- close firewall and disable reboot
  > can be replaced with firewall rules

  ```shell
  systemctl stop firewalld
  systemctl disable firewalld
  ```

- set data-time synchronization
  > run in all machines

  ```shell
  yum install -y ntpdate
  ntpdate time.windows.com
  ```

  validate
  
  ```shell
  date
  ```

- close selinux
  
  ```shell
  sed -i 's/enforcing/disabled/' /etc/selinux/config
  reboot
  ```
  
  validate

  ```shell
  # if not exist, you can install "selinux-utils"

  getenforce
  ```

- close swap memory

  ```shell
  sed -ri 's/.*swap.*/#&/' /etc/fstab
  reboot
  ```
  
  validate

  ```shell
  free -m
  ```

- pass the bridged IPv4 traffic to the chain of Iptables
  
  ```shell
  cat > /etc/sysctl.d/k8s.conf << EOF
  net.bridge.bridge-nf-call-ip6tables = 1
  net.bridge.bridge-nf-call-iptables = 1
  net.ipv4.ip_forward = 1
  vm.swappiness = 0 
  EOF

  # load
  modprobe br_netfilter

  # validate load
  lsmod | grep br_netfilter

  # to take effect
  sysctl --system
  ```

- enable ipvs
  
  ```shell
  # install
  yum -y install ipset ipvsadm

  # config
  cat > /etc/sysconfig/modules/ipvs.modules <<EOF
  #!/bin/bash
  modprobe -- ip_vs
  modprobe -- ip_vs_rr
  modprobe -- ip_vs_wrr
  modprobe -- ip_vs_sh
  modprobe -- nf_conntrack_ipv4
  EOF

  chmod 755 /etc/sysconfig/modules/ipvs.modules && bash /etc/sysconfig/modules/ipvs.modules && lsmod | grep -e ip_vs -e nf_conntrack_ipv4

  # validate
  lsmod | grep -e ipvs -e nf_conntrack_ipv4
  ```

## install containerd
[Script](../../scripts/prepare-container-runtime/install-cri-containerd-cni.sh) for Linux is ready, we only tested it on `Centos 7.9` and you can adapt it to your local system.

- install `runc`
  
  ```shell
  wget https://github.com/opencontainers/runc/releases/download/v1.1.0/runc.amd64 && mv runc.amd64 runc && chmod 777 runc

  mv runc /usr/local/bin/runc
  ln -s /usr/local/bin/runc /usr/local/sbin/
  ln -s /usr/local/bin/runc /usr/bin/
  ```

- download `cri-containerd-cni` and unpack to path
  > notice the system and architecture. We take Linux-amd64 as an example.

  ```shell
  wget https://github.com/containerd/containerd/releases/download/v1.6.2/cri-containerd-cni-1.6.2-linux-amd64.tar.gz

  tar -C / -xavf cri-containerd-cni-$version-linux-$arch.tar.gz
  ```

- add to $PATH
  
  ```shell
  cat >>/etc/profile<<EOF
  export PATH=$PATH:/usr/local/bin:/usr/local/sbin
  EOF

  cat >>~/.bashrc<<EOF
  export PATH=$PATH:/usr/local/bin:/usr/local/sbin
  EOF

  source /etc/profile
  source ~/.bashrc
  ```

  you can consider run `yum install -y containerd` for `Centos` or `apt-get install containerd` for `Ubuntu` instead.

- create config-file

  ```shell
  mkdir -p /etc/containerd
  containerd config default> /etc/containerd/config.toml
  ```

- set default `bridge` name from `cni0` to `containerd0` and reduce priority
  
  ```shell
  sed -i 's/"bridge": "cni0"/"bridge": "containerd0"/' /etc/cni/net.d/10-containerd-net.conflist
  mv /etc/cni/net.d/10-containerd-net.conflist /etc/cni/net.d/99-containerd-net.conflist
  ```

- to take effect
  
  ```shell
  systemctl daemon-reload 
  systemctl enable containerd
  systemctl restart containerd
  ```

- validate status
  
  ```shell
  systemctl status containerd -l
  ```

  if you meet some error while run containerd, try to remove your old libseccomp and install the latest version refer to [doc](scripts/prepare-container-runtime/update-libseccomp.md).

## adapt containerd to `LiteKube`
> [Script](../../scripts/prepare-container-runtime/containerd-to-k8s.sh) for Linux is ready, we only tested it on `Centos 7.9` and you can adapt it to your local system.

- set modules
  
  ```shell
  cat > /etc/modules-load.d/containerd.conf <<EOF
  overlay
  br_netfilter
  EOF

  modprobe overlay
  modprobe br_netfilter
  ```

- set cri-config
  
  ```shell
  cat > /etc/sysctl.d/99-kubernetes-cri.conf <<EOF
  net.bridge.bridge-nf-call-iptables  = 1
  net.ipv4.ip_forward                 = 1
  net.bridge.bridge-nf-call-ip6tables = 1
  EOF
  
  sysctl --system
  ```

- set `SystemdCgroup=true`

  ```shell
  sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml
  ```

- set image source to `registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.6`
   > optional, if you are in China, this is recommended
  
  ```shell
  sed -i 's/sandbox_image = "k8s.gcr.io\/pause:3.6"/sandbox_image = "registry.cn-hangzhou.aliyuncs.com\/google_containers\/pause:3.6"/' /etc/containerd/config.toml
  ```

- to take effect
  
  ```shell
  sudo systemctl restart containerd
  ```
