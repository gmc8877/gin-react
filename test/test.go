package main

import (
	"DEMO01/handler"
	"fmt"
	
)

func main() {
  handler.InitDB()
  db := handler.DB
  defer db.Close()
  var pubdate string
  err := db.QueryRow("select pubdate from channel_1 where title=1").Scan(&pubdate)
  if err != nil {
    fmt.Println("err1", err)
  }
  fmt.Println(pubdate[:10])
 
}