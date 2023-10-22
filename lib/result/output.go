package result

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"

	"slack/common"
	"slack/common/logger"
	"slack/gui/custom"
	"slack/gui/global"
	"slack/lib/util"
	"strings"
	"time"

	"fyne.io/fyne/v2/dialog"
)

// mode为结果类型
func ArrayOutput(data [][]string, mode string) {
	if _, err := os.Stat("./reports"); err != nil {
		os.Mkdir("./reports", 0664)
	}
	filename := mode + "-" + time.Now().Format("20060102_150405") + ".csv"
	file, _ := os.OpenFile("./reports/"+filename, os.O_CREATE|os.O_RDWR, os.ModePerm) // 创建结果文件
	file.WriteString("\xEF\xBB\xBF")
	scanFile := csv.NewWriter(file)
	for _, line := range data {
		scanFile.Write(line)
		scanFile.Flush()
	}
	dialog.ShowInformation("", "结果保存完毕!", global.Win)
}

func DistributeTaskReport(companyname, filename string, targets []string) {
	var temp []string
	err := os.MkdirAll("./reports/distribute/"+companyname, 0777)
	if err != nil {
		logger.Error(err)
	}
	file_name := fmt.Sprintf("./reports/distribute/%v/%v.txt", companyname, filename)
	if _, err2 := os.Stat(file_name); err2 == nil {
		b, _ := os.ReadFile(file_name)
		temp = append(temp, strings.Split(string(b), "\n")...)
		custom.Console.Append("[INF] " + file_name + ".txt文件已存在，若存在新内容会自动追加\n")
	}
	f, err := os.OpenFile(file_name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		logger.Error(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		temp = append(temp, s.Text())
	}
	temp = append(temp, targets...)
	write := bufio.NewWriter(f)
	dt := util.RemoveDuplicates[string](temp)
	for i, t := range dt {
		if t != "" && filename == "ips" {
			common.WaitPortScan = append(common.WaitPortScan, t)
		}
		if i != len(dt) {
			write.WriteString(t + "\n")
		} else {
			write.WriteString(t)
		}
	}
	write.Flush()
}
