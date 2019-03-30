package main

import "net/http"

func main() {
	http.Handle("/", http.FileServer(http.Dir("./dist")))
	if err := http.ListenAndServe("0.0.0.0:80", nil); err != nil {
		panic(err)
	}
}
