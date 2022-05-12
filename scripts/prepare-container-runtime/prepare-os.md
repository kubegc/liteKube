> Notice:
>
> * These commands may not be applicable to all systems
>
> * Required only by nodes that are required to run the `worker`

* close your firewall 

    ```shell
    systemctl stop firewalld
    systemctl disable firewalld
    ```

    > or you can try to enable ports for default: [2379, 6440-6443, 10257, 10259, 30000-32767] and other port you many need. This function is only experimental.

*   setting time synchronization for all cluster-nodes

 ```bash
yum install -y ntpdate		# for centos
# apt-get install ntpdate	# for ubuntu

ntpdate time.windows.com

# validate
date
 ```

* close selinux
```bash
sed -i 's/enforcing/disabled/' /etc/selinux/config

# reboot your machine
reboot

# validate selinux
getenforce

# get disable
```

* 关闭swap分区

```bash
sed -ri 's/.*swap.*/#&/' /etc/fstab

# 重启
reboot

#验证swap分区
free -m
```

* To pass the bridged ` IPv4` traffic to the chain of `iptables`, `k8s` may need 
```bash
cat > /etc/sysctl.d/k8s.conf << EOF
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
vm.swappiness = 0
EOF

# load br_netfilter modules
modprobe br_netfilter

# validate 
lsmod | grep br_netfilter

# take effect
sysctl --system
```

*   open ipvs
```bash
# install
yum -y install ipset ipvsadm	# for centos
apt-get install ipset ipvsadm	# for ubuntu

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
