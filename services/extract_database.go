package services

import (
	"context"
	"database/sql"
	"fmt"
	"slack-wails/core/portscan"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"slack-wails/lib/util"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
		flag, err = portscan.OracleConn(host, info.ServerName, info.Username, info.Password)
		dataSourceName = fmt.Sprintf("oracle://%s:%s@%s/%s", info.Username, info.Password, host, info.ServerName)
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

	// Query to get the total row count for the table
	var rowCount int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s`.`%s`", dbName, tableName)
	err = d.OtherDatabase.QueryRow(countQuery).Scan(&rowCount)
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[mysql] 获取总行数失败: %v", err))
	}

	return structs.RowData{
		Columns:   columns,
		Rows:      data,
		RowsCount: rowCount,
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
	// Query to get the total row count for the table
	var rowCount int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM [%s].dbo.[%s]", dbName, tableName)
	err = d.OtherDatabase.QueryRow(countQuery).Scan(&rowCount)
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[mysql] 获取总行数失败: %v", err))
	}

	return structs.RowData{
		Columns:   columns,
		Rows:      data,
		RowsCount: rowCount,
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
	// Query to get the total row count for the table
	var rowCount int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", schemaName, tableName)
	err = d.OtherDatabase.QueryRow(countQuery).Scan(&rowCount)
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[oracle] 获取总行数失败: %v", err))
	}

	return structs.RowData{
		Columns:   columns,
		Rows:      data,
		RowsCount: rowCount,
	}
}

var postgresSystemSchemas = []string{"pg_catalog", "template0", "template1", "information_schema", "postgres", "pg_toast"}

func (d *Database) FetchDatabaseInfoFromPostgres(info structs.DatabaseConnection) map[string][]string {
	d.PostgresInfo = &info
	// Query to retrieve all schemas in PostgreSQL
	schemas, err := d.OtherDatabase.Query("SELECT datname FROM pg_database WHERE datistemplate = false")
	if err != nil {
		gologger.Warning(d.ctx, fmt.Sprintf("[postgres] 查询数据库失败: %v", err))
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

		tables, err := getPostgresTables(info, schemaName)
		if err != nil {
			gologger.Warning(d.ctx, err)
			continue
		}

		databaseInfo[schemaName] = tables
	}

	if err := schemas.Err(); err != nil {
		gologger.Warning(d.ctx, fmt.Sprintf("[postgres] 遍历模式时出错: %v", err))
	}

	return databaseInfo
}

func getPostgresTables(info structs.DatabaseConnection, dbName string) ([]string, error) {
	// 切换到目标数据库
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", info.Host, info.Port, info.Username, info.Password, dbName)
	targetDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("[postgres] 连接数据库查询表[%s]失败: %v", dbName, err)
	}
	defer targetDB.Close()

	// 查询当前数据库中所有表
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	rows, err := targetDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("[postgres] 数据库[%s]查询表名失败: %v", dbName, err)
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			continue
		}
		tables = append(tables, tableName)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("[postgres] 遍历表名时出错: %v", err)
	}
	return tables, nil
}

func (d *Database) FetchTableInfoFromPostgres(schemaName, tableName string) structs.RowData {
	// 切换到目标数据库
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", d.PostgresInfo.Host, d.PostgresInfo.Port, d.PostgresInfo.Username, d.PostgresInfo.Password, schemaName)
	targetDB, err := sql.Open("postgres", connStr)
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[postgres] 连接数据库失败: %v", err))
	}
	defer targetDB.Close()

	// Construct SQL query to fetch the first three rows of the specified table
	sqlQuery := fmt.Sprintf(`SELECT * FROM %s LIMIT 3`, tableName)

	// Execute the query
	rows, err := targetDB.Query(sqlQuery)
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
	// Query to get the total row count for the table
	var rowCount int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err = targetDB.QueryRow(countQuery).Scan(&rowCount)
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[postgres] 获取总行数失败: %v", err))
	}

	return structs.RowData{
		Columns:   columns,
		Rows:      data,
		RowsCount: rowCount,
	}
}
