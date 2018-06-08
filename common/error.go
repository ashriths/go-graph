package common

import "log"

func NoError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
