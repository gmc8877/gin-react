package tools

import (
	
	
	"log"
	
)



func CheckErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}







