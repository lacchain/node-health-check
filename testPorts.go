package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/reiver/go-telnet"
)

var port *string

func testPorts() <-chan bool {
	port = flag.String("port", clientPort, "port on localhost to check")
	flag.Parse()
	c := make(chan bool)
	go func() {
		for i := 0; ; i++ {
			time.Sleep(delayTimeMinutes * time.Minute)
			msg := test()
			if msg {
				c <- msg //true ==> restart
			}			
		}
	}()

	return c
}

func test() bool {
	msg1 := testClientURL(clientURL)
	msg2 := testNodePort(nodeURL, nodePort)
	return msg1 || msg2
}

func testClientURL(url string) bool {
	resp, err := http.Get(url + ":" + *port + "/upcheck")
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("dead in client")
		return true
	}
	fmt.Println("Successful response from ClientPort:", clientPort)
	defer resp.Body.Close()
	return false
}

func testNodePort(ip, nodePort string) bool {
	conn, err := telnet.DialTo(ip + ":" + nodePort)
	if err != nil {
		fmt.Println("NodeUrl not responding")
		return true
	}
	fmt.Println("Successful response from Nodeport:", nodePort)
	conn.Write([]byte("hello world"))
	conn.Write([]byte("\n"))
	return false
}
