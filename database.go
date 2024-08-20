package main

import (
	"database/sql"
	"fmt"
	"os"
	"slack-wails/core/webscan"
	"slack-wails/lib/util"
	"sync"

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
	CREATE TABLE IF NOT EXISTS fofa_syntax ( name TEXT, content TEXT );
	CREATE TABLE IF NOT EXISTS agent_pool ( hosts TEXT );
	CREATE TABLE IF NOT EXISTS poc_workflow ( fingerprint TEXT, poc_path TEXT );`)
	return err == nil
}

var mutex sync.Mutex

func (d *Database) InitPocWorkflow() bool {
	var workflowFile = util.HomeDir() + "/slack/config/workflow.yaml"
	if _, err := os.Stat(workflowFile); err != nil {
		return false
	}
	var count int
	err := d.DB.QueryRow(`SELECT COUNT(*) FROM poc_workflow`).Scan(&count)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		count = 0
	}
	if count == 0 {
		// 初始化工作流，如果成功，写入数据库
		err = webscan.InitWorkflow(workflowFile)
		if err != nil {
			return false
		}
		// 将 webscan.WorkFlowDB 数据写入数据库
		tx, err := d.DB.Begin()
		if err != nil {
			return false
		}

		stmt, err := tx.Prepare(`INSERT INTO poc_workflow (fingerprint, poc_path) VALUES (?, ?)`)
		if err != nil {
			tx.Rollback()
			return false
		}
		defer stmt.Close()
		for fingerprint, pocPaths := range webscan.WorkFlowDB {
			for _, pocPath := range pocPaths {
				_, err := stmt.Exec(fingerprint, pocPath)
				if err != nil {
					tx.Rollback()
					return false
				}
			}
		}

		err = tx.Commit()
		return err == nil
	} else {
		// 载入数据库中的数据到 webscan.WorkFlowDB
		rows, err := d.DB.Query(`SELECT fingerprint, poc_path FROM poc_workflow`)
		if err != nil {
			return false
		}
		defer rows.Close()
		mutex.Lock()
		webscan.WorkFlowDB = make(map[string][]string)
		mutex.Unlock()
		for rows.Next() {
			var fingerprint, pocPath string
			if err := rows.Scan(&fingerprint, &pocPath); err != nil {
				continue
			}
			mutex.Lock()
			webscan.WorkFlowDB[fingerprint] = append(webscan.WorkFlowDB[fingerprint], pocPath)
			mutex.Unlock()
		}
		return rows.Err() == nil
	}
}

func (d *Database) GetFingerPocMap() map[string][]string {
	return webscan.WorkFlowDB
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
