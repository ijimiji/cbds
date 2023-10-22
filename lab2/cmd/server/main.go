package main

import (
	"lab2/pkg/server"
	"log"
)

func main() {
	log.Fatalln(server.New().Serve(":8081"))
}
