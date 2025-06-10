package services

import (
	"context"
	"database/sql"

	"github.com/xuri/excelize/v2"

	"encoding/json"
	"fmt"
	"os"
	"slack-wails/lib/gologger"
	"slack-wails/lib/report"
	"slack-wails/lib/structs"
	"slack-wails/lib/utils"
	"slack-wails/lib/utils/fileutil"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	ctx           context.Context
	DB            *sql.DB // 系统数据库
	lock          sync.RWMutex
	OtherDatabase *sql.DB                     // 数据库信息采集时的连接池
	MongoClient   *mongo.Client               // mongodb连接池
	PostgresInfo  *structs.DatabaseConnection // 用于临时存储postgres数据库连接信息，方便其他方法调用
}

func (d *Database) Startup(ctx context.Context) {
	d.ctx = ctx
}

func NewDatabase() *Database {
	os.Mkdir(utils.HomeDir()+"/slack", 0777) // 创建配置文件夹
	dp := utils.HomeDir() + "/slack/config.db"
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

// SQLite 检查字段是否已存在
func columnExists(db *sql.DB, tableName, columnName string) bool {
	query := `PRAGMA table_info(` + tableName + `);`
	rows, _ := db.Query(query)
	defer rows.Close()

	for rows.Next() {
		var (
			cid        int
			name       string
			fieldType  string
			notnull    int
			dflt_value interface{}
			pk         int
		)
		rows.Scan(&cid, &name, &fieldType, &notnull, &dflt_value, &pk)
		if name == columnName {
			return true
		}
	}
	return false
}
func (d *Database) CreateTable() bool {
	_, err := d.DB.Exec(`
		CREATE TABLE IF NOT EXISTS windows_size (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			width INTEGER,
			height INTEGER
		);
        CREATE TABLE IF NOT EXISTS hunter_syntax ( name TEXT, content TEXT );
        CREATE TABLE IF NOT EXISTS quake_syntax ( name TEXT, content TEXT );
        CREATE TABLE IF NOT EXISTS fofa_syntax ( name TEXT, content TEXT );
        CREATE TABLE IF NOT EXISTS agent_pool ( hosts TEXT );
        CREATE TABLE IF NOT EXISTS dirsearch ( path TEXT, times INTEGER );
        CREATE TABLE IF NOT EXISTS dbManager ( nanoid TEXT, scheme TEXT, host TEXT, port INTEGER, username TEXT, password TEXT, notes TEXT );
        CREATE TABLE IF NOT EXISTS scanTask ( task_id TEXT PRIMARY KEY, task_name TEXT, targets TEXT, failed INTEGER, vulnerability INTEGER );
        CREATE TABLE IF NOT EXISTS FingerprintInfo ( task_id TEXT, url TEXT, status INTEGER, length INTEGER, title TEXT, detect TEXT, is_waf INTEGER, waf TEXT, fingerprints TEXT, screenshot TEXT, host TEXT, scheme TEXT, port INTEGER );
        CREATE TABLE IF NOT EXISTS VulnerabilityInfo ( task_id TEXT, template_id TEXT, vuln_name TEXT, protocol TEXT, severity TEXT, vuln_url TEXT, extract TEXT, request TEXT, response TEXT, description TEXT, reference TEXT, response_time TEXT );
    `)
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[sqlite] create table: %s", err))
		return false
	}
	// 当数据不存在时，插入一条默认记录
	_, err = d.DB.Exec(`INSERT OR IGNORE INTO windows_size (id, width, height) VALUES (1, ?, ?)`, defaultWindowsWidth, defaultWindowsHeight)
	if err != nil {
		gologger.Debug(d.ctx, fmt.Sprintf("[sqlite] insert default windows_size: %s", err))
		return false
	}

	if !columnExists(d.DB, "FingerprintInfo", "host") {
		_, err := d.DB.Exec(`ALTER TABLE FingerprintInfo ADD COLUMN host TEXT`)
		if err != nil {
			return false
		}
	}
	if !columnExists(d.DB, "FingerprintInfo", "scheme") {
		_, err := d.DB.Exec(`ALTER TABLE FingerprintInfo ADD COLUMN scheme TEXT`)
		if err != nil {
			return false
		}
	}
	if !columnExists(d.DB, "FingerprintInfo", "port") {
		_, err := d.DB.Exec(`ALTER TABLE FingerprintInfo ADD COLUMN port INTEGER`)
		if err != nil {
			return false
		}
	}
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

const defaultWindowsWidth = 1280
const defaultWindowsHeight = 800

func (d *Database) SelectWindowsSize() structs.WindowsSize {
	var w, h int
	rows, err := d.DB.Query("SELECT width, height FROM windows_size WHERE id = 1")
	if err != nil {
		return structs.WindowsSize{
			Width:  defaultWindowsWidth,
			Height: defaultWindowsHeight,
		}
	}
	for rows.Next() {
		rows.Scan(&w, &h)
	}
	return structs.WindowsSize{
		Width:  w,
		Height: h,
	}
}

func (d *Database) SaveWindowsScreenSize(width, height int) bool {
	return d.ExecSqlStatement("UPDATE windows_size SET width = ?, height = ? WHERE id = 1", width, height)
}

func (d *Database) SelectAllAgentPool() (hosts []string) {
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

func (d *Database) InsertAgentPool(host string) bool {
	insertStmt := "INSERT INTO agent_pool(hosts) VALUES(?)"
	return d.ExecSqlStatement(insertStmt, host)
}

func (d *Database) DeleteAgentPool(host string) bool {
	deleteStmt := "DELETE FROM agent_pool WHERE hosts = ?"
	return d.ExecSqlStatement(deleteStmt, host)
}

func (d *Database) DeleteAllAgentPool() bool {
	deleteStmt := "DELETE FROM agent_pool"
	return d.ExecSqlStatement(deleteStmt)
}

func (d *Database) SelectAllSyntax(module string) (data []structs.SpaceEngineSyntax) {
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
		data = append(data, structs.SpaceEngineSyntax{
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

func (d *Database) GetAllPathsAndTimes() []structs.PathTimes {
	d.lock.RLock()         // 读操作前加锁
	defer d.lock.RUnlock() // 函数结束时解锁

	rows, err := d.DB.Query("SELECT path, times FROM dirsearch")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var results []structs.PathTimes

	for rows.Next() {
		var ds structs.PathTimes
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

func (d *Database) DeleteRecordByPath(path string) bool {
	return d.ExecSqlStatement("DELETE FROM dirsearch WHERE path = ?", path)
}

// 执行 SQL 删除语句，删除 times 为 1 的记录
func (d *Database) DeleteRecordsWithTimesEqualOne() bool {
	return d.ExecSqlStatement("DELETE FROM dirsearch WHERE times = 1")
}

// 检索所有扫描记录
func (d *Database) RetrieveAllScanTasks() []structs.TaskResult {
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

// 添加任务记录
func (d *Database) AddScanTask(taskid, taskname, targets string, failed, vulnerability int) bool {
	insertStmt := "INSERT INTO scanTask (task_id, task_name, targets, failed, vulnerability) VALUES (?, ?, ?, ?, ?)"
	return d.ExecSqlStatement(insertStmt, taskid, taskname, targets, failed, vulnerability)
}

// 修改扫描结果 - 失败数量，漏洞数量
func (d *Database) UpdateScanTaskWithResults(taskid string, failed, vulnerability int) bool {
	updateStmt := "UPDATE scanTask SET failed = ?, vulnerability = ? WHERE task_id = ?"
	return d.ExecSqlStatement(updateStmt, failed, vulnerability, taskid)
}

// 移除扫描记录
func (d *Database) RemoveScanTask(taskid string) bool {
	deleteStmt := "DELETE FROM scanTask WHERE task_id = ?"
	isSuccess := d.ExecSqlStatement(deleteStmt, taskid)
	if isSuccess {
		d.ExecSqlStatement("DELETE FROM FingerprintInfo WHERE task_id = ?", taskid)
		d.ExecSqlStatement("DELETE FROM VulnerabilityInfo WHERE task_id = ?", taskid)
	}
	return isSuccess
}

// 重命名任务
func (d *Database) RenameScanTask(taskid, taskname string) bool {
	updateStmt := "UPDATE scanTask SET task_name = ? WHERE task_id = ?"
	return d.ExecSqlStatement(updateStmt, taskname, taskid)
}

// 根据taskid检索指纹扫描的结果
func (d *Database) RetrieveFingerscanResults(taskid string) []structs.InfoResult {
	rows, err := d.DB.Query("SELECT * FROM FingerprintInfo WHERE task_id = ?;", taskid)
	if err != nil {
		gologger.Debug(d.ctx, err)
		return []structs.InfoResult{}
	}
	defer rows.Close()
	var results []structs.InfoResult
	for rows.Next() {
		var result structs.InfoResult
		var fingerprintsStr string
		var task_id string
		var host *string // 使用指针来处理可能的 NULL 值
		var scheme *string
		var port *int
		err = rows.Scan(&task_id, &result.URL, &result.StatusCode, &result.Length, &result.Title, &result.Detect, &result.IsWAF, &result.WAF, &fingerprintsStr, &result.Screenshot, &host, &scheme, &port)
		if err != nil {
			gologger.Debug(d.ctx, err)
			continue
		}
		if fingerprintsStr != "" {
			if strings.Contains(fingerprintsStr, ",") {
				result.Fingerprints = strings.Split(fingerprintsStr, ",")
			} else {
				result.Fingerprints = []string{fingerprintsStr}
			}
		} else {
			result.Fingerprints = []string{}
		}
		if port != nil {
			result.Port = *port
		}
		if host != nil {
			result.Host = *host
		}
		if scheme != nil {
			result.Scheme = *scheme
		}
		results = append(results, result)
	}
	return results
}

// 根据taskid检索漏洞扫描记录
func (d *Database) RetrievePocscanResults(taskid string) []structs.VulnerabilityInfo {
	rows, err := d.DB.Query("SELECT * FROM VulnerabilityInfo WHERE task_id = ?", taskid)
	if err != nil {
		return []structs.VulnerabilityInfo{}
	}
	defer rows.Close()
	var results []structs.VulnerabilityInfo
	for rows.Next() {
		var result structs.VulnerabilityInfo
		var task_id string
		var responseTime *string // 使用指针来处理可能的 NULL 值
		err = rows.Scan(&task_id, &result.ID, &result.Name, &result.Type, &result.Severity, &result.URL, &result.Extract, &result.Request, &result.Response, &result.Description, &result.Reference, &responseTime)
		if err != nil {
			gologger.Debug(d.ctx, err)
			continue
		}
		if responseTime != nil {
			result.ResponseTime = *responseTime // 只有在 responseTime 不为 NULL 时才赋值
		}
		results = append(results, result)
	}
	return results
}

// 添加指纹扫描结果
func (d *Database) AddFingerscanResult(result structs.InfoResult) bool {
	insertStmt := "INSERT INTO FingerprintInfo (task_id, url, status, length, title, detect, is_waf, waf, fingerprints, screenshot, host, scheme, port) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return d.ExecSqlStatement(insertStmt, result.TaskId, result.URL, result.StatusCode, result.Length, result.Title, result.Detect, result.IsWAF, result.WAF, strings.Join(result.Fingerprints, ","), result.Screenshot, result.Host, result.Scheme, result.Port)
}

// 添加漏洞扫描结果
func (d *Database) AddPocscanResult(result structs.VulnerabilityInfo) bool {
	insertStmt := "INSERT INTO VulnerabilityInfo (task_id, template_id, vuln_name, protocol, severity, vuln_url, extract, request, response, description, reference, response_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	return d.ExecSqlStatement(insertStmt, result.TaskId, result.ID, result.Name, result.Type, result.Severity, result.URL, result.Extract, result.Request, result.Response, result.Description, result.Reference, result.ResponseTime)
}

// 移除某个漏洞
func (d *Database) RemovePocscanResult(taskid, template_id, vuln_url string) bool {
	deleteStmt := "DELETE FROM VulnerabilityInfo WHERE task_id = ? AND template_id = ? AND vuln_url = ?"
	return d.ExecSqlStatement(deleteStmt, taskid, template_id, vuln_url)
}

// 移除某组指纹信息，用于删除探测http的基本状态，后续会由指纹探测重新写入
func (d *Database) RemoveFingerprintResult(taskid string, link []string) bool {
	// 如果链接列表为空，直接返回 false 表示操作未执行
	if len(link) == 0 {
		gologger.Info(d.ctx, "No link provided to remove fingerprint result")
		return true
	}

	// 构造占位符和参数列表
	placeholders := make([]string, len(link))
	params := make([]interface{}, len(link)+1)
	params[0] = taskid
	for i, l := range link {
		placeholders[i] = "?"
		params[i+1] = l
	}

	// 构造 SQL 语句
	deleteStmt := fmt.Sprintf(
		"DELETE FROM FingerprintInfo WHERE task_id = ? AND url IN (%s)",
		strings.Join(placeholders, ","),
	)

	// 执行 SQL 语句
	return d.ExecSqlStatement(deleteStmt, params...)
}

// 导出JSON报告
func (d *Database) ExportWebReportWithJson(reportpath string, tasks []structs.TaskResult) bool {
	var fingerprintsResults []structs.InfoResult
	var pocsResults []structs.VulnerabilityInfo
	var targets []string
	for _, task := range tasks {
		fingerprintsResult := d.RetrieveFingerscanResults(task.TaskId)
		pocsResult := d.RetrievePocscanResults(task.TaskId)
		fingerprintsResults = append(fingerprintsResults, fingerprintsResult...)
		pocsResults = append(pocsResults, pocsResult...)
		targets = append(targets, task.Targets)
	}
	result := structs.WebReport{
		Targets:      strings.Join(targets, "\n"),
		Fingerprints: fingerprintsResults,
		POCs:         pocsResults,
	}
	return fileutil.SaveJsonWithFormat(d.ctx, reportpath, result)
}

// 加载JSON报告
func (d *Database) ReadWebReportWithJson(reportpath string) (result structs.WebReport, err error) {
	data, err := os.ReadFile(reportpath)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &result)
	return
}

// 导出HTML报告
func (d *Database) ExportWebReportWithHtml(reportpath string, taskids []string) bool {
	var fingerprintsResults []structs.InfoResult
	var pocsResults []structs.VulnerabilityInfo
	for _, taskid := range taskids {
		fingerprintsResult := d.RetrieveFingerscanResults(taskid)
		pocsResult := d.RetrievePocscanResults(taskid)
		fingerprintsResults = append(fingerprintsResults, fingerprintsResult...)
		pocsResults = append(pocsResults, pocsResult...)
	}
	return os.WriteFile(reportpath, []byte(report.GenerateReport(fingerprintsResults, pocsResults)), 0644) == nil
}

// 导出EXCEL报告

func (d *Database) ExportWebReportWithExcel(reportpath string, tasks []structs.TaskResult) bool {
	var fingerprintsResults []structs.InfoResult
	var pocsResults []structs.VulnerabilityInfo
	var targets []string

	// 汇总任务数据
	for _, task := range tasks {
		fingerprintsResult := d.RetrieveFingerscanResults(task.TaskId)
		pocsResult := d.RetrievePocscanResults(task.TaskId)
		fingerprintsResults = append(fingerprintsResults, fingerprintsResult...)
		pocsResults = append(pocsResults, pocsResult...)
		targets = append(targets, task.Targets)
	}
	// 创建Excel文件
	f := excelize.NewFile()
	// 添加"Targets"工作表
	targetsSheet := "Targets"
	f.NewSheet(targetsSheet)
	f.SetCellValue(targetsSheet, "A1", "Targets")
	var alltargets []string
	for _, target := range targets {
		alltargets = append(alltargets, strings.Split(target, "\n")...)
	}
	for i, target := range alltargets {
		f.SetCellValue(targetsSheet, fmt.Sprintf("A%d", i+2), target)
	}
	f.DeleteSheet("Sheet1") // 删除默认的工作表
	// 添加"Fingerprints"工作表
	fingerprintsSheet := "Fingerprints"
	f.NewSheet(fingerprintsSheet)
	fingerprintsHeader := []string{"URL", "Scheme", "Host", "Port", "StatusCode", "Length", "Title", "Fingerprints", "IsWAF", "WAF", "Detect", "Screenshot"}
	for i, header := range fingerprintsHeader {
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("%s1", string(rune('A'+i))), header)
	}
	for i, result := range fingerprintsResults {
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("A%d", i+2), result.URL)
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("B%d", i+2), result.Scheme)
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("C%d", i+2), result.Host)
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("D%d", i+2), result.Port)
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("E%d", i+2), result.StatusCode)
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("F%d", i+2), result.Length)
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("G%d", i+2), result.Title)
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("H%d", i+2), strings.Join(result.Fingerprints, ","))
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("I%d", i+2), result.IsWAF)
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("J%d", i+2), result.WAF)
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("K%d", i+2), result.Detect)
		f.SetCellValue(fingerprintsSheet, fmt.Sprintf("L%d", i+2), result.Screenshot)
	}

	// 添加"POCs"工作表
	pocsSheet := "POCs"
	f.NewSheet(pocsSheet)
	pocsHeader := []string{"ID", "Name", "Description", "Reference", "Type", "Severity", "URL", "Request", "Response", "ResponseTime", "Extract"}
	for i, header := range pocsHeader {
		f.SetCellValue(pocsSheet, fmt.Sprintf("%s1", string(rune('A'+i))), header)
	}
	for i, result := range pocsResults {
		f.SetCellValue(pocsSheet, fmt.Sprintf("A%d", i+2), result.ID)
		f.SetCellValue(pocsSheet, fmt.Sprintf("B%d", i+2), result.Name)
		f.SetCellValue(pocsSheet, fmt.Sprintf("C%d", i+2), result.Description)
		f.SetCellValue(pocsSheet, fmt.Sprintf("D%d", i+2), result.Reference)
		f.SetCellValue(pocsSheet, fmt.Sprintf("E%d", i+2), result.Type)
		f.SetCellValue(pocsSheet, fmt.Sprintf("F%d", i+2), result.Severity)
		f.SetCellValue(pocsSheet, fmt.Sprintf("G%d", i+2), result.URL)
		f.SetCellValue(pocsSheet, fmt.Sprintf("H%d", i+2), result.Request)
		f.SetCellValue(pocsSheet, fmt.Sprintf("I%d", i+2), result.Response)
		f.SetCellValue(pocsSheet, fmt.Sprintf("J%d", i+2), result.ResponseTime)
		f.SetCellValue(pocsSheet, fmt.Sprintf("K%d", i+2), result.Extract)
	}

	// 保存Excel文件
	if err := f.SaveAs(reportpath); err != nil {
		gologger.Error(d.ctx, "Failed to save Excel file")
		return false
	}

	return true
}

func (d *Database) ExportJSReportWithExcel(reportpath string, results []structs.JSFindResult) bool {
	// 创建Excel文件
	f := excelize.NewFile()
	// 添加"Fingerprints"工作表
	resultSheet := "JSFinder Results"
	f.NewSheet(resultSheet)
	resultHeader := []string{"Source", "Method", "VulType", "Severity", "Length", "Request", "Response", "Filed"}
	for i, header := range resultHeader {
		f.SetCellValue(resultSheet, fmt.Sprintf("%s1", string(rune('A'+i))), header)
	}
	for i, result := range results {
		f.SetCellValue(resultSheet, fmt.Sprintf("A%d", i+2), result.Source)
		f.SetCellValue(resultSheet, fmt.Sprintf("B%d", i+2), result.Method)
		f.SetCellValue(resultSheet, fmt.Sprintf("C%d", i+2), result.VulType)
		f.SetCellValue(resultSheet, fmt.Sprintf("D%d", i+2), result.Severity)
		f.SetCellValue(resultSheet, fmt.Sprintf("E%d", i+2), result.Length)
		f.SetCellValue(resultSheet, fmt.Sprintf("F%d", i+2), result.Request)
		f.SetCellValue(resultSheet, fmt.Sprintf("G%d", i+2), result.Response)
		f.SetCellValue(resultSheet, fmt.Sprintf("H%d", i+2), result.Filed)
	}
	f.DeleteSheet("Sheet1") // 删除默认的工作表

	// 保存Excel文件
	if err := f.SaveAs(reportpath); err != nil {
		gologger.Error(d.ctx, "Failed to save Excel file : "+err.Error())
		return false
	}

	return true
}
