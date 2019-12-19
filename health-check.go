package main

import (
	"flag"
	"fmt"
	"github.com/reiver/go-telnet"
	"net/http"
	"os"
)

func main() {
	msg := "0"
	msg = TestClientUrl("4444")
	msg = TestNodePort("4040")
	fmt.Println(msg)
	os.Exit(0)
}

func TestClientUrl(clientPort string) string {
	port := flag.String("port", clientPort, "port on localhost to check")
	flag.Parse()

	resp, err := http.Get("http://127.0.0.1:" + *port + "/upcheck")
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("dead in client")
		return "1"
	}
	return "0"
}

func TestNodePort(nodePort string) string {
	conn, err := telnet.DialTo("localhost:" + nodePort)
	if err != nil {
		fmt.Println("NodeUrl not responding")
		return "1"
	}
	conn.Write([]byte("hello world"))
	conn.Write([]byte("\n"))
	return "0"
}
