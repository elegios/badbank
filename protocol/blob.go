package protocol

import (
	"bytes"
	"encoding/gob"
	"os"
	"sync"
)

type Blob struct {
	rwmutex sync.RWMutex
	blob    []string
}

const (
	delimiter   byte = 0
	MaxBlobSize      = 1024 * 5
)

const (
	LOGIN_INTRO = iota
	LOGIN_CARD_NUMBER
	LOGIN_PIN_CODE
	LOGIN_SUCCESS
	LOGIN_FAIL
	BALANCE
	MENU_BANNER
	MENU_BALANCE
	MENU_DEPOSIT
	MENU_WITHDRAW
	MENU_CHANGE_LANGUAGE
	MENU_QUIT
	CHANGE_AMOUNT
	DEPOSIT_CODE
	DEPOSIT_FAIL
	CHANGE_LANGUAGE_QUESTION
	LANGUAGE_WILL_CHANGE
)

func LoadLangBlob(filename string) (b *Blob) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	strings := make([]string, 0)

	dec := gob.NewDecoder(file)
	err = dec.Decode(strings) //May or may not work the way we want it to
	if err != nil {
		return
	}

	return &Blob{blob: strings}
}

func (b *Blob) SaveLangBlob(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	enc.Encode(b.blob)
}

func (b *Blob) Get(id int) string {
	b.rwmutex.RLock()
	defer b.rwmutex.RUnlock()
	return b.blob[id]
}

func (b *Blob) Set(values []string) {
	b.rwmutex.Lock()
	defer b.rwmutex.Unlock()
	b.blob = values
}

func (b *Blob) Encode() []byte {
	b.rwmutex.RLock()
	defer b.rwmutex.RUnlock()

	buf := new(bytes.Buffer)
	for _, s := range b.blob {
		buf.WriteString(s)
		buf.WriteByte(delimiter)
	}
	return buf.Bytes()
}

func DecodeBlob(b []byte) (blob []string) {
	blob = []string{}
	buf := bytes.NewBuffer(b)
	s, err := buf.ReadString(delimiter)
	for err == nil {
		blob = append(blob, s[:len(s)-1])
		s, err = buf.ReadString(delimiter)
	}
	return
}
