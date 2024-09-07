package main

import (
	"database/sql"
	"fmt"
	"os"
	"slack-wails/lib/util"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

type Database struct {
	DB   *sql.DB
	lock sync.RWMutex
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
	CREATE TABLE IF NOT EXISTS agent_pool ( hosts TEXT );
	CREATE TABLE IF NOT EXISTS dirsearch ( path TEXT, times INTEGER )`)
	return err == nil
}

func (d *Database) ExecSqlStatement(query string, args ...interface{}) bool {
	d.lock.Lock()         // 加锁，防止其他读写操作
	defer d.lock.Unlock() // 函数退出时解锁
	stmt, err := d.DB.Prepare(query)
	if err != nil {
		return false
	}
	defer stmt.Close()
	tx, err := d.DB.Begin()
	if err != nil {
		return false
	}
	_, err = stmt.Exec(args...)
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
	return d.ExecSqlStatement(insertStmt, name, content)
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
	return d.ExecSqlStatement(deleteStmt, name, content)
}

func (d *Database) UpdateOrInsertPath(path string) bool {
	d.lock.Lock()
	defer d.lock.Unlock()
	// 开始事务
	tx, err := d.DB.Begin()
	if err != nil {
		return false
	}
	// 尝试更新记录，如果path存在，则times增加1
	result, err := tx.Exec(`
        UPDATE dirsearch 
        SET times = times + 1 
        WHERE path = ?`, path)

	if err != nil {
		tx.Rollback()
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return false
	}

	// 如果更新受影响的行数为0，说明path不存在，需要插入新记录
	if rowsAffected == 0 {
		_, err = tx.Exec(`
            INSERT INTO dirsearch (path, times) 
            VALUES (?, 1)`, path)

		if err != nil {
			tx.Rollback()
			return false
		}
	}

	return tx.Commit() == nil
}

type pathTimes struct {
	Path  string
	Times int
}

func (d *Database) GetAllPathsAndTimes() []pathTimes {
	d.lock.RLock()         // 读操作前加锁
	defer d.lock.RUnlock() // 函数结束时解锁

	rows, err := d.DB.Query("SELECT path, times FROM dirsearch")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var results []pathTimes

	for rows.Next() {
		var ds pathTimes
		err := rows.Scan(&ds.Path, &ds.Times)
		if err != nil {
			return nil
		}
		results = append(results, ds)
	}

	// 检查是否有遍历错误
	if err = rows.Err(); err != nil {
		return nil
	}

	return results
}
