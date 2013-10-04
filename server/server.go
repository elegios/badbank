package main

import (
	"./database" // TODO kan vara i main. men hur?
	"./protocol"
	"fmt"
	"net"
	"os"
	"os/signal"
)

var db database.Database

func main() {
	db = database.SampleDatabase()
	go server(update, ":1338")
	go server(query, ":1337")
	waitforinterrupt() // select{}
}

func server(yield func(net.Conn), port string) {
	ln, err := net.Listen("tcp", port)
	derp(err)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Errorf("Error on client connect on query thread. %v\n", err)
			continue
		}
		go func(conn net.Conn) {
			yield(conn)
			conn.Close()
			// <3
		}(conn)
	}
}

func update(conn net.Conn) {
	message := getMessage(conn)
	if message.Opcode == protocol.Iam {
		conn.Write([]byte("TODO Cool string with updates encoded awesome!"))
	}
}

func query(conn net.Conn) {
	message, err := getMessage(conn)
	if message.Opcode != protocol.Login || err != nil {
		// send must come first
		return
	}
	account, err := db.Loggon(message.Big, message.Special)
	if err != nil {
		// Send login error
		return
	}
	defer account.Logoff()
	// Until closed connection or non change message.
	for {
		message, err := getMessage(conn)
		if err != nil || message.Opcode != protocol.Change {
			return
		}
		db.Change(message.Special, message.Big)
	}
}

func getMessage(conn net.Conn) (Message, error) {
	raw := make([]byte, 10)
	_, err := conn.Read(raw)
	var message Message
	message, err = protocol.Decode(raw)
	return message, err
}

func waitforinterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func derp(err error) {
	if err != nil {
		panic(err)
	}
}
