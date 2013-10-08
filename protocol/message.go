package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	Change = iota
	Iam
	Login
	Info
)

//Bitmasks for infomessages
const (
	login uint16 = 1 << iota
	change
)

const (
	opshift        uint   = 14
	charactershift uint   = 7
	charactermask  uint16 = (1 << charactershift) - 1
)

type Message struct {
	Opcode  uint8
	Special uint16
	Big     int64
}

func DecodeMessage(raw []byte) (message *Message, err error) {
	if len(raw) != 10 {
		return nil, errors.New("Message length must be 10 bytes!")
	}
	message = new(Message)
	buf := bytes.NewBuffer(raw)
	binary.Read(buf, binary.BigEndian, &message.Special) //cannot error, buf is long enough
	binary.Read(buf, binary.BigEndian, &message.Big)
	message.Opcode = uint8((message.Special & (3 << opshift)) >> opshift)
	message.Special &= ^(uint16(3) << opshift)
	return
}

func (m *Message) Encode() ([]byte, error) {
	if m.Opcode >= 1<<2 || m.Special >= 1<<14 {
		return nil, errors.New("The opcode is more than 2 bits or the special is more than 14 bits")
	}
	buf := bytes.NewBuffer(make([]byte, 0, 10))
	var special uint16 = (uint16(m.Opcode) << opshift) | m.Special
	binary.Write(buf, binary.BigEndian, special)
	binary.Write(buf, binary.BigEndian, m.Big)
	return buf.Bytes(), nil
}

func (m *Message) GetASCII() string {
	return string(
		[]byte{
			byte((m.Special & (charactermask << charactershift)) >> charactershift),
			byte(m.Special & charactermask),
		},
	)
}

func (m *Message) SetASCII(str string) {
	b := []byte(str)
	m.Special = (uint16(b[0]) << charactershift) | uint16(b[1])
}

func (m *Message) IsLoginSuccess() bool {
	return m.Opcode == Info && m.Special&login == login
}

func (m *Message) IsChangeSuccess() bool {
	return m.Opcode == Info && m.Special&change == change
}

func (m *Message) SetLoginSuccess(success bool) {
	if success {
		m.Special |= login
	} else {
		m.Special ^= login
	}
}

func (m *Message) SetChangeSuccess(success bool) {
	if success {
		m.Special |= change
	} else {
		m.Special ^= change
	}
}
