package main

import (
	"flag"
	"fmt"
	"github.com/reiver/go-telnet"
	"net/http"
	"os"
)

func main() {
	msg1 := TestClientUrl("http://127.0.0.1", "4444")
	msg2 := TestNodePort("127.0.0.1", "4040")
	if msg1 || msg2 {
		fmt.Println("1") //restart
	} else {
		fmt.Print("0")
	}
	os.Exit(0)
}

func TestClientUrl(url, clientPort string) bool {
	port := flag.String("port", clientPort, "port on localhost to check")
	flag.Parse()

	resp, err := http.Get(url + ":" + *port + "/upcheck")
	if err != nil || resp.StatusCode != 200 {
		//fmt.Println("dead in client")
		return true
	}
	return false
}

func TestNodePort(ip, nodePort string) bool {
	conn, err := telnet.DialTo(ip + ":" + nodePort)
	if err != nil {
		//fmt.Println("NodeUrl not responding")
		return true
	}
	conn.Write([]byte("hello world"))
	conn.Write([]byte("\n"))
	return false
}
