package sqlite

import (
	"database/sql"
	"sync"

	sqlite "github.com/mattn/go-sqlite3"
)

const (
	sqliteDriverName = "sqlite3-ext"
)

var (
	sqliteLock *sync.Mutex
	sqliteConn *sqlite.SQLiteConn
)

func init() {
	sql.Register(sqliteDriverName, &sqlite.SQLiteDriver{
		ConnectHook: func(conn *sqlite.SQLiteConn) error {
			sqliteConn = conn
			return nil
		},
	})
}

func openSQLiteConn(dsn string) (*sql.DB, *sqlite.SQLiteConn, error) {
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
