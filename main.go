package main

import (
	"fmt"
	"net/http"
    "example.com/goyasu/api"
    "example.com/goyasu/db/migrations"
)

func main() {
	migrations.SetupDatabase()

	http.HandleFunc("/", api.UploadForm)
	http.HandleFunc("/view", api.ViewPosts)
	http.HandleFunc("/submit-post", api.SubmitPost)
    http.HandleFunc("/edit-post", api.EditPost)
    http.HandleFunc("/update-post", api.UpdatePost)
	fmt.Println("サーバーを起動しています...")
	http.ListenAndServe(":8080", nil)
}

