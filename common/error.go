package common

import "log"

func NoError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func LogError(e error){
	if e != nil{
		log.Println(e)
	}
}
