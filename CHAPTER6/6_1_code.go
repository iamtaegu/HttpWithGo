package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func temp() {
	resp, err := http.Get("https://example.com")
	if err != nil {
		panic(err)
	}
	// 이 스코프를 벗어난 곳에서 반드시 닫는다
	defer resp.Body.Close()
	// ioutil.ReadAll로 서버 응답을 끝까지 일괄적으로 읽는다
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(resp.Body)
	log.Println(body)
}
