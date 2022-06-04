package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
)

var (
	tempHistory = []byte{}
	users       = make(map[net.Conn]string)
)

func main() {
	pengue, err := ioutil.ReadFile("assets/penguen.txt")
	if err != nil {
		log.Fatal(err)
	}
	var mut sync.Mutex
	port := "8989"
	args := os.Args[1:]
	Largs := len(args)
	if Largs != 0 {
		if validPort(args[0]) && Largs == 1 {
			port = args[0]
		} else {
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("Server is listening port: " + port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Close()
			continue
		}

		if len(users) <= 9 {
			conn.Write([]byte(pengue))
			conn.Write([]byte("\n[ENTER YOUR NAME]:"))
			go handleConnection(conn, &mut)
		} else {
			conn.Write([]byte("Server is full, cannot connect. Try again later\n"))
			log.Print("Server is full")
			conn.Close()
		}
	}
}
