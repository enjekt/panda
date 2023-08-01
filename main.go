package main

import (
	"log"
	"panda/restserver"
)

func main() {
	log.Println("Starting Panda Service...")
	restserver.NewRestServer()

}
