package lang

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2023
	Description:
		Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
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
