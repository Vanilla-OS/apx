package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around apt to make it works inside a container
	from outside, directly on the host.
*/

import (
	"net/http"
)

func CheckConnection() bool {
	_, err := http.Get("https://google.com") // TODO: use a better way to check connection
	if err != nil {
		return false
	}
	return true
}
