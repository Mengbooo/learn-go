package main

import (
	"encoding/json"
	"fmt"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Animal interface {
	Speak() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "wang wang"
}

type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "miao miao"
}

func printTypes() {
	var (
		b   bool    = true
		i   int     = 42
		i64 int64   = 64
		f   float64 = 3.14
		s   string  = "hello"
		by  byte    = 'A'
		r   rune    = '中'
	)

	fmt.Println("练习 1：类型实验")
	fmt.Printf("bool: %v, type: %T\n", b, b)
	fmt.Printf("int: %v, type: %T\n", i, i)
	fmt.Printf("int64: %v, type: %T\n", i64, i64)
	fmt.Printf("float64: %v, type: %T\n", f, f)
	fmt.Printf("string: %v, type: %T\n", s, s)
	fmt.Printf("byte: %v, type: %T\n", by, by)
	fmt.Printf("rune: %v, type: %T\n", r, r)
}

func printSliceExperiment() {
	fmt.Println("\n练习 2：切片实验")

	var nums []int
	for i := 1; i <= 10; i++ {
		nums = append(nums, i)
		fmt.Printf("append %d -> value: %v, len: %d, cap: %d\n", i, nums, len(nums), cap(nums))
	}
}

func printTodoJSON() {
	fmt.Println("\n练习 3：结构体 + JSON")

	todo := Todo{
		ID:    1,
		Title: "learn go",
		Done:  false,
	}

	data, err := json.Marshal(todo)
	if err != nil {
		fmt.Println("marshal json error:", err)
		return
	}

	fmt.Println(string(data))
}

func printAnimals() {
	fmt.Println("\n练习 4：接口实验")

	animals := []Animal{
		Dog{Name: "Buddy"},
		Cat{Name: "Kitty"},
	}

	for _, animal := range animals {
		fmt.Printf("%T says: %s\n", animal, animal.Speak())
	}
}

func main() {
	printTypes()
	printSliceExperiment()
	printTodoJSON()
	printAnimals()
}
