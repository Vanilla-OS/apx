package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Pietro di Caprio <pietro@fabricators.ltd>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"fmt"

	"github.com/vanilla-os/apx/v2/settings"
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
		// localisation features aren't available at this stage, so this error can't be translated
		fmt.Println("ERROR: Unable to find apx configuration files")
		return nil
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
		// localisation features aren't available at this stage, so this error can't be translated
		fmt.Println("ERROR: Unable to find apx configuration files")
		return nil
	}
	return apx
}
