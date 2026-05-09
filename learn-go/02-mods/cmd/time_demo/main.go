package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println("now:", now.Format(time.RFC3339))

	expireAt := now.Add(30 * time.Minute)
	fmt.Println("expireAt:", expireAt.Format(time.RFC3339))
}
