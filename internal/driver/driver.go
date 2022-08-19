package driver

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

//DB holds the database connection
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConnections = 10
const maxIdleDbConnections = 5
const maxLifetimeDbConnections = 5 * time.Minute

//CoonectSQL creates database connection to postgres
func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)
	if err != nil {
		log.Println("cannot connect to DB")
		panic(err)
	}

	db.SetMaxOpenConns(maxOpenDbConnections)
	db.SetConnMaxLifetime(maxLifetimeDbConnections)
	db.SetMaxIdleConns(maxIdleDbConnections)
	dbConn.SQL = db

	err = testDB(db)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
