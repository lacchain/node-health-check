package main

import (
	"fmt"
	"github.com/coreos/go-systemd/daemon"
)

func main() {
	c := fanIn(executeReadJavaProcess(), testPorts())
	daemon.SdNotify(false, "READY=1")
	restart := false
	for !restart {
		restart = <-c
	}

	fmt.Println("Restarting", processName, "process...")
}
