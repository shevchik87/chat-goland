package main

import (
	"fmt"
)

func main() {
	a := App{}
	a.Initialize("postgres", "postgres", "chat")

	a.Run(":8080")

	fmt.Print(a)
}