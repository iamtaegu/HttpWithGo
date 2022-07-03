package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

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

func main() {
	//var httpServer http.Server
	server := &http.Server{
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequestClientCert,
			//ClientAuth: tls.RequireAndVerifyClientCert, //가장 엄격한 클라이언트 인증서를 요구
			MinVersion: tls.VersionTLS12,
		},
		Addr: ":18443",
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/digest", handlerDigest)
	//log.Println("start http listening :18888")
	log.Println("start https listening :18443")
	err := server.ListenAndServeTLS("../인증서/server.crt", "../인증서/server.key")
	//err := http.ListenAndServeTLS(":18443", "../인증서/server.crt", "../인증서/server.key", nil)
	//httpServer.Addr = ":18888"
	//log.Println(httpServer.ListenAndServe())
	log.Println(err)
}
