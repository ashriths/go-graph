package system

import "log"

var Logging bool

func Logln(v ...interface{}) {
	if Logging {
		log.Println(v)
	}
}

func Logf(format string, v ...interface{}) {
	if Logging {
		log.Printf(format, v...)
	}
}
