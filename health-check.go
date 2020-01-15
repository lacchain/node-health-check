package main

import (
	"fmt"
)

func main() {
	c := fanIn(executeReadJavaProcess(), testPorts())
	for <-c != true {
		fmt.Println(<-c)
	}

	fmt.Println("Restarting", processName, "process...")
}
