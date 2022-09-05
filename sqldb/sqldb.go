package sqldb

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	ID      int
	N       string
	C       string
	Created string
}
type Comment struct {
	C_ID      int
	PostID    int
	C_N       string
	C_C       string
	C_Created string
}

func Open() *sql.DB {

	db, err := sql.Open("mysql", fmt.Sprint("root:", os.Getenv("dbpw"), "@tcp(localhost:3306)/crud_go"))

	if err != nil {
		fmt.Println("open error")
		panic(err.Error())
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Not connected")
		panic(err.Error())
	} else {
		fmt.Println("Connected")
	}

	return db
}

func Insert(n string, c string, id int, v bool) error {

	n = strings.TrimSpace(n)
	c = strings.TrimSpace(c)

	if n == "" || c == "" {
		return fmt.Errorf("empty fields")
	}

	xss_inj := regexp.MustCompile(`(<script>|<\/script>|'#)`)

	if xss_inj.MatchString(n) || xss_inj.MatchString(c) {
		return fmt.Errorf("no hacking")
	}

	db, err := sql.Open("mysql", fmt.Sprint("root:", os.Getenv("dbpw"), "@tcp(localhost:3306)/crud_go"))

	if err != nil {
		return fmt.Errorf("there seem to be an issue right now")
	}

	defer db.Close()

	var res sql.Result

	if v {
		res, err = db.Exec("INSERT INTO posts (n, c, Created) VALUES (?, ?, ?)", n, c, time.Now().String()[0:19])
	} else {
		res, err = db.Exec("INSERT INTO comments (postid, c_n, c_c, c_Created) VALUES (?, ?, ?, ?)", id, n, c, time.Now().String()[0:19])
	}

	rows_affected, _ := res.RowsAffected()

	if err != nil || rows_affected != 1 {
		fmt.Println(err)
		return fmt.Errorf("there seem to be an issue right now")
	}

	return nil
}

func SelectAll() ([]Post, error) {

	db, err := sql.Open("mysql", fmt.Sprint("root:", os.Getenv("dbpw"), "@tcp(localhost:3306)/crud_go"))

	var posts []Post

	if err != nil {
		return posts, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM posts")

	if err != nil {
		return posts, err
	}

	defer rows.Close()

	for rows.Next() {

		var p Post

		if err := rows.Scan(&p.ID, &p.N, &p.C, &p.Created); err != nil {
			return posts, err
		}

		posts = append(posts, p)
	}
	fmt.Println(posts)

	return posts, nil
}

func SelectOne(id int) (Post, error) {

	db, err := sql.Open("mysql", fmt.Sprint("root:", os.Getenv("dbpw"), "@tcp(localhost:3306)/crud_go"))
	var p Post

	if err != nil {
		return p, fmt.Errorf("there seem to be an issue right now")
	}

	defer db.Close()

	if err := db.QueryRow("SELECT * FROM posts WHERE id=?", id).Scan(&p.ID, &p.N, &p.C, &p.Created); err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("the post you requested does not exist")
		}
		return p, fmt.Errorf("there seem to be an issue right now")
	}

	fmt.Println(p)

	return p, nil

}

func GetComments() ([]Comment, error) {
	db, err := sql.Open("mysql", fmt.Sprint("root:", os.Getenv("dbpw"), "@tcp(localhost:3306)/crud_go"))
	var comments []Comment

	if err != nil {
		return comments, fmt.Errorf("there seem to be an issue right now")
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM comments")

	if err != nil {
		return comments, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Comment

		if err := rows.Scan(&c.C_ID, &c.PostID, &c.C_N, &c.C_C, &c.C_Created); err != nil {
			return comments, err
		}

		comments = append(comments, c)
	}
	fmt.Println(comments)

	return comments, nil
}

func Delete(id int) bool {
	db, err := sql.Open("mysql", fmt.Sprint("root:", os.Getenv("dbpw"), "@tcp(localhost:3306)/crud_go"))

	if err != nil {
		return false
	}

	defer db.Close()

	res, err := db.Exec(`DELETE FROM posts WHERE id=?`, id)

	if err != nil {
		return false
	}

	if rows_affected, _ := res.RowsAffected(); rows_affected != 1 {
		return false
	}

	fmt.Println(res)

	return true

}

func Update(n string, c string, id int) error {

	n = strings.TrimSpace(n)
	c = strings.TrimSpace(c)

	if n == "" && c == "" {
		return fmt.Errorf("all fields are empty")
	}

	db, err := sql.Open("mysql", fmt.Sprint("root:", os.Getenv("dbpw"), "@tcp(localhost:3306)/crud_go"))

	if err != nil {
		return fmt.Errorf("there seem to be an issue right now")
	}

	defer db.Close()

	res, err := db.Exec(`UPDATE posts SET n=IF(LENGTH(?)=0, n, ?), c=IF(LENGTH(?)=0, c, ?) WHERE id=?`, n, n, c, c, id)

	if err != nil {
		return fmt.Errorf("there seem to be an issue right now")
	}

	if rows_affected, _ := res.RowsAffected(); rows_affected != 1 {
		return fmt.Errorf("there seem to be an issue right now")
	}

	return nil
}
