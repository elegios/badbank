package database

import (
	"errors"
)

type Account struct {
	saldo    int64
	pin      uint16
	loggedon chan bool // Empty = not logged on
}

type AccountNumber int64
type Database struct {
	data map[AccountNumber]*Account
}

func SampleDatabase() Database {
	return Database{map[AccountNumber]*Account{
		42:            &Account{10, 1234, make(chan bool, 1)},
		1234567890123: &Account{0x7FFFFFFFFFFFFFFF, 4444, make(chan bool, 1)},
		46:            &Account{1E12, 7777, make(chan bool, 1)},
	}}
}

func (account *Account) verify(code int8) bool {
	return code%2 != 0
}

func (account *Account) Change(vercode int8, diff int64) (int64, error) {
	if diff < 0 && account.verify(vercode) == false {
		return 0, errors.New("Wrong vercode!")
	}
	account.saldo += diff
	return account.saldo, nil
}

func (account *Account) Logoff() {
	<-account.loggedon
}

func (database *Database) Loggon(nr AccountNumber, pin uint16) (*Account, error) {
	a := database.data[nr]
	if a.pin != pin {
		return nil, errors.New("Wrong pin!")
	}
	select {
	case a.loggedon <- true:
		return a, nil
	default:
		return nil, errors.New("Already logged on.")
	}
}
