package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
)

func cookiePost() {
	jar, err := cookiejar.New(nil) // 쿠키 저장할 cookiejar 인스턴스
	if err != nil {
		panic(err)
	}
	client := http.Client{ // 쿠키 저장할 수 있는 http.CLient 인스턴스
		Jar: jar,
	}
	for i := 0; i < 2; i++ {
		resp, err := client.Get("http://localhost:18888/cookie")
		if err != nil {
			panic(err)
		}
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			panic(err)
		}
		log.Println(string(dump))
	}
}

func proxyPost() {
	proxyUrl, err := url.Parse("http://localhost:18888")
	if err != nil {
		panic(err)
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	resp, err := client.Get("http://github.com")
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}

func fileAccessPost() {

}

func main() {
	//cookiePost()
	//proxyPost()
	fileAccessPost()
}
