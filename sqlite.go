package sqlite

import (
	"database/sql"

	sqlite "github.com/mattn/go-sqlite3"
)

const (
	backupPagesPerStep = 1000
)

type SQLite struct {
	DrvConn *sqlite.SQLiteConn
	Conn    *sql.DB
}

func NewSQLite(dsn string) (*SQLite, error) {
	conn, drvConn, err := openSQLiteConn(dsn)
	if err != nil {
		return nil, err
	}
	return &SQLite{DrvConn: drvConn, Conn: conn}, nil
}

func copySQLite(src, dst *sqlite.SQLiteConn) error {
	backup, err := dst.Backup("main", src, "main")
	if err != nil {
		return err
	}
	for {
		allDone, err := backup.Step(backupPagesPerStep)
		if err != nil {
			return err
		}
		if allDone {
			break
		}
	}
	return nil
}

func (s *SQLite) CopyTo(dsn string) error {
	_, drvConn, err := openSQLiteConn(dsn)
	if err != nil {
		return err
	}
	defer drvConn.Close()
	err = copySQLite(drvConn, s.drvConn)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLite) CopyFrom(dsn string) error {
	_, drvConn, err := openSQLiteConn(dsn)
	if err != nil {
		return err
	}
	defer drvConn.Close()
	err = copySQLite(s.drvConn, drvConn)
	if err != nil {
		return err
	}
	return nil
}
