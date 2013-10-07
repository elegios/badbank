package main

import (
	"../protocol"
	"./database"
	"fmt"
	"net"
)

var db database.Database

func main() {
	db = database.SampleDatabase()
	go server(update, ":1338")
	go server(query, ":1337")
	select {}
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
	message, err := getMessage(conn)
	if message.Opcode == protocol.Iam || err != nil {
		sampleblob := new(protocol.Blob)
		sampleblob.Set([]string{
			"LOGIN_INTRO\n",
			"LOGIN_CARD_NUMBER\n",
			"LOGIN_PIN_CODE\n",
			"LOGIN_SUCCESS\n",
			"LOGIN_FAIL\n",
			"BALANCE\n",
			"MENU_BANNER\n",
			"MENU_BALANCE\n",
			"MENU_DEPOSIT\n",
			"MENU_WITHDRAW\n",
			"MENU_CHANGE_LANGUAGE\n",
			"MENU_QUIT\n",
			"CHANGE_AMOUNT\n",
			"DEPOSIT_CODE\n",
			"DEPOSIT_FAIL\n",
			"CHANGE_LANGUAGE_QUESTION\n",
			"LANGUAGE_WILL_CHANGE\n",
		})
		conn.Write(sampleblob.Encode())
	}
}

func query(conn net.Conn) {
	message, err := getMessage(conn)
	if message.Opcode != protocol.Login || err != nil {
		sendMessage(conn, protocol.Message{protocol.Info, protocol.Error, 0})
		return
	}
	account, err := db.Loggon(message.Big, message.Special)
	if err != nil {
		sendMessage(conn, protocol.Message{protocol.Info, protocol.Badlogon, 0})
		return
	}
	sendMessage(conn, protocol.Message{protocol.Info, protocol.Logon, 0})
	defer account.Logoff()
	for {
		// Loop until closed connection or non change message.
		message, err := getMessage(conn)
		if err != nil || message.Opcode != protocol.Change {
			sendMessage(conn, protocol.Message{protocol.Info, protocol.Error, 0})
			return
		}
		saldo, err := account.Change(message.Special, message.Big)
		if err != nil {
			sendMessage(conn, protocol.Message{protocol.Info, protocol.VercodeErr, 0})
			continue
		}
		sendMessage(conn, protocol.Message{protocol.Info, protocol.Ok, saldo})
	}
}

func getMessage(conn net.Conn) (*protocol.Message, error) {
	raw := make([]byte, 10)
	conn.Read(raw)
	return protocol.DecodeMessage(raw)
}

func sendMessage(conn net.Conn, message protocol.Message) (bytes []byte, err error) {
	bytes, err = message.Encode()
	_, err = conn.Write(bytes)
	return
}

func derp(err error) {
	if err != nil {
		panic(err)
	}
}
