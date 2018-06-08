package storage

import "log"

var Logging bool

func logln(v ...interface{})  {
	if Logging{
		log.Println(v)
	}
}

func logf(format string, v ...interface{})  {
	if Logging{
		log.Printf(format, v...)
	}
}
