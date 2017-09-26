package sqlite

import (
	"database/sql"
	"sync"

	sqlite3 "github.com/mattn/go-sqlite3"
)

const (
	sqliteDriverName = "sqlite3-ext"
)

var (
	sqliteLock = &sync.Mutex{}
	sqliteConn *sqlite3.SQLiteConn
)

func init() {
	sql.Register(sqliteDriverName, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			sqliteConn = conn
			return nil
		},
	})
}

func openSQLiteConn(dsn string) (*sql.DB, *sqlite3.SQLiteConn, error) {
	sqliteLock.Lock()
	defer sqliteLock.Unlock()

	conn, err := sql.Open(sqliteDriverName, dsn)
	if err != nil {
		return nil, nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, nil, err
	}
	drvConn := sqliteConn
	return conn, drvConn, nil
}
