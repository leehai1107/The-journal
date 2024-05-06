package main

import (
	"log"

	"github.com/leehai1107/The-journal/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Printf("error while execute: %s", err.Error())
	}
}
