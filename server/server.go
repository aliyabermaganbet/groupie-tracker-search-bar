package server

import (
	"fmt"
	"net/http"
)

func Server() {
	http.HandleFunc("/", Homepage)
	http.HandleFunc("/hello", Search)
	http.HandleFunc("/result/", DisplayOutput)
	fmt.Println("click to the http://localhost:8081")
	http.ListenAndServe("localhost:8081", nil)
}
