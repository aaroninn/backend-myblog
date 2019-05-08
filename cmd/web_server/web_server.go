package main

import (
	// "github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./dist")))
	if err := http.ListenAndServe("0.0.0.0:80", nil); err != nil {
		panic(err)
	}
	// engine := gin.Default()
	// engine.Static("/static", "./dist/static")
	// engine.StaticFile("/", "./dist/index.html")
	// engine.Run(":80")
}
