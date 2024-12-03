package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLDB struct {
	db *sql.DB
}

func New() *MySQLDB {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "gameapp_local",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(fmt.Errorf("can't open mysql db -> %w", err))
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQLDB{db: db}
}
