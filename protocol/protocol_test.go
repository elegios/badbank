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

func TestDecodeBigPositive(t *testing.T) {
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
	raw, err := message.Encode()
	if err != nil {
		t.Fail()
	}
	opcode := uint8(raw[0]) // inte helt sant men för 1 så
	t.Log(raw, opcode)
	if opcode != 1<<6 {
		t.Fail()
	}
}

func TestEncodeSpecial(t *testing.T) {
	message := Message{0, 10, 0}
	raw, err := message.Encode()
	if err != nil {
		t.Fail()
	}
	special := uint8(raw[1]) // inte helt sant men för 10 så
	t.Log(raw, special)
	if special != 10 {
		t.Fail()
	}
}

func TestEncodeLargeSpecial(t *testing.T) {
	message := Message{0, 0x4000, 0}
	_, err := message.Encode()
	if err == nil {
		t.Fail()
	}
}

func TestEncodeLargeOpcode(t *testing.T) {
	message := Message{4, 0, 0}
	_, err := message.Encode()
	if err == nil {
		t.Fail()
	}
}

func TestEncodeBig(t *testing.T) {
	message := Message{0, 0, -10}
	raw, err := message.Encode()
	if err != nil {
		t.Fail()
	}
	var big int64
	buf := bytes.NewBuffer(raw[2:])
	binary.Read(buf, binary.BigEndian, &big)
	t.Log(raw, big)
	if big != -10 {
		t.Fail()
	}
}

func TestBothWays(t *testing.T) {
	messages := []Message{
		Message{3, 465, 1920347102},
		Message{1, 7293, -1234904575772},
		Message{0, 0, 12334123},
		Message{0, 0, 0},
		Message{2, 8, 12},
	}
	for _, m := range messages {
		b, err := m.Encode()
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if m2, err := Decode(b); !(err == nil && *m2 == m) {
			t.Log(err, m2)
			t.Fail()
		}
	}
}

func TestSetASCII(t *testing.T) {
	message := new(Message)
	message.SetASCII("en")
	if message.Special != 0x32EE {
		t.Log(message.Special)
		t.Fail()
	}
}

func TestGetASCII(t *testing.T) {
	message := &Message{0, 0x39F6, 0}
	if message.GetASCII() != "sv" {
		t.Fail()
	}
}

func TestASCIIBothWays(t *testing.T) {
	list := []string{"sv", "en", "fr", "br", "lo"}
	m := new(Message)
	for _, s := range list {
		m.SetASCII(s)
		if s != m.GetASCII() {
			t.Fail()
		}
	}
}
