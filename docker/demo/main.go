package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("start web")
	http.HandleFunc("/go", myHandler)
	http.ListenAndServe(":8080", nil)
}

// handler函数
func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr, "连接成功")
	// 请求方式：GET POST DELETE PUT UPDATE
	fmt.Println("method:", r.Method)
	// /go
	fmt.Println("url:", r.URL.Path)
	fmt.Println("header:", r.Header)
	fmt.Println("body:", r.Body)
	// 回复
	w.Write([]byte("邓瀚宇 哈哈哈哈哈"))
}
