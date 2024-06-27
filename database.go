package main

import (
	"database/sql"
	"fmt"
	"os"
	"slack-wails/lib/util"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() *Database {
	os.Mkdir(util.HomeDir()+"/slack", 0777) // 创建配置文件夹
	dp := util.HomeDir() + "/slack/config.db"
	db, err := sql.Open("sqlite3", dp) // 创建数据库文件
	if err != nil {
		return &Database{
			DB: nil,
		}
	}
	err = db.Ping()
	if err != nil {
		return &Database{
			DB: nil,
		}
	}
	return &Database{
		DB: db,
	}
}

func (d *Database) Check() bool {
	return d.DB == nil
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
	return err == nil
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
	return err == nil
}

func (d *Database) CreateTable() bool {
	_, err := d.DB.Exec(`CREATE TABLE IF NOT EXISTS hunter_syntax ( name TEXT, content TEXT );
	CREATE TABLE IF NOT EXISTS quake_syntax ( name TEXT, content TEXT );
	CREATE TABLE IF NOT EXISTS agent_pool ( hosts TEXT );`)
	return err == nil
}

type Syntax struct {
	Name    string
	Content string
}

func (d *Database) SelectAllSyntax(module string) (data []Syntax) {
	rows, err := d.DB.Query(fmt.Sprintf(`SELECT name, content FROM %v;`, chooseSyntaxDbName(module)))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name, content string
		err = rows.Scan(&name, &content)
		if err != nil {
			return
		}
		data = append(data, Syntax{
			Name:    name,
			Content: content,
		})
	}
	return
}

func (d *Database) InsertFavGrammarFiled(module, name, content string) bool {
	stmt, err := d.DB.Prepare(fmt.Sprintf("INSERT INTO %v(name, content) VALUES(?,?)", chooseSyntaxDbName(module)))
	if err != nil {
		return false
	}
	defer stmt.Close()
	tx, err := d.DB.Begin()
	if err != nil {
		return false
	}
	_, err = stmt.Exec(name, content)
	if err != nil {
		tx.Rollback()
		logger.NewDefaultLogger().Debug(err.Error())
	}
	err = tx.Commit()
	return err == nil
}

func chooseSyntaxDbName(name string) string {
	switch name {
	case "quake":
		return "quake_syntax"
	case "hunter":
		return "hunter_syntax"
	default:
		return "fofa_syntax"
	}
}

func (d *Database) RemoveFavGrammarFiled(module, name, content string) bool {
	stmt, err := d.DB.Prepare(fmt.Sprintf("DELETE FROM %v WHERE name = ? AND content = ?", chooseSyntaxDbName(module)))
	if err != nil {
		return false
	}
	defer stmt.Close()
	tx, err := d.DB.Begin()
	if err != nil {
		return false
	}
	_, err = stmt.Exec(name, content)
	if err != nil {
		tx.Rollback()
		logger.NewDefaultLogger().Debug(err.Error())
	}
	err = tx.Commit()
	return err == nil
}
