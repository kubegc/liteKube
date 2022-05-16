
# Catalogue
- [Catalogue](#catalogue)
- [simple start](#simple-start)
- [By dockerfile](#by-dockerfile)
# simple start
> `golang` and `gcc` environment are required

- leader

    ```shell
    cd LiteKube/cmd/leader
    go build -o leader .
    ```

- worker

    ```shell
    cd LiteKube/cmd/worker
    go build -o worker .
    ```

# By [dockerfile](Dockerfile)
> assum you start your work in `/mywork/`

1. download code from github

   ```shell
   cd /mywork
   git clone https://github.com/Litekube/LiteKube.git 
   ```

2. build image by docker

   ```shell
   cd /mywork/LiteKube/build/
   docker build -t litekube/centos-go:v1 .
   ```

   if you need proxy, you can use proxy of your host-device and run:

   ```shell
   cd /mywork/LiteKube/build/
   export http_proxy="your proxy"
   export https_proxy="your proxy"
   docker build --network=host -t litekube/centos-go:v1 .
   ```

3. start to build binaries for LiteKube

   ```shell
   chmod +x /mywork/LiteKube/scripts/build/build.sh
   docker run -v /mywork/LiteKube:/LiteKube --name=compile-litekube litekube/centos-go:v1 /LiteKube/scripts/build/build.sh
   ```

   now, you can view binaries in `/mywork/LiteKube/build/outputs/`. 
   
   > we only provide two version in this container. 
   >
   > - the same arch with your machine for Linux
   > - `Armv7l` for Linux
