package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

var html []byte

// HTML을 브라우저에 송신
// 정적 HTML 반환
func handlerHtml(w http.ResponseWriter, r *http.Request) {
	// Pusher로 캐스팅 가능하면 푸시한다
	w.Header().Add("Content-Type", "text/html")
	w.Write(html)
}

//소수를 브라우저에 송신
// server-sent events 전송 핸들러
// 필요한 헤더를 설정하고, 루프 안에서 소수를 출력하고, Flush()를 호출하는 처리 반복
func handlerPrimeSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	closeNotify := w.(http.CloseNotifier).CloseNotify()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Add("Connection", "keep-alive")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	var num int64 = 1
	for id := 1; id <= 100; id++ {
		// 통신이 끊겨도 종료
		select {
		case <-closeNotify:
			fmt.Println("Connection closed form client")
			return
		default:
			// do nothing
		}
		for {
			num++
			// 확률론적으로 소수를 구한다
			if big.NewInt(num).ProbablyPrime(20) {
				fmt.Println(num)
				fmt.Fprintf(w, "data: {\"id\": %d, \"number\": %d}\n\n", id, num)
				flusher.Flush()
				time.Sleep(time.Second)
				break
			}
		}
		time.Sleep(time.Second)
	}
	//100개가 넘으면 송신 종료
	fmt.Println("Connection closed from server")
}

func main() {
	var err error
	html, err = ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", handlerHtml)
	http.HandleFunc("/prime", handlerPrimeSSE)
	fmt.Println("start http listening :18888")
	err = http.ListenAndServe(":18888", nil)

	fmt.Println(err)
}
