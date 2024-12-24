package migrator

import (
	"GoGameApp/repository/mysql"
	"database/sql"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect    string
	dbConfig   mysql.Config
	migrations *migrate.FileMigrationSource
}

func New(dbConfig mysql.Config) Migrator {
	migrations := &migrate.FileMigrationSource{
		Dir: "./repository/mysql/migrations",
	}
	return Migrator{dbConfig: dbConfig, dialect: "mysql", migrations: migrations}
}

func (m Migrator) Up() {
	db, err := sql.Open(m.dialect, fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true",
		m.dbConfig.User, m.dbConfig.Passwd, m.dbConfig.Addr, m.dbConfig.DBName))
	if err != nil {
		panic(fmt.Errorf("can't open mysql db -> %w", err))
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		panic(fmt.Errorf("can't apply migrations -> %w", err))
	}

	fmt.Printf("INFO - sql-migrate: Applied %d migartions.\n", n)
}

func (m Migrator) Down() {
	db, err := sql.Open(m.dialect, fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true",
		m.dbConfig.User, m.dbConfig.Passwd, m.dbConfig.Addr, m.dbConfig.DBName))
	if err != nil {
		panic(fmt.Errorf("can't open mysql db -> %w", err))
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		panic(fmt.Errorf("can't rollback migrations -> %w", err))
	}

	fmt.Printf("Rollbacked %d migartions.\n", n)
}

func (m Migrator) Status() {
	//? have status?
}
