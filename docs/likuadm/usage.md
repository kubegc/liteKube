# Catalogue
- [Catalogue](#catalogue)
- [Introduce](#introduce)
- [Command-Line parameters](#command-line-parameters)
# Introduce

`likuadm` is a simple tool for accessing the LiteKube control plane. Usually you need to use this tool to get the token value for the `worker` to join into the cluster.

# Command-Line parameters
> You can view the usage details directly by `--help`. Here give some key options information.

- check health
  
  ```shell
  ./likuadm health
  ```

- create bootstrap-token
  > Administrator account required. You can create token info to allow `worker` join into cluster. 
    
  ```shell
  ./likuadm create-token
  ```

  give `--life=10` to set valid only in `10` minutes(default). if `-1`, it means permanent.

- create service account
  > Administrator account required. 
  
  ```shell
  ./likuadm create-account
  ```
  
    give `--life=10` to set valid only in `10` minutes(default). if `-1`, it means permanent. If you set `--admin=true` while `--admin=false` is default, it will have the same permissions as you currently have, that is, full control.

  To specify the use of an account, you can specify in `[global]: --token`

- list service accounts
  > Administrator account required. 

  ```shell
  ./likuadm list-accounts
  ```

- delete account
  > Administrator account required. 

  ```shell
  ./likuadm delete-account --token=$(delete-token)
  ```
