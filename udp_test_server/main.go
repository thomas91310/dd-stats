package main

import (
	"fmt"
	"net"
)

func main() {
	udpEndpoint, err := net.ResolveUDPAddr("udp", ":8125")
	if err != nil {
		panic(err)
	}

	connection, err := net.ListenUDP("udp", udpEndpoint)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	buf := make([]byte, 1024)

	for {
		n, _, err := connection.ReadFromUDP(buf)
		fmt.Println("Received: ", string(buf[0:n]))

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
