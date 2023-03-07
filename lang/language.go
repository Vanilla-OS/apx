package lang

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around apt to make it works inside a container
	from outside, directly on the host.
*/

import (
	"fmt"
	"log"
	"os"
)

func GetText(lang string, name string) string {
	content, err := os.ReadFile(fmt.Sprintf("lang/%s/%s", lang, name))
	if err != nil {
		log.Println("FATAL! Can't get translation file!")
		log.Fatal(err)
	}
	return string(content)
}
