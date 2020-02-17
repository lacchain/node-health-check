# Node Health Check #

## Introduction
* This software guarantees availability of **orion** transaction manager. The code aims to work hand in hand with the operating system. If for some reason Orion fails then the software restarts orion service. 

* Currently this code is being used in the [lacchain network](https://github.com/lacchain/besu-network) by **writer** nodes that execute private transactions. The installation has been automated in the lacchain-network repository by using ansible. 

* When running, the software monitors the following:
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
If you have installed this software by using the  ansible for [lacchain network](https://github.com/lacchain/besu-network) then you can check the logs by using:
```shell
$ journalctl -fu health-check.service
```