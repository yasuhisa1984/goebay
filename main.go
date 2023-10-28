package main

import (
	"fmt"
	"net/http"
)

func main() {
	setupDatabase()

	http.HandleFunc("/", uploadForm)
	http.HandleFunc("/view", viewPosts)
	http.HandleFunc("/submit-post", submitPost)
	fmt.Println("サーバーを起動しています...")
	http.ListenAndServe(":8080", nil)
}

