package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	User                 string
	Passwd               string
	Net                  string
	Addr                 string
	DBName               string
	AllowNativePasswords bool
}

type MySQLDB struct {
	config Config
	db     *sql.DB
}

func New(config Config) *MySQLDB {

	cfg := mysql.Config{
		User:                 config.User,
		Passwd:               config.Passwd,
		Net:                  config.Net,
		Addr:                 config.Addr,
		DBName:               config.DBName,
		AllowNativePasswords: config.AllowNativePasswords,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(fmt.Errorf("can't open mysql db -> %w", err))
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQLDB{config: config, db: db}
}
