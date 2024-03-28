package main

import (
	"github.com/flemming-petersen/promptbattle/server"
)

func main() {
	app := server.NewServer()

	app.Run()
}
