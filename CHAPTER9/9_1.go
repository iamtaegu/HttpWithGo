package main

import (
	"fmt"
	"net/http"
)

func main_9_1() {
	resp, err := http.Get("https://google.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Printf("Protocol Version: %s\n", resp.Proto)
}
