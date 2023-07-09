package core

import (
	"github.com/vanilla-os/apx/settings"
)

var apx *Apx

type Apx struct {
	Cnf *settings.Config
}

func NewApx(cnf *settings.Config) *Apx {
	apx = &Apx{
		Cnf: cnf,
	}

	err := apx.EssentialChecks()
	if err != nil {
		panic(err)
	}

	return apx
}

func NewStandardApx() *Apx {
	cnf, err := settings.GetApxDefaultConfig()
	if err != nil {
		panic(err)
	}

	apx = &Apx{
		Cnf: cnf,
	}

	err = apx.EssentialChecks()
	if err != nil {
		panic(err)
	}

	return apx
}
