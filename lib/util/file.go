package util

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
)

func BufferWriteAppend(filename string, param string) error {
	fileHandle, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0660)
	if err != nil {
		return err
	}
	defer fileHandle.Close()
	// NewWriter 默认缓冲区大小是 4096
	// 需要使用自定义缓冲区的writer 使用 NewWriterSize()方法
	buf := bufio.NewWriter(fileHandle)
	// 字节写入
	//buf.Write([]byte(param))
	// 字符串写入
	buf.WriteString(param + "\n")
	// 将缓冲中的数据写入
	return buf.Flush()
}

func OpenFolder(path string) {
	var command string
	switch runtime.GOOS {
	case "linux":
		command = "xdg-open"
	case "windows":
		command = "explorer"
	case "darwin":
		command = "open"
	}
	cmd := exec.Command(command, path)
	cmd.Start()
}
