package task3

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Student struct {
	gorm.Model
	Name  string
	Age   int
	Grade string
}

func ConnectDb() *gorm.DB {
	dsn := "root:wtms.123@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("操作失败: %v", err)
	}
	return db
}

func InsrtStd() {
	db := ConnectDb()
	db.AutoMigrate(&Student{})
	std := Student{Name: "bob", Age: 13, Grade: "一年级"}
	db.Create(&std)
}

func QueryGt18() {
	db := ConnectDb()
	var stds []Student
	db.Debug().Where("age > ?", 18).Find(&stds)
	fmt.Println("查询结果", stds)
}

func UpdateTo4() {
	db := ConnectDb()
	db.Debug().Model(&Student{}).Where("name = ?", "张三").Update("grade", "四年级")
}

func Deletelt15() {
	db := ConnectDb()
	db.Debug().Where("age < ?", 15).Delete(&Student{})
}
func (std *Student) Delete() (err error) {
	db := ConnectDb()
	return db.Debug().Unscoped().Where("age < ?", 15).Delete(std).Error
}
