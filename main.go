package main

import (
	"go-practice/task3"
)

func main() {
	db := task3.ConnectDb()
	var comment task3.Comment
	db.Debug().Where("id=?", 4).First(&comment)
	db.Debug().Delete(&comment)

}
