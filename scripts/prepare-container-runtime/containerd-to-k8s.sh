#!/bin/bash

# run as: sudo ./containerd-to-k8s.sh

cat > /etc/modules-load.d/containerd.conf <<EOF
overlay
br_netfilter
EOF

modprobe overlay
modprobe br_netfilter

cat > /etc/sysctl.d/99-kubernetes-cri.conf <<EOF
net.bridge.bridge-nf-call-iptables  = 1
net.ipv4.ip_forward                 = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF

sysctl --system
# echo ""
# echo "-------------------------------- --notice------------------------------------------------"
# echo "|                                                                                       |"
# echo "| something here need your hand to finish the work, follow:                             |"
# echo "| 1. edit \"/etc/containerd/config.toml\"                                               |"
# echo "| 2. set runc with <systemd cgroup>, usually it look like:                              |"
# echo "|    [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]                     |"
# echo "|      ...                                                                              |"
# echo "|      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]           |"
# echo "|        SystemdCgroup = true # set to <true> please                                    |" 
# echo "| 3. run command: \"sudo systemctl restart containerd\", and you will finish your work. |"
# echo "| 4. run command: \"sudo systemctl status containerd  -l \" to check your work.         |"
# echo "-----------------------------------------------------------------------------------------"

sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml
sed -i 's/sandbox_image = "k8s.gcr.io\/pause:3.6"/sandbox_image = "registry.cn-hangzhou.aliyuncs.com\/google_containers\/pause:3.6"/' /etc/containerd/config.toml
sudo systemctl restart containerd