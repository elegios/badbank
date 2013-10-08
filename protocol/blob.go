package protocol

import (
	"bytes"
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
	BLOBLENGTH
)

func (b *Blob) GetSlice() []string {
	b.rwmutex.RLock()
	defer b.rwmutex.RUnlock()
	slice := make([]string, len(b.blob))
	copy(slice, b.blob)
	return slice
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
