package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Listening on port 6379")
	server()
}

func server() {
	//Create server
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	//Listen for connections
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	//Main loop for server
	for {
		buf := make([]byte, 1024)

		//read messages from client
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from the client: ", err.Error())
			os.Exit(1)
		}

		//ignore request and send back OK
		conn.Write([]byte("+OK\r\n"))
	}
}
