package sqlite

import (
	"database/sql"
	"github.com/Litekube/network-controller/utils"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
)

const (
	dbDriverName = "sqlite3"
	dbDir        = ".litekube/network-controller/server/"
	dbName       = "litekube-nc.db"
)

var db *sql.DB

func GetDb() *sql.DB {
	return db
}

func InitSqlite(dbPath string) (err error) {
	if dbPath == "" {
		dbPath = filepath.Join(utils.GetHomeDir(), dbDir)
	}
	dbPath = filepath.Join(dbPath, "db")
	utils.CreateDir(dbPath)
	db, err = sql.Open(dbDriverName, filepath.Join(dbPath, dbName))
	if err != nil {
		return
	}
	err = createTable()
	if err != nil {
		return
	}
	return
}

func createTable() error {
	// create table network_mgr
	sql := `create table if not exists "network_mgr" (
		"id" integer primary key autoincrement,
		"token" text not null unique,
		"state" integer not null,
		"bind_ip" text not null default "",
		"create_time" timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime')),
    	"update_time"    timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime'))
	)`

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	// trigger for update_time
	sql = `
	CREATE TRIGGER if not exists update_time_trigger UPDATE OF id,token,state,bind_ip,create_time ON network_mgr
	BEGIN
	  UPDATE network_mgr SET update_time=datetime(CURRENT_TIMESTAMP, 'localtime') WHERE id=OLD.id;
	END
	`
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// create bootstrap token mgr
	sql = `create table if not exists "token_mgr" (
		"id" integer primary key autoincrement,
		"token" text not null unique,
		"expire_time" timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime')),
		"create_time" timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime')),
    	"update_time"    timestamp default (datetime(CURRENT_TIMESTAMP, 'localtime'))
	)`
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	// trigger for update_time
	sql = `
	CREATE TRIGGER if not exists update_time_trigger2 UPDATE OF id,token,create_time,expire_time ON token_mgr
	BEGIN
	  UPDATE token_mgr SET update_time=datetime(CURRENT_TIMESTAMP, 'localtime') WHERE id=OLD.id;
	END
	`

	return err
}
