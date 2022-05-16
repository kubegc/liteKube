# Catalogue

- [Catalogue](#catalogue)
- [Deployment](#deployment)
  - [Introduce](#introduce)
  - [Configuration](#configuration)
# Deployment

## Introduce

`kubectl`, the kubernetes command-line tool, allows you to run commands against Kubernetes clusters. LiteKube do no change to this part and running environment is configed automatically after leader starts.

## Configuration
you can add `kubectl` binary to your `$PATH`, such as:

```shell
mv ./kubectl /usr/bin/
```

once `leader` run ok, `kubectl` will be ready.