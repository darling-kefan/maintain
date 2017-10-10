package model

import (
	"fmt"
	"time"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Mtdb struct{
	conn *sql.DB
}

func NewMtdb(conf MySQLConfig) (*Mtdb, error) {
	conn, err := sql.Open("mysql", conf.dataStoreName())
	if err != nil {
		return nil, fmt.Errorf("mysql: Could not get a connection: %v", err)
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("mysql: Could not establish a good connection: %v", err)
	}

	db := &Mtdb{
		conn: conn,
	}
	return db, nil
}

func (db *Mtdb) Close() {
	db.conn.Close()
}

const TIME_FORMAT = "2006-01-02 15:04:05"

func parseTimestamp(s sql.NullString) (t time.Time, err error){
	if s.Valid {
		t, err = time.Parse(TIME_FORMAT, s.String)
	} else {
		t = time.Time{}
	}
	return
}

func execAffectingRows(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	r, err := stmt.Exec(args...)
	if err != nil {
		return r, fmt.Errorf("mysql: could not execute statement: %v", err)
	}
	_, err = r.RowsAffected()
	if err != nil {
		return r, fmt.Errorf("mysql: could not get rows affected: %v", err)
	}
	return r, nil
}
