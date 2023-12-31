package api

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"database/sql"
	"os"
	 _ "github.com/mattn/go-sqlite3"
)

func ViewPosts(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータからメッセージを取得
	message := r.URL.Query().Get("message")

	db, err := sql.Open("sqlite3", "db/blog.db")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// メッセージを表示
	if message != "" {
		fmt.Fprintf(w, "<p>%s</p>", message)
	}

	rows, err := db.Query("SELECT id, title, content FROM posts")
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []struct {
		Id      int
		Title   string
		Content string
	}
	for rows.Next() {
		var id int
		var title string
		var content string
		err = rows.Scan(&id, &title, &content)
		if err != nil {
			http.Error(w, "Failed to read post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, struct {
			Id      int
			Title   string
			Content string
		}{Id: id, Title: title, Content: content})
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/view_posts.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Posts []struct {
			Id      int
			Title   string
			Content string
		}
	}{
		Title: "ブログ投稿一覧",
		Posts: posts,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
    id := r.PostFormValue("id")
    title := r.PostFormValue("title")
    content := r.PostFormValue("content")

    db, err := sql.Open("sqlite3", "db/blog.db")
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    _, err = db.Exec("UPDATE posts SET title=?, content=? WHERE id=?", title, content, id)
    if err != nil {
        http.Error(w, "Failed to update the post", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/view", http.StatusSeeOther)
}


func SubmitPost(w http.ResponseWriter, r *http.Request) {
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

    // リダイレクト先とメッセージをクエリパラメータとして設定
	http.Redirect(w, r, "/view?message=ブログを投稿しました", http.StatusSeeOther)
}


func UploadForm(w http.ResponseWriter, r *http.Request) {
	// ベーステンプレートを読み込み
	tmpl, err := template.ParseFiles("templates/base.html", "templates/upload_form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// フォームをレンダリング（upload_form.html を呼び出す）
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
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



func EditPost(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("sqlite3", "db/blog.db")  // パスを修正
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        fmt.Println("Database error:", err)  // エラーの詳細をログに出力
        return
    }
    defer db.Close()

    row := db.QueryRow("SELECT id, title, content FROM posts WHERE id=?", id)

    var post struct {
        ID      int
        Title   string
        Content string
    }
    err = row.Scan(&post.ID, &post.Title, &post.Content)
    if err != nil {
        http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
        fmt.Println("Row scan error:", err)  // エラーの詳細をログに出力
        return
    }

    tmpl, err := template.ParseFiles("templates/base.html", "templates/edit_post.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Println("Template parse error:", err)  // エラーの詳細をログに出力
        return
    }

    tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Post": post,
    })
}


