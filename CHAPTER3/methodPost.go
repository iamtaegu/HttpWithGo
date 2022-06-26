package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func fileSend() {
	// os.File 오브젝트는 io.Reader I/F 만족
	file, err := os.Open("server.go")
	if err != nil {
		panic(err)
	}

	/*values := url.Values{
		"test": {"value"},
	}*/
	//resp, err := http.PostForm("http://localhost:18888", values)
	// io.Reader 형식으로 전달
	// 텍스트 전달에는 문자열이 아닌 io.Reader 형식으로 인터페이스화 해야함
	resp, err := http.Post("http://localhost:18888", "text/plain", file)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
}

func stringSend() {
	reader := strings.NewReader("텍스트") // io.Reader 인터페이스화
	resp, err := http.Post("http://localhost:18888", "text/plain", reader)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
}

func main() {
	fileSend()
	stringSend()
}
