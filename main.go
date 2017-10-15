package main

import ()

func main() {
	a := App{}
	a.Initialize("postgres", "postgres", "chat")

	a.Run(":8080")
}
