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
	"fmt"
	"net"
)

func CheckConnection(host, port string) bool {
	testDomain := fmt.Sprintf("%s:%s", host, port)
	_, err := net.Dial("tcp", testDomain)

	return err == nil
}
