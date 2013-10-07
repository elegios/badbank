package main

import (
	"fmt"
	"github.com/elegios/badbank/protocol"
	"net"
	"os"
)

var (
	blob = protocol.LoadLangBlob(BLOB_FILE)
)

const (
	MAIN_PORT    = 9090
	BLOB_PORT    = 9091
	BLOB_FILE    = "blobfile"
	DEFAULT_BLOB = "en"
)

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	mainConn, erp := net.Dial("tcp", fmt.Sprintf("%s:%d", os.Args[1], MAIN_PORT))
	d(erp)
	blobConn, erp := net.Dial("tcp", fmt.Sprintf("%s:%d", os.Args[1], BLOB_PORT))
	d(erp)

	blobChan, langChan := handleBlobConn(blobConn)
	if blob == nil {
		langChan <- DEFAULT_BLOB
		blob = new(protocol.Blob)
		blob.Set(<-blobChan)
	}

	if !login(mainConn) {
		os.Exit(1)
	}

	go func() {
		for {
			loggedInInteraction(mainConn, langChan)
		}
	}()

	for b := range blobChan {
		blob.Set(b)
	}
}

func d(err error) {
	if err != nil {
		panic(err)
	}
}
