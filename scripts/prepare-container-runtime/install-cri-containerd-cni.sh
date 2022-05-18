#!/bin/bash

# run as: sudo ./install-containerd.sh $version $arch
# default version=1.6.4
# default arch=amd64

version=0
arch=amd64
if [ $# -lt 2 ]
then
    version=1.6.2
    arch=amd64
else
    version=$1
    arch=$2
fi

# echo "install runc:"

# rm -r /usr/local/sbin/runc
# rm -r /usr/bin/runc

# if [[ ! -f runc.$arch ]] ; then
#     wget https://github.com/opencontainers/runc/releases/download/v1.1.0/runc.$arch
# fi

# cp runc.$arch runc
# chmod 777 runc
# mv runc /usr/local/bin/runc

# ln -s /usr/local/bin/runc /usr/local/sbin/
# ln -s /usr/local/bin/runc /usr/bin/

# echo "success to install runc."

if [[ ! -f cri-containerd-cni-$version-linux-$arch.tar.gz ]] ; then
    if [ $arch -eq "arm" ] ; then
        wget https://github.com/Litekube/LiteKube/releases/download/release-v0.1.0/cri-containerd-cni-$version-linux-$arch.tar.gz
    else
        wget https://github.com/containerd/containerd/releases/download/v$version/cri-containerd-cni-$version-linux-$arch.tar.gz
    fi 
fi

if [[ ! -f cri-containerd-cni-$version-linux-$arch.tar.gz ]] ; then
    echo "fail to download cri-containerd-cni-$version-linux-$arch.tar.gz, exit."
    exit
fi

tar -C / -xavf cri-containerd-cni-$version-linux-$arch.tar.gz

ln -s /usr/local/sbin/runc /usr/local/bin/
ln -s /usr/local/sbin/runc /usr/bin/

cat >>/etc/profile<<EOF
export PATH=$PATH:/usr/local/bin:/usr/local/sbin
EOF

cat >>~/.bashrc<<EOF
export PATH=$PATH:/usr/local/bin:/usr/local/sbin
EOF

source /etc/profile
source ~/.bashrc

mkdir -p /etc/containerd
containerd config default> /etc/containerd/config.toml
echo "default containerd config is write to \"/etc/containerd/config.toml\". If you want to config it, you can refer to <https://github.com/containerd/containerd/blob/main/docs/man/containerd-config.toml.5.md>"

sed -i 's/"bridge": "cni0"/"bridge": "containerd0"/' /etc/cni/net.d/10-containerd-net.conflist
mv /etc/cni/net.d/10-containerd-net.conflist /etc/cni/net.d/99-containerd-net.conflist

systemctl daemon-reload
systemctl enable containerd
systemctl restart containerd

echo "success to install containerd."

echo "if you meet some error while run container, try to remove your old libseccomp and install the latest version refer to scripts/prepare-container-runtime/update-libseccomp.md."

