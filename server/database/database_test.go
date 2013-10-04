package database

import (
	"testing"
)

func TestVerify(t *testing.T) {
	a := new(Account)
	if a.verify(9) != true {
		t.Fail()
	}
	if a.verify(10) != false {
		t.Fail()
	}
}

func TestChangeWithdwral(t *testing.T) {
	a := Account{100, 0, make(chan bool, 1)}
	saldo, _ := a.Change(9, -100)
	if saldo != 0 {
		t.Fail()
	}
	saldo, _ = a.Change(11, -100)
	if saldo != -100 {
		t.Fail()
	}
}

func TestChangeVercode(t *testing.T) {
	a := Account{100, 0, make(chan bool, 1)}
	_, err := a.Change(10, -100)
	if err == nil {
		t.Fail()
	}
}

func TestChangeDeposit(t *testing.T) {
	a := Account{100, 0, make(chan bool, 1)}
	_, err := a.Change(10, 100)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestLoggon(t *testing.T) {
	d := SampleDatabase()
	a, err := d.Loggon(42, 1234)
	if !(err == nil && a.saldo == 10) {
		t.Log(a)
		t.Fail()
	}
}

func TestLoggonWrongPin(t *testing.T) {
	d := SampleDatabase()
	_, err := d.Loggon(42, 0000)
	if err == nil {
		t.Fail()
	}
}

func TestLoggonAlreadyOn(t *testing.T) {
	d := SampleDatabase()
	d.Loggon(42, 1234)
	_, err := d.Loggon(42, 1234)
	if err == nil {
		t.Log(err)
		t.Fail()
	}
}
