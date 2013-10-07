package main

import (
	"fmt"
	"github.com/elegios/badbank/protocol"
	"net"
	"time"
)

const (
	timeout = time.Second
)

var buffer []byte

func handleBlobConn(blobConn net.Conn) (<-chan []string, chan<- string) {
	buffer = make([]byte, protocol.MaxBlobSize)
	out := make(chan []string)
	in := make(chan string, 1)
	go func() {
		for {
			select {
			case lang := <-in:
				askForBlob(blobConn, lang)
			default:
				readBlob(blobConn, out)
			}
		}
	}()
	return out, in
}

func askForBlob(blobConn net.Conn, lang string) {
	m := &protocol.Message{protocol.Iam, 0, 0}
	m.SetASCII(lang)
	b, erp := m.Encode()
	d(erp)
	_, erp = blobConn.Write(b)
	d(erp)
}

func readBlob(blobConn net.Conn, out chan<- []string) {
	blobConn.SetReadDeadline(time.Now().Add(timeout))
	n, err := blobConn.Read(buffer) //TODO, check correctness/error handling
	if err != nil {
		fmt.Print(err)
		return
	}
	strings := protocol.DecodeBlob(buffer[:n])
	if err != nil {
		fmt.Print("died2")
		return //TODO: print some error, the blob package was malformed
	}
	out <- strings
}
