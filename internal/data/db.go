package data

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"log"
	"os"
	"path/filepath"
	"time"
)

//go:generate go run github.com/fpawel/elco/cmd/utils/sqlstr/...

var (
	DBx *sqlx.DB
	DB  *reform.DB
)

func Years() (years []int) {
	err := DBx.Select(&years, `SELECT DISTINCT year FROM act`)
	if err != nil {
		log.Println("ERROR:", err)
	}
	return
}

func NextActNumber() (n int) {

	_ = DBx.Get(&n, `
SELECT act_number + 1
FROM act
WHERE year = ?
ORDER BY act_number DESC LIMIT 1;`, time.Now().Year())
	if n == 0 {
		n = 1
	}
	return

}

func init() {
	fileName := filepath.Join(filepath.Dir(os.Args[0]), "vimec.sqlite")

	log.Println("data:", fileName)

	conn, err := sql.Open("sqlite3", fileName)
	if err != nil {
		panic(err)
	}
	conn.SetMaxIdleConns(1)
	conn.SetMaxOpenConns(1)
	conn.SetConnMaxLifetime(0)

	if _, err = conn.Exec(SQLCreate); err != nil {
		panic(err)
	}

	DBx = sqlx.NewDb(conn, "sqlite3")
	DB = reform.NewDB(conn, sqlite3.Dialect, nil)
}
