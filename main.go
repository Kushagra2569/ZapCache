package main

import (
	"fmt"
	"net"
	"strings"
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

	//Initialise NewAof for data persistence
	aof, err := NewAof("database.aof")
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

		if val.typ != TYP_ARRAY {
			fmt.Println("Invalid Request, Expected array")
			continue
		}

		if len(val.array) == 0 {
			fmt.Println("Invalid Request, Expected array length greater than 0")
			continue
		}
		command := strings.ToUpper(val.array[0].bulk)
		args := val.array[1:]

		writer := NewWriter(conn)
		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid Command")
			writer.Write(Value{typ: TYP_STRING, str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(val)
		}

		result := handler(args)
		writer.Write(result)
	}
}
