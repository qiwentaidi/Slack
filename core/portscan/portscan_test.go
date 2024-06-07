package portscan

import (
	"context"
	"testing"
)

func TestWrapperTcpWithTimeout(t *testing.T) {
	// // Set up test data
	// network := "tcp"
	// address := "www.baidu.com:80"
	// timeout := time.Second * 3

	// // Call the function under test
	// conn, err := WrapperTcpWithTimeout(network, address, timeout)

	// // Check if the function returned an error
	// if err != nil {
	// 	t.Errorf("WrapperTcpWithTimeout() returned an error: %v", err)
	// }

	// // Check if the returned connection is nil
	// if conn == nil {
	// 	t.Error("WrapperTcpWithTimeout() returned a nil connection")
	// }

	// // Close the connection
	// conn.Close()
	// TcpScan([]string{"118.31.42.197"}, []int{21, 22}, 10)
	PortBrute(context.TODO(), "167.172.69.161", []string{"root"}, []string{"123456"})
}
