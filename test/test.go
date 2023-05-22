package main

import (
	"DEMO01/handler"
	"DEMO01/tools"
	"time"
)

func main() {
  handler.InitDB()
  defer handler.DB.Close()
  db := handler.DB
	password_hash := tools.CodeHash("567777")
	smtp, err := db.Prepare("insert into root_info (user_name, password, submission_date)  VALUES (?,?, ?)")
	tools.CheckErr(err)
	_, err = smtp.Exec("13614758877", password_hash, time.Now().Format("2006-01-02"))
	tools.CheckErr(err)
}