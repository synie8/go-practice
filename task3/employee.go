package task3

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type QueryLogger struct{}

func (ql *QueryLogger) LogQuery(query string, args []interface{}) {
	log.Printf("Query: %s\nArgs: %v\n", query, args)
}

type Employee struct {
	ID         uint64 `gorm:"primary_key"`
	Name       string
	Department string
	Salary     float32
}

var db *sqlx.DB

func (e *Employee) InitDB() (db *sqlx.DB, err error) {
	dsn := "root:wtms.123@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return db, err
}
func (e *Employee) Create() {
	gormDb := ConnectDb()
	gormDb.AutoMigrate(&Employee{})

}
func (e *Employee) CreateData1() {
	gormDb := ConnectDb()
	var employees = []Employee{{Name: "张1", Department: "技术部", Salary: 80000}, {Name: "张2", Department: "技术部", Salary: 70000}, {Name: "张3", Department: "销售部", Salary: 70000}}

	gormDb.Debug().Create(&employees)
}
func (e *Employee) Query1(department string) []Employee {
	var result []Employee
	err := db.Select(&result, "select id,name,department,salary from employees where department=?", department)
	if err != nil {
		fmt.Printf("Query1 failed, err:%v\n", err)
	}
	defer db.Close()
	return result
}
func (e *Employee) Query2() Employee {
	var result Employee
	sqlStr := `select a.id,a.name,a.department,a.salary 
			from employees  a
			join (select max(salary) as salary,name from employees  group by name)  b on a.name=b.name and a.salary=b.salary
			
			`
	err := db.Get(&result, sqlStr)
	if err != nil {
		fmt.Printf("Query2 failed, err:%v\n", err)
	}
	defer db.Close()
	return result
}
