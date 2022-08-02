# Node Health Check #

## Introduction
* The node-health-check service can be used to guarante availability of the **orion** transaction manager. The code aims to work hand in hand with the operating system. If for some reason the Orion service fails, then the node-health-check service restarts it. 

* The installation of this service has been automated in the recommended process for the deployment of new nodes described in the [lacchain network](https://github.com/lacchain/besu-network) repository, that uses Ansible. Therefore, it's being used in the [besu-network](https://github.com/lacchain/besu-network) by **writer** nodes that execute private transactions.

* The node-health-check service monitors the following:
    1. Orion Node url
    2. Orion Client url
    3. Orion Java process: Monitors the heap and old space.

## Requirements
* Golang

### Clone Repository ####
```shell
$ https://github.com/lacchain/node-health-check.git
$ cd node-health-check/
```

### Build process the code ###
```shell
$ export GO111MODULE=off && go get ./... && go build -o health-check
```

### Running the code ###
Run the code with:
```shell
$ go run health-check
```

### Looking for logs ###
If you have installed this software by using the ansible from the [besu-network](https://github.com/lacchain/besu-network), you can check the logs by using:
```shell
$ journalctl -fu health-check.service
```
