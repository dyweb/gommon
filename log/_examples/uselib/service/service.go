package service

import (
	"errors"

	"github.com/dyweb/gommon/log/_examples/uselib/storage"
)

type Auth struct {
	s storage.Driver
}

func NewAuth(driver storage.Driver) *Auth {
	return &Auth{
		s: driver,
	}
}

func (a *Auth) Check(user string, password string) error {
	//log2.NewIdentityFromCallerOld(0)
	p, err := a.s.Get(user)
	if err != nil {
		return err
	}
	if p == password {
		return nil
	}
	return errors.New("invalid password")
}
