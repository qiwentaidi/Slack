package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"regexp"
	"slack-wails/core/portscan"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/ssh"
)

type Tools struct{}

var (
	FscanRegs = map[string][]string{
		"FTP":       {"[+] ftp"},
		"SSH":       {"[+] SSH"},
		"Mssql":     {"[+] mssql"},
		"Oracle":    {"[+] oracle"},
		"Mysql":     {"[+] mysql"},
		"RDP":       {"[+] RDP"},
		"Redis":     {"[+] Redis"},
		"Postgres":  {"[+] Postgres"},
		"Mongodb":   {"[+] Mongodb"},
		"Memcached": {"[+] Memcached"},
		"MS17-010":  {"[+] MS17-010"},
		"DC INFO":   {"[+]DC"},
		"POC":       {"poc"},
		// 其他信息
		"Web INFO":        {"[*] WebTitle"},
		"INFO":            {"[+] InfoScan"},
		"Hikvison Camera": {"len:2512", "len:600", "len:481", "len:480"},
	}
	NetInfoReg = regexp.MustCompile(`\[\*]((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}(\s+\[\-\>](\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|[a-zA-Z0-9\-]+))+`)
)

func (t *Tools) FormatOutput(content string) map[string][]string {
	var result = make(map[string][]string)
	nets := NetInfoReg.FindAllString(content, -1)
	if len(nets) > 0 {
		var multiNetInfo []string
		for _, net := range FormatNetInfo(nets) {
			if len(net.IPs) >= 1 {
				multiNetInfo = append(multiNetInfo, fmt.Sprintf("%s: \n%s\n", net.Hostname, strings.Join(net.IPs, "\n")))
			}
		}
		result["NetInfo"] = multiNetInfo
	}
	lines := strings.Split(content, "\n")
	for name, reg := range FscanRegs {
		match := MatchLines(name, reg, lines)
		if len(match) > 0 {
			result[name] = match
		}
	}
	return result
}

type NetInfo struct {
	Hostname string
	IPs      []string
}

func FormatNetInfo(nets []string) []NetInfo {
	var netinfo []NetInfo
	for _, net := range nets {
		lines := strings.Split(net, "\n")
		var hostname string
		var IPs []string
		for index, line := range lines {
			line = strings.Replace(line, "   [->]", "", -1)
			if index == 1 {
				hostname = line
			} else if index > 1 {
				IPs = append(IPs, line)
			}
		}
		netinfo = append(netinfo, NetInfo{
			Hostname: hostname,
			IPs:      IPs,
		})
	}
	return netinfo
}

func MatchLines(name string, contains, lines []string) []string {
	var temp []string
	for _, v := range lines {
		for _, c := range contains {
			if strings.Contains(strings.ToLower(v), strings.ToLower(c)) {
				temp = append(temp, v)
			}
		}
	}
	return temp
}

func (t *Tools) ConnectAndExecute(protocol, ip, port string, username, password string) string {
	var commond, result string
	var err error
	host := fmt.Sprintf("%s:%s", ip, port)
	switch strings.ToLower(protocol) {
	case "ftp":
		commond = "ls"
		result, err = executeFtp(host, username, password)
	case "ssh":
		commands := []string{"whoami", "id", "ip a"}
		for _, cmd := range commands {
			output, err := executeSshCommand(host, username, password, cmd)
			if err == nil {
				result += fmt.Sprintf("[Commond] %s\n%s\n", cmd, output)
			}
		}
		return result
	case "mysql":
		commond = "SHOW DATABASES"
		result, err = executeMysqlQuery(host, username, password, commond)
	case "mssql":
		commond = "SELECT name FROM sys.databases"
		result, err = executeMssqlQuery(ip, port, username, password, commond)
	case "oracle":
		commond = "SELECT USERNAME FROM ALL_USERS"
		result, err = executeOracleQuery(host, username, password, commond)
	case "postgres":
		commond = "SELECT schema_name FROM information_schema.schemata"
		result, err = executePostgresQuery(host, username, password, commond)
	case "redis":
		commond = "info"
		result, err = executeRedisInfo(host, password)
	case "memcached":
		commond = "stats"
		result, err = executeMemcached(host)
	case "mongodb":
		commond = "show databases;"
		result, err = executeMongodbQuery(host)
	}
	if err != nil {
		return fmt.Sprintf("[Error] %v", err)
	}
	return fmt.Sprint("[Commond] " + commond + "\n" + result)
}

func executeSshCommand(host, username, password, command string) (string, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout: 5 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	// Create a session to execute the command
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Run the command and capture the output
	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %v", err)
	}

	return string(output), nil
}

func executeMongodbQuery(host string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create MongoDB URI
	mongoURI := fmt.Sprintf("mongodb://%s", host)

	// Define client options with or without credentials
	clientOpts := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, clientOpts)

	// Get the total number of databases
	if client == nil {
		return "", fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	databaseNames, err := client.ListDatabaseNames(context.TODO(), bson.D{})
	if err != nil {
		return "", fmt.Errorf("can't list MongoDB database names: %v", err)
	}

	client.Disconnect(context.TODO())
	return strings.Join(databaseNames, ","), nil
}

func executeMysqlQuery(host, username, password, command string) (string, error) {
	var databaseNames []string
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v)/mysql?charset=utf8&timeout=%v", username, password, host, 5*time.Second)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return "", fmt.Errorf("[mysql] 连接数据库失败: %v", err)
	}
	databases, err := db.Query(command)
	if err != nil {
		return "", fmt.Errorf("[mysql] 查询数据库失败: %v", err)
	}
	defer db.Close()
	for databases.Next() {
		var dbName string
		if err := databases.Scan(&dbName); err != nil {
			continue
		}
		databaseNames = append(databaseNames, dbName)
	}
	return strings.Join(databaseNames, ","), nil
}

func executeMssqlQuery(ip, port, username, password, command string) (string, error) {
	var databaseNames []string
	dataSourceName := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%v;encrypt=disable;timeout=%v", ip, username, password, port, 5*time.Second)
	db, err := sql.Open("mssql", dataSourceName)
	if err != nil {
		return "", fmt.Errorf("[sqlserver] 连接数据库失败: %v", err)
	}
	databases, err := db.Query(command)
	if err != nil {
		return "", fmt.Errorf("[sqlserver] 查询数据库失败: %v", err)
	}
	defer db.Close()
	for databases.Next() {
		var dbName string
		if err := databases.Scan(&dbName); err != nil {
			continue
		}
		databaseNames = append(databaseNames, dbName)
	}
	return strings.Join(databaseNames, ","), nil
}

func executeOracleQuery(host, username, password, command string) (string, error) {
	var databaseNames []string
	dataSourceName := fmt.Sprintf("oracle://%s:%s@%s/orcl", username, password, host)
	db, err := sql.Open("oracle", dataSourceName)
	if err != nil {
		return "", fmt.Errorf("[oracle] 连接数据库失败: %v", err)
	}
	databases, err := db.Query(command)
	if err != nil {
		return "", fmt.Errorf("[oracle] 查询数据库失败: %v", err)
	}
	defer db.Close()
	for databases.Next() {
		var dbName string
		if err := databases.Scan(&dbName); err != nil {
			continue
		}
		databaseNames = append(databaseNames, dbName)
	}
	return strings.Join(databaseNames, ","), nil
}

func executePostgresQuery(host, username, password, command string) (string, error) {
	var databaseNames []string
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v/postgres?sslmode=disable", username, password, host)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return "", fmt.Errorf("[postgres] 连接数据库失败: %v", err)
	}
	databases, err := db.Query(command)
	if err != nil {
		return "", fmt.Errorf("[postgres] 查询数据库失败: %v", err)
	}
	defer db.Close()
	for databases.Next() {
		var dbName string
		if err := databases.Scan(&dbName); err != nil {
			continue
		}
		databaseNames = append(databaseNames, dbName)
	}
	return strings.Join(databaseNames, ","), nil
}

func executeMemcached(host string) (string, error) {
	client, err := portscan.WrapperTcpWithTimeout("tcp", host, 5*time.Second)
	defer func() {
		if client != nil {
			client.Close()
		}
	}()
	if err == nil {
		err = client.SetDeadline(time.Now().Add(10 * time.Second))
		if err == nil {
			_, err = client.Write([]byte("stats\n")) //Set the key randomly to prevent the key on the server from being overwritten
			if err == nil {
				rev := make([]byte, 1024)
				n, err := client.Read(rev)
				if err == nil {
					return string(rev[:n]), nil
				} else {
					return "", errors.New("can't read" + host + "stat response")
				}
			}
		}
	}
	return "", errors.New("memcached://" + host + " is no unauthorized access")
}

func executeFtp(host, username, password string) (string, error) {
	var result string
	conn, err := ftp.Dial(host, ftp.DialWithTimeout(5*time.Second))
	if err == nil {
		err = conn.Login(username, password)
		if err == nil {
			dirs, err := conn.List("")
			if err == nil {
				if len(dirs) > 0 {
					for i := 0; i < len(dirs); i++ {
						if len(dirs[i].Name) > 50 {
							result += "\n" + dirs[i].Name[:50]
						} else {
							result += "\n" + dirs[i].Name
						}
						if i == 5 {
							break
						}
					}
				}
			}
		}
		return result, err
	}
	return result, err
}

func executeRedisInfo(host, password string) (string, error) {
	// Establish a TCP connection with a timeout
	conn, err := portscan.WrapperTcpWithTimeout("tcp", host, 5*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Set the read deadline for the connection
	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return "", err
	}

	// Check if the connection is unauthorized or requires authentication
	if password == "unauthorized" {
		return sendCommand(conn, "info\r\n")
	}

	// Authenticate with the given password
	authResponse, err := sendCommand(conn, fmt.Sprintf("auth %s\r\n", password))
	if err != nil {
		return "", err
	}

	// Check if the authentication was successful
	if !strings.Contains(authResponse, "+OK") {
		return "", errors.New("password is incorrect")
	}

	// Send the INFO command after successful authentication
	return sendCommand(conn, "info\r\n")
}

// sendCommand sends a command to the Redis server and returns the response
func sendCommand(conn net.Conn, command string) (string, error) {
	_, err := conn.Write([]byte(command))
	if err != nil {
		return "", err
	}

	// Read the response from the server
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:n]), nil
}
