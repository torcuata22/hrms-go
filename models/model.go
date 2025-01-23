package model

type Employee struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    int     `json:"age"`
}
