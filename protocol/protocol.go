package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Opcode int8

const (
	Iam = iota
	Login
	Change
	Info
)

type Message struct {
	Opcode  Opcode
	Special uint16
	Big     int64
}

type Infocode uint16

func Decode(raw []byte) (Message, error) {
	var message Message
	if len(raw) != 10 {
		return message, errors.New("Message length must be 10 bytes!")
	}
	message.Opcode = Opcode(raw[0] >> 6)
	raw[0] = raw[0] & 0x3F
	buf := bytes.NewBuffer(raw)
	err := binary.Read(buf, binary.BigEndian, &message.Special)
	err = binary.Read(buf, binary.BigEndian, &message.Big)
	if err != nil {
		return message, errors.New("Message format is wierd.")
	}
	return message, nil
}

func Encode(message Message) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, message.Special)
	binary.Write(buf, binary.BigEndian, message.Big)
	code := buf.Bytes()
	buf.Reset()
	binary.Write(buf, binary.BigEndian, message.Opcode)
	code[0] = code[0] | buf.Bytes()[0]
	return code
}
