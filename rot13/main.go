package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}
	defer lis.Close()

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println(err)
		}

		go handle(conn)
	}

}

func handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	fmt.Printf("Accepted connection from %v\n", addr)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		bs := rot13(line)
		fmt.Fprintf(conn, "%s -> %s\n", line, string(bs))

	}

	defer conn.Close()
	fmt.Printf("Client %v closed the connection\n", addr)
}

func rot13(str string) []byte {
	bs := make([]byte, len(str))
	for i, r := range strings.ToLower(str) {
		if r >= 'a' && r <= 'z' {
			c := r - 'a' + 13
			if c >= 26 {
				c = c % 26
			}
			c += 'a'
			bs[i] = byte(c)
		} else {
			bs[i] = byte(r)
		}
	}
	return bs
}
