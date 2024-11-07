package main

import (
	"fmt"
	"net"
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
		resp := NewResp(conn)
		val, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		_ = val

		writer := NewWriter(conn)
		writer.Write(Value{typ: TYP_STRING, str: "OK"})
	}
}
