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
func (d *Database) CreateTable() bool {
	_, err := d.DB.Exec(`CREATE TABLE IF NOT EXISTS hunter_syntax ( name TEXT, content TEXT );
	CREATE TABLE IF NOT EXISTS quake_syntax ( name TEXT, content TEXT );
	CREATE TABLE IF NOT EXISTS fofa_syntax ( name TEXT, content TEXT );
	CREATE TABLE IF NOT EXISTS agent_pool ( hosts TEXT );`)
	return err == nil
}

func (d *Database) ExecSqlStatement(query string, args ...[]string) bool {
	stmt, err := d.DB.Prepare(query)
	if err != nil {
		return false
	}
	defer stmt.Close()
	tx, err := d.DB.Begin()
	if err != nil {
		return false
	}
	_, err = stmt.Exec(args)
	if err != nil {
		tx.Rollback()
		logger.NewDefaultLogger().Debug(err.Error())
	}
	err = tx.Commit()
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
	insertStmt := fmt.Sprintf("INSERT INTO %v(name, content) VALUES(?,?)", chooseSyntaxDbName(module))
	return d.ExecSqlStatement(insertStmt, []string{name, content})
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
	deleteStmt := fmt.Sprintf("DELETE FROM %v WHERE name = ? AND content = ?", chooseSyntaxDbName(module))
	return d.ExecSqlStatement(deleteStmt, []string{name, content})
}
