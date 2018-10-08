package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var db map[string]string

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}
	defer lis.Close()
	db = make(map[string]string)
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
	fmt.Printf("Accepted Connection From %v\n", addr)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		slice := strings.Fields(line)
		if len(slice) < 1 {
			continue
		}
		handleCommand(conn, slice)
	}

	defer conn.Close()
	fmt.Printf("%v connection closed\n", addr)
}

func handleCommand(conn net.Conn, slice []string) {
	switch slice[0] {
	case "GET":
		if len(slice) != 2 {
			fmt.Fprintln(conn, "INVALID ARGS")
			break
		}
		k := slice[1]
		val := db[k]
		fmt.Fprintln(conn, val)
	case "SET":
		if len(slice) != 3 {
			fmt.Fprintln(conn, "INVALID ARGS")
			break
		}
		k := slice[1]
		v := slice[2]
		db[k] = v
	case "DELETE":
		if len(slice) != 2 {
			fmt.Fprintln(conn, "INVALID ARGS")
			break
		}
		k := slice[1]
		delete(db, k)
		fmt.Fprintf(conn, "deleted %s\n", k)
	default:
		fmt.Fprintln(conn, "INVALID COMMAND "+slice[0]+"\r\n")
	}
}
