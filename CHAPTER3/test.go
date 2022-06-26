package main

import (
	"log"
	"os"
)

func main() {
	path, _ := os.Getwd();

	log.Println(path)
}
