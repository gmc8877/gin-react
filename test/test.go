package main

import (

	"fmt"

	"github.com/satori/go.uuid"
)

func main() {
	
	
    
  u2 := uuid.NewV4()
  res := u2.String()
  fmt.Println(res)
}