package main

import (
	"github.com/flemming-petersen/promptbattle/server"
)

func main() {
	app := server.NewServer()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
