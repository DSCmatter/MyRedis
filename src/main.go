// Server Setup

package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close() // // close connection once finished

	// An infinite loop and receive commands from clients and respond to them.
	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		// modify to receive ping command from handler.go
		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		// Example value object upon setting a name:

		// Value{
		// 	typ: "array",
		// 	array: []Value{
		// 		Value{typ: "bulk", bulk: "SET"},
		// 		Value{typ: "bulk", bulk: "name"},
		// 		Value{typ: "bulk", bulk: "Ahmed"},
		// 	},
		// }

		// The code above will make the command and args look like this:

		// Perform validations to make sure command is array and not empty.
		command := strings.ToUpper(value.array[0].bulk) // SET
		args := value.array[1:]

		// usage of the Writer
		writer := NewWriter(conn)

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
