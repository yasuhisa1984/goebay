package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"database/sql"
	"os"
	 _ "github.com/mattn/go-sqlite3"
)

func postForm(w http.ResponseWriter, r *http.Request) {
	formHTML := `
	<!DOCTYPE html>
	<html>
	<head>
	    <title>ブログ投稿</title>
	</head>
	<body>
	    <h1>新しいブログ記事を投稿</h1>
	    <form action="/submit-post" method="post">
	        <div>Title: <input type="text" name="title"></div>
	        <div>Content: <textarea name="content"></textarea></div>
	        <div><input type="submit" value="投稿"></div>
	    </form>
	</body>
	</html>`
	w.Write([]byte(formHTML))
}


func setupDatabase() {
	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS posts (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"content" TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
}


func submitPost(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")

	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO posts(title, content) VALUES (?, ?)")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(title, content)
	if err != nil {
		http.Error(w, "Failed to save the post", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ブログ記事を投稿しました!"))
}


func uploadForm(w http.ResponseWriter, r *http.Request) {
	// HTMLフォームを読み込む
	tmpl, err := template.ParseFiles("upload_form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// フォームをレンダリング
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10MBのファイルまでを処理

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("ファイルの取得に失敗:", err)
		return
	}
	defer file.Close()

	fmt.Printf("アップロードファイル名: %s\n", handler.Filename)

	// ファイルを保存する
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("ファイルの保存に失敗:", err)
		return
	}
	defer f.Close()

	// ファイルのコピー
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println("ファイルのコピーに失敗:", err)
		return
	}

	fmt.Fprintf(w, "ファイル %s がアップロードされました。", handler.Filename)
}

func viewPosts(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT title, content FROM posts")
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	postsHTML := "<h1>ブログ投稿一覧</h1>"
	for rows.Next() {
		var title string
		var content string
		err = rows.Scan(&title, &content)
		if err != nil {
			http.Error(w, "Failed to read post", http.StatusInternalServerError)
			return
		}
		postsHTML += "<h2>" + title + "</h2><p>" + content + "</p><hr>"
	}

	w.Write([]byte(postsHTML))
}

