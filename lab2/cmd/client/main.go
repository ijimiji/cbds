package main

import (
	"lab2/pkg/client"
	"log"
)

func main() {
	log.Fatalln(client.New().Serve(":8080"))
}
