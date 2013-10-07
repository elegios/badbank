package main

import (
	"fmt"
	"github.com/elegios/badbank/protocol"
	"net"
	"os"
)

func login(mainConn net.Conn) (success bool) {
	loginMessage := new(protocol.Message)
	loginMessage.Opcode = protocol.Login

	fmt.Print(blob.Get(protocol.LOGIN_INTRO))

	fmt.Print(blob.Get(protocol.LOGIN_CARD_NUMBER))
	fmt.Scanf("%d", &loginMessage.Big)

	fmt.Print(blob.Get(protocol.LOGIN_PIN_CODE))
	fmt.Scanf("%d", &loginMessage.Special)

	sendMessage(loginMessage, mainConn)
	loginMessage = readMessage(mainConn)
	success = loginMessage.IsLoginSuccess()

	if success {
		fmt.Print(blob.Get(protocol.LOGIN_SUCCESS))
	} else {
		fmt.Print(blob.Get(protocol.LOGIN_FAIL))
	}

	fmt.Printf("%s %d\n", blob.Get(protocol.BALANCE), loginMessage.Big)

	return
}

func loggedInInteraction(mainConn net.Conn, langChan chan<- string) {
	fmt.Print(blob.Get(protocol.MENU_BANNER))
	fmt.Printf("1: %s", blob.Get(protocol.MENU_BALANCE))
	fmt.Printf("2: %s", blob.Get(protocol.MENU_DEPOSIT))
	fmt.Printf("3: %s", blob.Get(protocol.MENU_WITHDRAW))
	fmt.Printf("4: %s", blob.Get(protocol.MENU_CHANGE_LANGUAGE))
	fmt.Printf("5: %s", blob.Get(protocol.MENU_QUIT))

	var choice uint
	fmt.Scanf("%d", &choice)

	m := &protocol.Message{protocol.Change, 0, 0}

	switch choice {
	case 1:

	case 2:
		fmt.Print(blob.Get(protocol.CHANGE_AMOUNT))
		fmt.Scanf("%d", &m.Big)

	case 3:
		fmt.Print(blob.Get(protocol.CHANGE_AMOUNT))
		fmt.Scanf("%d", &m.Big)
		m.Big *= -1
		fmt.Print(blob.Get(protocol.DEPOSIT_CODE))
		fmt.Scanf("%d", &m.Special)

	case 4:
		var lang string
		fmt.Print(blob.Get(protocol.CHANGE_LANGUAGE_QUESTION))
		fmt.Scanf("%s", &lang)
		langChan <- lang
		fmt.Print(blob.Get(protocol.LANGUAGE_WILL_CHANGE))
		return

	case 5:
		os.Exit(0)

	default:
		return
	}

	sendMessage(m, mainConn)
	m = readMessage(mainConn)
	if !m.IsChangeSuccess() {
		fmt.Print(blob.Get(protocol.DEPOSIT_FAIL))
	}
	fmt.Printf("%s %d\n", blob.Get(protocol.BALANCE), m.Big)
}

func sendMessage(message *protocol.Message, conn net.Conn) {
	b, _ := message.Encode()
	conn.Write(b)
}

func readMessage(conn net.Conn) (m *protocol.Message) {
	b := make([]byte, 10)
	n, _ := conn.Read(b) //TODO: check error handling
	if n != 10 {
		panic(n)
	}
	m, erp := protocol.DecodeMessage(b)
	d(erp)
	return
}
