> assume LiteKube is store in fold: $FOLD

# Easy to start

```shell
cd build/
docker build -t litekube/centos-go:v1 .

# if you need proxy, you can set proxy for your host and run:
# docker build --network=host -t  litekube/centos-go:v1 .

chmox +x $ProjectPath/scripts/build/build.sh
docker run -v $ProjectPath:/LiteKube --name=compile-litekube litekube/centos-go:v1 /LiteKube/scripts/build/build.sh
```
