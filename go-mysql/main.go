// Created by quicksandzn@gmail.com on 2018/7/25
package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"log"
	"io/ioutil"
	"os/exec"
	"os"
	"strings"
)

type Config struct {
	Mysql struct {
		Path       string `yaml:"path"`
		DriverName string `yaml:"driverName"`
	}
}

const ConfigPath = "./config/conf.yml"

func main() {

	config := Config{}
	buffer, err := ioutil.ReadFile(ConfigPath)
	failOnError(err, "read config error")
	err = yaml.Unmarshal(buffer, &config)
	failOnError(err, "yml convert error")

	db, err := sql.Open(config.Mysql.DriverName, config.Mysql.Path)
	failOnError(err, "connection mysql error")
	defer db.Close()

	//insert(db)
	//update(db)
	query(db)
	//delete(db)
}

// update
func delete(db *sql.DB) {
	stmt, err := db.Prepare("delete from book where id = ?")
	failOnError(err, "delete data error")
	defer stmt.Close()
	res, err := stmt.Exec(2)
	failOnError(err, "delete data error")
	num, err := res.RowsAffected()
	failOnError(err, "delete data error")
	if num > 0 {
		log.Println("delete success")
	} else {
		log.Println("delete error")
	}
}

// update
func update(db *sql.DB) {
	stmt, err := db.Prepare("UPDATE  book set author = ? where id = ?")
	failOnError(err, "update data error")
	defer stmt.Close()
	res, err := stmt.Exec("ddd", 1)
	failOnError(err, "update data error")
	num, err := res.RowsAffected()
	failOnError(err, "update data error")
	if num > 0 {
		log.Println("update success")
	} else {
		log.Println("update error")
	}
}

// insert
func insert(db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO book(id,title,author) VALUES( ?, ?,? )")
	failOnError(err, "update data error")
	defer stmt.Close()
	res, err := stmt.Exec(1, "book1", "book1")
	failOnError(err, "update data error")
	num, err := res.RowsAffected()
	failOnError(err, "update data error")
	if num > 0 {
		log.Println("insert success")
	} else {
		log.Println("insert error")
	}
}

// query
func query(db *sql.DB) {
	r, err := db.Query("select id,title,author from book")
	failOnError(err, "select data error")
	for r.Next() {
		var id int
		var title string
		var author string
		r.Columns()
		err = r.Scan(&id, &title, &author)
		failOnError(err, "select data error")
		log.Println(id, title, author)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	failOnError(err, "getCurrentPath error")
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}
