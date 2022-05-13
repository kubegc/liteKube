<h1 align="center">How to build</h1>

1. you can simple complie as follow:

    > `golang` and `gcc` environment are required

    * leader

        ```shell
        cd LiteKube/cmd/leader
        go build -o leader .
        ```

    * worker

        ```shell
        cd LiteKube/cmd/worker
        go build -o worker .
        ```

2. We still recommend that you use container-based compilation, and we provide a [Dockerfile](Dockerfile) to help you build your image. You can follow the rules for writing DockerFile and compile them to your schema, which may require a little knowledge of container.  The current version of the image, at least, provides a native go-compilation environment and a cross-compilation environment for the `armv7l` architecture. You can run by `Docker`as follow:

    ```shell
    # assum you start your work in /mywork/
    
    # download project first
    cd /mywork
    git clone https://github.com/Litekube/LiteKube.git 
    
    cd /mywork/LiteKube/build/
    docker build -t litekube/centos-go:v1 .
    # if you need proxy, you can use proxy of your host-device and run:
    # docker build --network=host -t  litekube/centos-go:v1 .
    
    chmox +x /mywork/LiteKube/scripts/build/build.sh
    docker run -v /mywork/LiteKube:/LiteKube --name=compile-litekube litekube/centos-go:v1 /LiteKube/scripts/build/build.sh
    ```

    then you can view binary in `/mywork/LiteKube/build/outputs/`