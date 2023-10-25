package main

import (
	"fmt"
	"net/http"
)

func main() {
	setupDatabase() // データベースのセットアップ
	http.HandleFunc("/", postForm)
	http.HandleFunc("/submit-post", submitPost)
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/view-posts", viewPosts)
	fmt.Println("サーバーを起動しています...")
	http.ListenAndServe(":8080", nil)
}

