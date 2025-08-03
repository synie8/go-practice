package task2

import "fmt"

type Person struct {
	Name string
	Age  int
}
type Employee struct {
	Person
	EmployeeID string
}

func (employee *Employee) PrintInfo() {
	fmt.Println("employee:", employee)
}
