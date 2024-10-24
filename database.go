package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slack-wails/core/portscan"
	"slack-wails/core/webscan"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"slack-wails/lib/util"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	ctx           context.Context
	DB            *sql.DB // 系统数据库
	lock          sync.RWMutex
	OtherDatabase *sql.DB       // 数据库信息采集时的连接池
	MongoClient   *mongo.Client // mongodb连接池
}

func (d *Database) startup(ctx context.Context) {
	d.ctx = ctx
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
	CREATE TABLE IF NOT EXISTS dirsearch ( path TEXT, times INTEGER );
	CREATE TABLE IF NOT EXISTS dbManager ( nanoid TEXT, scheme TEXT, host TEXT, port INTEGER, username TEXT, password TEXT, notes TEXT );
	CREATE TABLE IF NOT EXISTS scanTask ( task_id TEXT PRIMARY KEY, task_name TEXT, targets TEXT, failed INTEGER, vulnerability INTEGER );
	CREATE TABLE IF NOT EXISTS FingerprintInfo ( task_id TEXT, url TEXT, status INTEGER, length INTEGER, title TEXT, detect TEXT, is_waf INTEGER, waf TEXT, fingerprints TEXT, screenshot TEXT );
	CREATE TABLE IF NOT EXISTS VulnerabilityInfo ( task_id TEXT, template_id TEXT, vuln_name TEXT, protocol TEXT, severity TEXT, vuln_url TEXT, extract TEXT, request TEXT, response TEXT, description TEXT, reference TEXT );
	`)
	return err == nil
}

func (d *Database) ExecSqlStatement(query string, args ...interface{}) bool {
	d.lock.Lock()         // 加锁，防止其他读写操作
	defer d.lock.Unlock() // 函数退出时解锁
	stmt, err := d.DB.Prepare(query)
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[sqlite] exec sql statement step 1: %s", err))
		return false
	}
	defer stmt.Close()
	tx, err := d.DB.Begin()
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[sqlite] exec sql statement step 2: %s", err))
		return false
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		tx.Rollback()
		gologger.Debug(d.ctx, fmt.Sprintf("[sqlite] exec sql statement step 3: %s", err))
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

func (d *Database) GetAllDatabaseConnections() (dcs []structs.DatabaseConnection) {
	rows, err := d.DB.Query(`SELECT * FROM dbManager;`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var nanoid, scheme, host, username, password, notes string
		var port int
		err = rows.Scan(&nanoid, &scheme, &host, &port, &username, &password, &notes)
		if err != nil {
			return
		}
		dcs = append(dcs, structs.DatabaseConnection{
			Nanoid:   nanoid,
			Scheme:   scheme,
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
			Notes:    notes,
		})
	}
	return dcs
}

func (d *Database) AddConnection(info structs.DatabaseConnection) bool {
	return d.ExecSqlStatement("INSERT INTO dbManager (nanoid, scheme, host, port, username, password, notes) VALUES (?, ?, ?, ?, ?, ?, ?)", info.Nanoid, info.Scheme, info.Host, info.Port, info.Username, info.Password, info.Notes)
}

func (d *Database) RemoveConnection(nanoid string) bool {
	return d.ExecSqlStatement("DELETE FROM dbManager WHERE nanoid = ?", nanoid)
}

func (d *Database) UpdateConnection(info structs.DatabaseConnection) bool {
	return d.ExecSqlStatement("UPDATE dbManager SET scheme = ?, host = ?, port = ? , username = ?, password = ?, notes = ? WHERE nanoid = ?", info.Scheme, info.Host, info.Port, info.Username, info.Password, info.Notes, info.Nanoid)
}

func (d *Database) ConnectDatabase(info structs.DatabaseConnection) bool {
	var (
		flag           bool
		err            error
		dataSourceName string
	)
	host := fmt.Sprintf("%s:%d", info.Host, info.Port)

	// Determine connection based on the scheme
	switch info.Scheme {
	case "mysql":
		flag, err = portscan.MysqlConn(host, info.Username, info.Password)
		dataSourceName = fmt.Sprintf("%v:%v@tcp(%v)/mysql?charset=utf8&timeout=%v", info.Username, info.Password, host, 10*time.Second)
	case "mssql":
		flag, err = portscan.MssqlConn(host, info.Username, info.Password)
		dataSourceName = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%v;encrypt=disable;timeout=%v", info.Host, info.Username, info.Password, info.Port, 10*time.Second)
	case "oracle":
		flag, err = portscan.OracleConn(host, info.Username, info.Password)
		dataSourceName = fmt.Sprintf("oracle://%s:%s@%s/orcl", info.Username, info.Password, host)
	case "postgres":
		flag, err = portscan.PostgresConn(host, info.Username, info.Password)
		dataSourceName = fmt.Sprintf("postgres://%v:%v@%v/postgres?sslmode=disable", info.Username, info.Password, host)
	case "mongodb":
		flag, err = portscan.MongodbConn(host, info.Username, info.Password)
		if err == nil && flag {
			d.MongoClient, err = d.ConnectMongodb(host, info.Username, info.Password)
			if err != nil {
				d.showErrorMessage(err.Error())
				return false
			}
			return true
		}
	default:
		return false
	}

	// Handle connection failure
	if err != nil || !flag {
		d.showErrorMessage(err.Error())
		return false
	}

	// Connect to other databases
	d.OtherDatabase, err = sql.Open(info.Scheme, dataSourceName)
	if err != nil {
		d.showErrorMessage("认证正确，但无法连接数据库")
		return false
	}
	return true
}

// Helper function to show error messages
func (d *Database) showErrorMessage(message string) {
	runtime.MessageDialog(d.ctx, runtime.MessageDialogOptions{
		Title:         "提示",
		Message:       message,
		Type:          runtime.ErrorDialog,
		DefaultButton: "Ok",
	})
}

func (d *Database) ConnectMongodb(host, user, pass string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create MongoDB URI
	mongoURI := fmt.Sprintf("mongodb://%s", host)

	// Define client options with or without credentials
	clientOpts := options.Client().ApplyURI(mongoURI)
	if user != "" && pass != "" {
		credentials := options.Credential{
			Username: user,
			Password: pass,
		}
		clientOpts.SetAuth(credentials)
	}

	// Connect to MongoDB
	return mongo.Connect(ctx, clientOpts)
}

func (d *Database) FetchDatabaseinfoFromMongodb() map[string][]string {
	var databasesInfo = make(map[string][]string)
	// Get the total number of databases
	if d.MongoClient == nil {
		return databasesInfo
	}
	databaseNames, err := d.MongoClient.ListDatabaseNames(context.TODO(), bson.D{})
	if err != nil {
		gologger.Debug(d.ctx, err)
		return databasesInfo
	}
	// Loop through each database and count the collections
	for _, dbName := range databaseNames {
		db := d.MongoClient.Database(dbName)
		collections, err := db.ListCollectionNames(context.TODO(), bson.D{})
		if err != nil {
			gologger.Debug(d.ctx, err)
			continue
		}
		databasesInfo[dbName] = collections
	}
	return databasesInfo
}

func (d *Database) DisconnectDatabase(scheme string) bool {
	if scheme == "mongodb" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return d.MongoClient.Disconnect(ctx) == nil
	}
	return d.OtherDatabase.Close() == nil
}

var mysql_system_db = []string{"performance_schema", "information_schema", "mysql", "sys"}

func (d *Database) FetchDatabaseinfoFromMysql() map[string][]string {
	// Retrieve all databases
	databases, err := d.OtherDatabase.Query("SHOW DATABASES")
	if err != nil {
		gologger.Warning(d.ctx, fmt.Sprintf("[mysql] 查询数据库失败: %v", err))
		return nil
	}
	defer databases.Close()

	// Iterate over each database
	var databasesInfo = make(map[string][]string)
	for databases.Next() {
		var dbName string
		if err := databases.Scan(&dbName); err != nil {
			gologger.Warning(d.ctx, fmt.Sprintf("[mysql] 扫描数据库名称失败: %v", err))
			continue
		}
		if util.ArrayContains(dbName, mysql_system_db) {
			continue
		}

		// Retrieve tables for the current database
		tables, err := d.OtherDatabase.Query("SHOW TABLES FROM " + dbName)
		if err != nil {
			gologger.Warning(d.ctx, fmt.Sprintf("[mysql] 查询表失败: %v", err))
			continue
		}
		defer tables.Close()

		// Iterate over tables in the current database
		for tables.Next() {
			var tableName string
			if err := tables.Scan(&tableName); err != nil {
				gologger.Warning(d.ctx, fmt.Sprintf("[mysql] 扫描表名称失败: %v", err))
				continue
			}
			databasesInfo[dbName] = append(databasesInfo[dbName], tableName)
		}

		// Check for any error encountered during iteration
		if err := tables.Err(); err != nil {
			gologger.Warning(d.ctx, fmt.Sprintf("[mysql] 遍历表时出错: %v", err))
		}
	}

	// Check for any error encountered during iteration
	if err := databases.Err(); err != nil {
		gologger.Warning(d.ctx, fmt.Sprintf("[mysql] 遍历数据库时出错: %v", err))
	}

	return databasesInfo
}

func (d *Database) FetchTableInfoFromMysql(dbName, tableName string) structs.RowData {
	// Construct SQL query to fetch the first three rows of the specified table
	sqlQuery := fmt.Sprintf("SELECT * FROM `%s`.`%s` LIMIT 3", dbName, tableName)
	_, err := d.OtherDatabase.Exec("set global show_compatibility_56 = ON")
	if err != nil {
		gologger.Debug(d.ctx, err)
	}
	// Execute the query
	rows, err := d.OtherDatabase.Query(sqlQuery)
	if err != nil {
		gologger.Debug(d.ctx, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[mysql] 获取列名失败: %v", err))
	}

	// 创建一个二维切片来存储数据
	var data [][]interface{}

	// 遍历每一行
	for rows.Next() {
		// 创建一个切片来存储每一行的值
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		// 扫描行数据
		if err := rows.Scan(values...); err != nil {
			gologger.Debug(d.ctx, fmt.Sprintf("[mysql] 扫描行数据失败: %v", err))

		}

		// 将每一行的数据附加到切片中
		row := make([]interface{}, len(columns))
		for i, v := range values {
			row[i] = *(v.(*interface{}))
		}
		data = append(data, row)
	}

	return structs.RowData{
		Columns: columns,
		Rows:    data,
	}
}

var sqlserver_system_db = []string{"master", "tempdb", "model", "msdb"}

func (d *Database) FetchDatabaseinfoFromSqlServer() map[string][]string {
	// Retrieve all databases
	databases, err := d.OtherDatabase.Query("SELECT name FROM sys.databases")
	if err != nil {
		gologger.Warning(d.ctx, fmt.Sprintf("[sqlserver] 查询数据库失败: %v", err))
		return nil
	}
	defer databases.Close()

	// Iterate over each database
	var databasesInfo = make(map[string][]string)
	for databases.Next() {
		var dbName string
		if err := databases.Scan(&dbName); err != nil {
			gologger.Warning(d.ctx, fmt.Sprintf("[sqlserver] 扫描数据库名称失败: %v", err))
			continue
		}
		if util.ArrayContains(dbName, sqlserver_system_db) {
			continue
		}

		// Retrieve tables for the current database
		tables, err := d.OtherDatabase.Query(fmt.Sprintf("SELECT TABLE_NAME FROM [%s].INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE'", dbName))
		if err != nil {
			gologger.Warning(d.ctx, fmt.Sprintf("[sqlserver] 查询表失败: %v", err))
			continue
		}
		defer tables.Close()

		// Iterate over tables in the current database
		for tables.Next() {
			var tableName string
			if err := tables.Scan(&tableName); err != nil {
				gologger.Warning(d.ctx, fmt.Sprintf("[sqlserver] 扫描表名称失败: %v", err))
				continue
			}
			databasesInfo[dbName] = append(databasesInfo[dbName], tableName)
		}

		// Check for any error encountered during iteration
		if err := tables.Err(); err != nil {
			gologger.Warning(d.ctx, fmt.Sprintf("[sqlserver] 遍历表时出错: %v", err))
		}
	}

	// Check for any error encountered during iteration
	if err := databases.Err(); err != nil {
		gologger.Warning(d.ctx, fmt.Sprintf("[sqlserver] 遍历数据库时出错: %v", err))
	}

	return databasesInfo
}

func (d *Database) FetchTableInfoFromSqlServer(dbName, tableName string) structs.RowData {
	// Construct SQL query to fetch the first three rows of the specified table
	sqlQuery := fmt.Sprintf("SELECT TOP 3 * FROM [%s].dbo.[%s]", dbName, tableName)

	// Execute the query
	rows, err := d.OtherDatabase.Query(sqlQuery)
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[sqlserver] 查询数据失败: %v", err))
		return structs.RowData{}
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[sqlserver] 获取列名失败: %v", err))
	}

	// Create a two-dimensional slice to store data
	var data [][]interface{}

	// Iterate over each row
	for rows.Next() {
		// Create a slice to store each row's values
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		// Scan the row data
		if err := rows.Scan(values...); err != nil {
			gologger.Debug(d.ctx, fmt.Sprintf("[sqlserver] 扫描行数据失败: %v", err))
		}

		// Append each row's data to the slice
		row := make([]interface{}, len(columns))
		for i, v := range values {
			row[i] = *(v.(*interface{}))
		}
		data = append(data, row)
	}

	return structs.RowData{
		Columns: columns,
		Rows:    data,
	}
}

var oracleSystemSchemas = []string{"SYS", "SYSTEM", "OUTLN", "XDB", "DBSNMP", "APPQOSSYS", "CTXSYS", "ORDDATA"}

func (d *Database) FetchDatabaseInfoFromOracle() map[string][]string {
	// Query to retrieve all schemas (users) from the Oracle database
	schemas, err := d.OtherDatabase.Query("SELECT USERNAME FROM ALL_USERS")
	if err != nil {
		gologger.Warning(d.ctx, fmt.Sprintf("[oracle] 查询模式失败: %v", err))
		return nil
	}
	defer schemas.Close()

	var databaseInfo = make(map[string][]string)

	for schemas.Next() {
		var schemaName string
		if err := schemas.Scan(&schemaName); err != nil {
			gologger.Warning(d.ctx, fmt.Sprintf("[oracle] 扫描模式名称失败: %v", err))
			continue
		}

		// Skip system schemas
		if util.ArrayContains(schemaName, oracleSystemSchemas) {
			continue
		}

		// Query to retrieve tables for the current schema
		tables, err := d.OtherDatabase.Query(fmt.Sprintf("SELECT TABLE_NAME FROM ALL_TABLES WHERE OWNER = '%s'", schemaName))
		if err != nil {
			gologger.Warning(d.ctx, fmt.Sprintf("[oracle] 查询表失败: %v", err))
			continue
		}
		defer tables.Close()

		for tables.Next() {
			var tableName string
			if err := tables.Scan(&tableName); err != nil {
				gologger.Warning(d.ctx, fmt.Sprintf("[oracle] 扫描表名称失败: %v", err))
				continue
			}
			databaseInfo[schemaName] = append(databaseInfo[schemaName], tableName)
		}

		if err := tables.Err(); err != nil {
			gologger.Warning(d.ctx, fmt.Sprintf("[oracle] 遍历表时出错: %v", err))
		}
	}

	if err := schemas.Err(); err != nil {
		gologger.Warning(d.ctx, fmt.Sprintf("[oracle] 遍历模式时出错: %v", err))
	}

	return databaseInfo
}

func (d *Database) FetchTableInfoFromOracle(schemaName, tableName string) structs.RowData {
	// Construct SQL query to fetch the first three rows of the specified table
	sqlQuery := fmt.Sprintf("SELECT * FROM %s.%s FETCH FIRST 3 ROWS ONLY", schemaName, tableName)

	// Execute the query
	rows, err := d.OtherDatabase.Query(sqlQuery)
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[oracle] 查询表数据失败: %v", err))
		return structs.RowData{}
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[oracle] 获取列名失败: %v", err))
		return structs.RowData{}
	}

	// Create a 2D slice to store the data
	var data [][]interface{}

	// Iterate through each row
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		if err := rows.Scan(values...); err != nil {
			gologger.Debug(d.ctx, fmt.Sprintf("[oracle] 扫描行数据失败: %v", err))
			continue
		}

		row := make([]interface{}, len(columns))
		for i, v := range values {
			row[i] = *(v.(*interface{}))
		}
		data = append(data, row)
	}

	return structs.RowData{
		Columns: columns,
		Rows:    data,
	}
}

var postgresSystemSchemas = []string{"pg_catalog", "information_schema", "pg_toast"}

func (d *Database) FetchDatabaseInfoFromPostgres() map[string][]string {
	// Query to retrieve all schemas in PostgreSQL
	schemas, err := d.OtherDatabase.Query("SELECT schema_name FROM information_schema.schemata")
	if err != nil {
		gologger.Warning(d.ctx, fmt.Sprintf("[postgres] 查询模式失败: %v", err))
		return nil
	}
	defer schemas.Close()

	var databaseInfo = make(map[string][]string)

	for schemas.Next() {
		var schemaName string
		if err := schemas.Scan(&schemaName); err != nil {
			gologger.Warning(d.ctx, fmt.Sprintf("[postgres] 扫描模式名称失败: %v", err))
			continue
		}

		// Skip system schemas
		if util.ArrayContains(schemaName, postgresSystemSchemas) {
			continue
		}

		// Query to retrieve tables for the current schema
		tables, err := d.OtherDatabase.Query(fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'", schemaName))
		if err != nil {
			gologger.Warning(d.ctx, fmt.Sprintf("[postgres] 查询表失败: %v", err))
			continue
		}
		defer tables.Close()

		for tables.Next() {
			var tableName string
			if err := tables.Scan(&tableName); err != nil {
				gologger.Warning(d.ctx, fmt.Sprintf("[postgres] 扫描表名称失败: %v", err))
				continue
			}
			databaseInfo[schemaName] = append(databaseInfo[schemaName], tableName)
		}

		if err := tables.Err(); err != nil {
			gologger.Warning(d.ctx, fmt.Sprintf("[postgres] 遍历表时出错: %v", err))
		}
	}

	if err := schemas.Err(); err != nil {
		gologger.Warning(d.ctx, fmt.Sprintf("[postgres] 遍历模式时出错: %v", err))
	}

	return databaseInfo
}

func (d *Database) FetchTableInfoFromPostgres(schemaName, tableName string) structs.RowData {
	// Construct SQL query to fetch the first three rows of the specified table
	sqlQuery := fmt.Sprintf("SELECT * FROM %s.%s LIMIT 3", schemaName, tableName)

	// Execute the query
	rows, err := d.OtherDatabase.Query(sqlQuery)
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[postgres] 查询表数据失败: %v", err))
		return structs.RowData{}
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[postgres] 获取列名失败: %v", err))
		return structs.RowData{}
	}

	// Create a 2D slice to store the data
	var data [][]interface{}

	// Iterate through each row
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		if err := rows.Scan(values...); err != nil {
			gologger.Debug(d.ctx, fmt.Sprintf("[postgres] 扫描行数据失败: %v", err))
			continue
		}

		row := make([]interface{}, len(columns))
		for i, v := range values {
			row[i] = *(v.(*interface{}))
		}
		data = append(data, row)
	}

	return structs.RowData{
		Columns: columns,
		Rows:    data,
	}
}

func (d *Database) GetAllScanTask() []structs.TaskResult {
	rows, err := d.DB.Query(`SELECT * FROM scanTask;`)
	if err != nil {
		return []structs.TaskResult{}
	}
	defer rows.Close()
	var tasks []structs.TaskResult
	for rows.Next() {
		var task structs.TaskResult
		err = rows.Scan(&task.TaskId, &task.TaskName, &task.Targets, &task.Failed, &task.Vulnerability)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks
}

func (d *Database) InsertScanTask(taskid, taskname, targets string, failed, vulnerability int) bool {
	insertStmt := "INSERT INTO scanTask (task_id, task_name, targets, failed, vulnerability) VALUES (?, ?, ?, ?, ?)"
	return d.ExecSqlStatement(insertStmt, taskid, taskname, targets, failed, vulnerability)
}

// 改变扫描结果，漏洞数量
func (d *Database) UpdateScanWithResult(taskid string, failed, vulnerability int) bool {
	updateStmt := "UPDATE scanTask SET failed = ?, vulnerability = ? WHERE task_id = ?"
	return d.ExecSqlStatement(updateStmt, failed, vulnerability, taskid)
}

func (d *Database) DeleteScanTask(taskid string) bool {
	deleteStmt := "DELETE FROM scanTask WHERE task_id = ?"
	isSuccess := d.ExecSqlStatement(deleteStmt, taskid)
	if isSuccess {
		d.ExecSqlStatement("DELETE FROM FingerprintInfo WHERE task_id = ?", taskid)
		d.ExecSqlStatement("DELETE FROM VulnerabilityInfo WHERE task_id = ?", taskid)
	}
	return isSuccess
}

func (d *Database) SelectFingerscanResult(taskid string) []webscan.InfoResult {
	rows, err := d.DB.Query("SELECT * FROM FingerprintInfo WHERE task_id = ?;", taskid)
	if err != nil {
		gologger.Debug(d.ctx, err)
		return []webscan.InfoResult{}
	}
	defer rows.Close()
	var results []webscan.InfoResult
	for rows.Next() {
		var result webscan.InfoResult
		var fingerprintsStr string
		var task_id string
		err = rows.Scan(&task_id, &result.URL, &result.StatusCode, &result.Length, &result.Title, &result.Detect, &result.IsWAF, &result.WAF, &fingerprintsStr, &result.Screenshot)
		if err != nil {
			gologger.Debug(d.ctx, err)
			continue
		}
		result.Fingerprints = strings.Split(fingerprintsStr, ",")
		results = append(results, result)
	}
	return results
}

func (d *Database) SelectPocscanResult(taskid string) []webscan.VulnerabilityInfo {
	rows, err := d.DB.Query("SELECT * FROM VulnerabilityInfo WHERE task_id = ?", taskid)
	if err != nil {
		return []webscan.VulnerabilityInfo{}
	}
	defer rows.Close()
	var results []webscan.VulnerabilityInfo
	for rows.Next() {
		var result webscan.VulnerabilityInfo
		var task_id string
		err = rows.Scan(&task_id, &result.ID, &result.Name, &result.Type, &result.Risk, &result.URL, &result.Extract, &result.Request, &result.Response, &result.Description, &result.Reference)
		if err != nil {
			continue
		}
		results = append(results, result)
	}
	return results
}

func (d *Database) InsertFingerscanResult(taskid string, result webscan.InfoResult) bool {
	insertStmt := "INSERT INTO FingerprintInfo (task_id, url, status, length, title, detect, is_waf, waf, fingerprints, screenshot) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return d.ExecSqlStatement(insertStmt, taskid, result.URL, result.StatusCode, result.Length, result.Title, result.Detect, result.IsWAF, result.WAF, strings.Join(result.Fingerprints, ","), result.Screenshot)
}

func (d *Database) InsertPocscanResult(taskid string, result webscan.VulnerabilityInfo) bool {
	insertStmt := "INSERT INTO VulnerabilityInfo (task_id, template_id, vuln_name, protocol, severity, vuln_url, extract, request, response, description, reference) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return d.ExecSqlStatement(insertStmt, taskid, result.ID, result.Name, result.Type, result.Risk, result.URL, result.Extract, result.Request, result.Response, result.Description, result.Reference)
}
