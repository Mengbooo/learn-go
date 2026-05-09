package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
}

func main() {
	user := User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
		Admin: true,
	}

	data, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	fmt.Println("marshal:", string(data))

	raw := `{"id":2,"name":"Bob","email":"bob@example.com","admin":false}`

	var decoded User
	if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
		panic(err)
	}

	fmt.Printf("unmarshal: %+v\n", decoded)
}
