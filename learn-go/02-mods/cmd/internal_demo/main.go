package main

import (
	"fmt"

	"example.com/learn-go/02-mods/internal/greet"
)

func main() {
	fmt.Println(greet.Hello("alice"))
}
