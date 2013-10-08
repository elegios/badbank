package main

import (
	"fmt"
	"github.com/elegios/badbank/protocol"
	"github.com/elegios/badbank/server/database"
	"io"
	"net"
)

var db database.Database

func main() {
	db = database.SampleDatabase()
	go server(update, ":1338")
	go server(query, ":1337")
	go gui()
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

var langs map[string]*protocol.Blob

func init() {
	langs = make(map[string]*protocol.Blob)
	langs["bg"] = new(protocol.Blob)
	langs["bg"].Set([]string{
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
}

func gui() {
	fmt.Println("Welcome to Badbank Server!")
	for {
		fmt.Print("Change or create language (en): ")
		var lang string
		fmt.Scanf("%s", &lang)
		if len(lang) == 0 {
			lang = "en"
		}
		blob := langs[lang]
		var fg []string
		if blob == nil {
			// Create lang
			blob = new(protocol.Blob)
			fg = make([]string, protocol.BLOBLENGTH)
		} else {
			fg = blob.GetSlice()
		}
		bg := langs["bg"].GetSlice()
		for i := range fg {
			fmt.Printf("%s = ", bg[i])
			var scan string
			fmt.Scanf("%s", &scan)
			fg[i] = scan + "\n"
		}
		blob.Set(fg)
		langs[lang] = blob
	}
}

func update(conn net.Conn) {
	lang := "bg"
	for {
		message, err := getMessage(conn)
		if err == io.EOF {
			return
		}
		if message.Opcode == protocol.Iam || err != nil {
			lang = message.GetASCII()
			blob := langs[lang]
			if blob == nil {
				lang = "bg"
				blob = langs[lang]
			}
			conn.Write(blob.Encode())
		}
	}
}

func query(conn net.Conn) {
	message, err := getMessage(conn)
	m := &protocol.Message{protocol.Info, 0, 0}
	m.SetLoginSuccess(false)
	if message.Opcode != protocol.Login || err != nil {
		sendMessage(conn, m)
		return
	}
	account, err := db.Loggon(message.Big, message.Special)
	if err != nil {
		sendMessage(conn, m)
		return
	}
	m.Big, _ = account.Change(0, 0)
	m.SetLoginSuccess(true)
	sendMessage(conn, m)
	defer account.Logoff()
	for {
		// Loop until closed connection or non change message.
		message, err := getMessage(conn)
		m := &protocol.Message{protocol.Info, 0, 0}
		m.SetChangeSuccess(false)
		if err != nil || message.Opcode != protocol.Change {
			m.Big, _ = account.Change(0, 0)
			sendMessage(conn, m)
			return
		}
		saldo, err := account.Change(message.Special, message.Big)
		m.Big = saldo
		if err != nil {
			sendMessage(conn, m)
			continue
		}
		m.SetChangeSuccess(true)
		sendMessage(conn, m)
	}
}

func getMessage(conn net.Conn) (*protocol.Message, error) {
	raw := make([]byte, 10)
	_, err := conn.Read(raw)
	if err != nil {
		return nil, err
	}
	message, err := protocol.DecodeMessage(raw)
	return message, err
}

func sendMessage(conn net.Conn, message *protocol.Message) (bytes []byte, err error) {
	bytes, err = message.Encode()
	_, err = conn.Write(bytes)
	return
}

func derp(err error) {
	if err != nil {
		panic(err)
	}
}
