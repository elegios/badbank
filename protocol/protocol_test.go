package protocol

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestDecodeOpcodeIam(t *testing.T) {
	raw := make([]byte, 10)
	message, _ := Decode(raw)
	if message.Opcode != Iam {
		t.Fail()
	}
}

func TestDecodeOpcodeInfo(t *testing.T) {
	raw := make([]byte, 10)
	raw[0] = 0xFF
	message, _ := Decode(raw)
	if message.Opcode != Info {
		t.Log(message.Opcode)
		t.Fail()
	}
}

func TestDecodeSpecial(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, uint16(1))
	binary.Write(buf, binary.BigEndian, int64(1))
	message, _ := Decode(buf.Bytes())
	if message.Special != 1 {
		t.Fail()
	}
}

func TestDecodeSpecialOverflow(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, uint16(0xFFFF))
	binary.Write(buf, binary.BigEndian, int64(1))
	message, _ := Decode(buf.Bytes())
	t.Log(message)
	if message.Special != 0xFFFF/4 {
		t.Fail()
	}
}

func TestDecodeBigPositiv(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0, 0})
	binary.Write(buf, binary.BigEndian, int64(1))
	message, _ := Decode(buf.Bytes())
	t.Log(message)
	if message.Big != 1 {
		t.Fail()
	}
}

func TestDecodeBigNegative(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0, 0})
	binary.Write(buf, binary.BigEndian, int64(-1))
	message, _ := Decode(buf.Bytes())
	if message.Big != -1 {
		t.Fail()
	}
}

func TestEncodeOpcode(t *testing.T) {
	message := Message{1, 0, 0}
	raw := Encode(message)
	opcode := int8(raw[0]) // inte helt sant men för 1 så
	t.Log(raw, opcode)
	if opcode != 1 {
		t.Fail()
	}
}

func TestEncodeSpecial(t *testing.T) {
	message := Message{0, 10, 0}
	raw := Encode(message)
	special := uint8(raw[1]) // inte helt sant men för 10 så
	t.Log(raw, special)
	if special != 10 {
		t.Fail()
	}
}

func TestEncodeBig(t *testing.T) {
	message := Message{0, 0, -10}
	raw := Encode(message)
	var big int64
	buf := bytes.NewBuffer(raw[2:])
	binary.Read(buf, binary.BigEndian, &big)
	t.Log(raw, big)
	if big != -10 {
		t.Fail()
	}
}
