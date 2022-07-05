package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/k0kubun/pp"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>hello</body></html>\n")
}

func handlerDigest(w http.ResponseWriter, r *http.Request) {
	pp.Printf("URL: %s\n", r.URL.String())
	pp.Printf("Query: %s\n", r.URL.Query())
	pp.Printf("Proto: %s\n", r.Proto)
	pp.Printf("Method: %s\n", r.Method)
	pp.Printf("Header: %s\n", r.Header)
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("--body--\n%s\n", string(body))
	if _, ok := r.Header["Authorization"]; !ok {
		fmt.Printf("--if--")
		w.Header().Add(
			"WWW-Authenticate",
			`Digest releam="Secret Zone",
			nonce="TgLc25U2BQA=f510a2780473e18e6587be702c2e67fe2b04afd",
			algorithm=MD5,
			qop="auth"`)
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		fmt.Printf("--else--")
		fmt.Fprintf(w, "<html><body>secret page</body></html>\n")
	}
}

func handlerUpgrade(w http.ResponseWriter, r *http.Request) {
	// 이 엔드포인트에서는 변경 외는 받아들이지 않는다
	if r.Header.Get("Connection") != "Upgrade" || r.Header.Get("Upgrade") != "MyProtocol" {
		w.WriteHeader(400)
		return
	}
	fmt.Println("Upgrade to MyProtocol")

	// 소켓을 획득
	hjacker := w.(http.Hijacker)
	conn, readWriter, err := hjacker.Hijack()
	if err != nil {
		panic(err)
		return
	}
	defer conn.Close()

	//프로토콜이 바뀐다는 응답을 보낸다
	response := http.Response{
		StatusCode: 101,
		Header:     make(http.Header),
	}
	response.Header.Set("Upgrade", "MyProtocol")
	response.Header.Set("Connection", "Upgrade")
	response.Write(conn)

	// 오리지널 통신 시작
	for i := 1; i <= 10; i++ {
		fmt.Fprintf(readWriter, "%d\n", i)
		fmt.Println("->", i)
		readWriter.Flush()                      // Trigger "chunked" encoding and send a chunk...
		recv, err := readWriter.ReadBytes('\n') // 개행이 올때까지 읽고, 모아서 반환
		if err == io.EOF {
			break
		}
		fmt.Printf("<- %s", string(recv))
		time.Sleep(500 * time.Millisecond)
	}
}

func handlerChunkedResponse(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("expected http.ResponseWriter to be an http.Flusher")
	}
	for i := 1; i <= 10; i++ {
		fmt.Fprintf(w, "Chunk #%d\n", i)
		flusher.Flush()
		time.Sleep(500 * time.Millisecond)
	}
	flusher.Flush()
}

func main() {
	var httpServer http.Server

	/*server := &http.Server{
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequestClientCert,
			//ClientAuth: tls.RequireAndVerifyClientCert, //가장 엄격한 클라이언트 인증서를 요구
			MinVersion: tls.VersionTLS12,
		},
		Addr: ":18443",
	}*/

	http.HandleFunc("/", handler)
	http.HandleFunc("/digest", handlerDigest)
	http.HandleFunc("/upgrade", handlerUpgrade)
	http.HandleFunc("/chunked", handlerChunkedResponse)
	log.Println("start http listening :18888")
	/*log.Println("start https listening :18443")
	err := server.ListenAndServeTLS("./인증서/server.crt", "./인증서/server.key")*/
	//err := http.ListenAndServeTLS(":18443", "../인증서/server.crt", "../인증서/server.key", nil)
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
	//log.Println(err)
}
