package main

import (
	"io/ioutil"
	"log"
	"net/http"	
)

func main () {
	
	code_6_1 ();
	
}

func code_6_1() {
	resp, err := http.Get("https://example.com")
	if err != nil {
		panic(err)
	}
	// 이 스코프를 벗어난 곳에서 반드시 닫는다
	// 함수가 종료됐을때 호출되는 후처리 연산자
	defer resp.Body.Close()
	// ioutil.ReadAll로 서버 응답을 끝까지 일괄적으로 읽는다
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(resp.Body)
	log.Println(body)
}