package main

import (
	"bufio"
	"encoding/base64"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slack-wails/lib/update"
	"slack-wails/lib/util"
	"strings"
)

// File struct 文件操作
type File struct {
	configPath string
}

func NewFile() *File {
	return &File{
		configPath: util.HomeDir() + "/slack/config/",
	}
}

// 开始就要检测
func (f *File) UserHomeDir() string {
	return util.HomeDir()
}

func (f *File) PathBase(p string) string {
	return filepath.Base(p)
}

func (f *File) OpenFolder(path string) string {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err.Error()
	}
	return ""
}

func (f *File) CheckFileStat(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func (f *File) GetFileContent(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "文件不存在"
	}
	return string(b)
}

func (f *File) UpdatePocFile() string {
	if err := update.UpdatePoc(f.configPath); err != nil {
		return err.Error()
	}
	return ""
}

func (f *File) InitConfig() bool {
	return update.InitConfig(f.configPath)
}

func (*File) InitMemo(filepath, content string) bool {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return false
	}
	_, err = f.WriteString(content)
	return err == nil
}

func (*File) ReadMemo(filepath string) map[string]string {
	file, err := os.Open(filepath)
	if err != nil {
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var key string
	var value strings.Builder
	keyValueMap := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// This is a key line
			if key != "" {
				// Save the previous key-value pair
				keyValueMap[key] = value.String()
				value.Reset()
			}
			key = line[1 : len(line)-1] // Remove brackets
		} else {
			// This is a value line
			value.WriteString(line + "\n")
		}
	}
	// Save the last key-value pair
	if key != "" {
		keyValueMap[key] = value.String()
	}
	return keyValueMap
}

func (*File) Mkdir(dir string) bool {
	return os.Mkdir(dir, 0777) == nil
}

func (*File) WriteFile(filetype, path, content string) bool {
	var buf []byte
	switch filetype {
	case "base64":
		buf, _ = base64.StdEncoding.DecodeString(content)
	// txt
	default:
		buf = []byte(content)
	}
	err := os.WriteFile(path, buf, 0644)
	return err == nil
}

// type Records struct {
// 	Fields []string
// }

// func (*File) HunterRemoveDuplicates(filename string) bool {
// 	// 打开 CSV 文件
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		fmt.Printf("Cannot access file %s: %v\n", filename, err)
// 		return false
// 	}
// 	defer file.Close()

// 	// 创建 CSV 读取器
// 	reader := csv.NewReader(file)

// 	// 读取 CSV 文件头
// 	headers, err := reader.Read()
// 	if err != nil {
// 		fmt.Printf("Error reading headers: %v\n", err)
// 		return false
// 	}
// 	// 使用 map 去重
// 	urlRecords := make(map[string]Records)
// 	for {
// 		record, err := reader.Read()
// 		if err != nil {
// 			if err.Error() == "EOF" {
// 				break
// 			}
// 			fmt.Printf("Error reading record: %v\n", err)
// 		}
// 		// 获取URL字段进行初步去重
// 		url := record[4]
// 		urls = append(urls, url)
// 		if _, exists := urlRecords[url]; !exists {
// 			urlRecords[url] = Records{
// 				Fields: record,
// 			}
// 		}
// 	}

// 	// 第二轮去重，按 ip-port-title
// 	uniqueRecords := make(map[string]Records)
// 	for _, record := range urlRecords {
// 		fields := record.Fields
// 		ip, port, title := fields[0], fields[1], fields[5]
// 		key := fmt.Sprintf("%s-%s-%s", ip, port, title)
// 		ips = append(ips, ip)
// 		if _, exists := uniqueRecords[key]; !exists {
// 			uniqueRecords[key] = Records{
// 				Fields: fields,
// 			}
// 		}
// 	}

// 	// 创建输出文件
// 	outFile, err := os.Create(getOutputFilename(filename))
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 		return false
// 	}
// 	defer outFile.Close()

// 	// 创建 CSV 写入器
// 	writer := csv.NewWriter(outFile)
// 	defer writer.Flush()

// 	// 写入 CSV 文件头
// 	if err := writer.Write(headers); err != nil {
// 		fmt.Printf("Error writing headers: %v\n", err)
// 		return false
// 	}

// 	// 写入唯一记录
// 	for _, record := range uniqueRecords {
// 		if err := writer.Write(record.Fields); err != nil {
// 			fmt.Printf("Error writing record: %v\n", err)
// 			return false
// 		}
// 	}
// 	return true
// }

// func getOutputFilename(inputFile string) string {
// 	ext := filepath.Ext(inputFile)
// 	name := strings.TrimSuffix(inputFile, ext)
// 	return fmt.Sprintf("%s_removeDuplicates%s", name, ext)
// }
