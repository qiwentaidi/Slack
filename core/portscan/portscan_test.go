package portscan

import (
	"fmt"
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

	// synscan --------------------
	// var id int32
	// SynScan(context.Background(), "1.1.1.1", []uint16{53, 80, 443, 8080}, &id)

	// mongodb connect

	flag, err := MongodbConn("192.168.0.107:27017", "", "")
	fmt.Printf("flag1: %v\n", flag)
	fmt.Printf("err: %v\n", err)
}
