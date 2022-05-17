FROM centos:7

# docker run -v $ProjectPath:/LiteKube 
LABEL maintainer="a_flying_fish@outlook.com"

# build golang environment 
WORKDIR /

RUN yum install wget -y && \
    yum install git -y && \
    yum install tar -y && \
    yum install gcc -y

RUN mkdir -p /usr/local && \
    mkdir -p /home/go && \
    wget --no-check-certificate https://go.dev/dl/go1.18.2.linux-amd64.tar.gz && \
    tar -C /usr/local/ -zxvf go1.18.2.linux-amd64.tar.gz && \
    rm -r go1.18.2.linux-amd64.tar.gz

ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=/home/go
ENV GOROOT=/usr/local/go
ENV GOPROXY=goproxy.cn

RUN go version

# build arm-gcc enviroment
RUN yum install glibc.i686 -y && \
    yum install fontconfig.i686 -y && \
    yum install libXrender.i686 -y && \
    yum install libstdc++.i686 -y

RUN mkdir -p /usr/local/toolchain/ && \ 
    wget http://releases.linaro.org/components/toolchain/binaries/7.5-2019.12/arm-linux-gnueabihf/gcc-linaro-7.5.0-2019.12-i686_arm-linux-gnueabihf.tar.xz && \
    xz -d gcc-linaro-7.5.0-2019.12-i686_arm-linux-gnueabihf.tar.xz && \
    tar xvf gcc-linaro-7.5.0-2019.12-i686_arm-linux-gnueabihf.tar -C /usr/local/toolchain/ && \
    rm -r gcc-linaro-7.5.0-2019.12-i686_arm-linux-gnueabihf.tar*

ENV PATH=$PATH:/usr/local/toolchain/gcc-linaro-7.5.0-2019.12-i686_arm-linux-gnueabihf/bin

RUN mkdir -p /LiteKube

WORKDIR /LiteKube

# RUN chmod +x scripts/build/build.sh
CMD ["bin/bash"]
