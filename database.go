package main

import (
	"database/sql"
	"fmt"
	"slack-wails/lib/util"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() *Database {
	dp := util.HomeDir() + "/slack/config.db"
	db, err := sql.Open("sqlite3", dp)
	if err != nil {
		logger.NewDefaultLogger().Error("Failed to open database:" + err.Error())
		return &Database{
			DB: nil,
		}
	}
	err = db.Ping()
	if err != nil {
		logger.NewDefaultLogger().Error("Failed to ping database:" + err.Error())
		return &Database{
			DB: nil,
		}
	}
	return &Database{
		DB: db,
	}
}

func (d *Database) Check() bool {
	if d.DB == nil {
		return false
	}
	return true
}

func (d *Database) CreateTable() bool {
	_, err := d.DB.Exec(`CREATE TABLE IF NOT EXISTS agent_pool ( hots TEXT );`)
	if err != nil {
		return false
	}
	return true
}

func (d *Database) InsertAgentPool(host string) bool {
	stmt, err := d.DB.Prepare("INSERT INTO agent_pool(hosts) VALUES(?)")
	if err != nil {
		return false
	}
	defer stmt.Close()
	tx, err := d.DB.Begin()
	if err != nil {
		return false
	}
	_, err = stmt.Exec(host)
	if err != nil {
		tx.Rollback()
		logger.NewDefaultLogger().Debug(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return false
	}
	return true
}

func (d *Database) SearchAgentPool() (hosts []string) {
	var host string
	rows, err := d.DB.Query("SELECT hosts FROM agent_pool")
	if err != nil {
		return hosts
	}
	for rows.Next() {
		rows.Scan(&host)
		hosts = append(hosts, host)
	}
	return hosts
}

func (d *Database) DeleteAgentPoolField(host string) bool {
	deleteStmt := `DELETE FROM agent_pool WHERE hosts = ?;`
	res, err := d.DB.Exec(deleteStmt, host)
	if err != nil {
		return false
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false
	}
	if rowsAffected >= 1 {
		return true
	}
	return false
}

func (d *Database) DeleteAllField(tableName string) bool {
	deleteStmt := fmt.Sprintf("DELETE FROM %v;", tableName)
	_, err := d.DB.Exec(deleteStmt, tableName)
	if err != nil {
		return false
	}
	return true
}
