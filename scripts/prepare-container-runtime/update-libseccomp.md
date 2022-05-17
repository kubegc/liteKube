## update libseccomp

The latest version of Containerd (2022-3-28) requires libseccomp &gt;`2.4`, but centos 7.9 is `2.3` by default, we manually install (must uninstall the system built-in, otherwise it will not take effect)

## take Centos 7.9 as an example

```shell
# view verison now
rpm -qa | grep libseccomp
libseccomp-devel-2.3.1-4.el7.x86_64
libseccomp-2.3.1-4.el7.x86_64

# uninstall
rpm -e libseccomp-devel-2.3.1-4.el7.x86_64 --nodeps
rpm -e libseccomp-2.3.1-4.el7.x86_64 --nodeps

# install
yum install gcc make gperf
wget https://github.com/seccomp/libseccomp/releases/download/v2.5.3/libseccomp-2.5.3.tar.gz
tar -zxvf libseccomp-2.5.3.tar.gz
cd libseccomp-2.5.3
./configure --prefix=/usr --disable-static
make && make install

systemctl restart containerd
