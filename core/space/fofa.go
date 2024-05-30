package space

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slack-wails/lib/clients"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

const TipApi = "https://api.fofa.info/v1/search/tip?"

type FofaConfig struct {
	AppId      string
	PrivateKey string
}

type TipsResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []Data `json:"data"`
}

type Data struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	RCode   string `json:"r_code"`
}

func (f *FofaConfig) GetTips(key string) ([]byte, error) {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signParam := "q" + key + "ts" + ts
	params := url.Values{}
	params.Set("q", key)
	params.Set("ts", ts)
	params.Set("sign", f.GetInputSign(signParam))
	params.Set("app_id", f.AppId)
	resp, err := http.Get(TipApi + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return b, nil
}

func (f *FofaConfig) GetInputSign(inputString string) string {
	data := []byte(inputString)
	keyBytes, err := base64.StdEncoding.DecodeString(f.PrivateKey)
	if err != nil {
		return ""
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(keyBytes)
	if err != nil {
		return ""
	}
	hash := sha256.New()
	hash.Write(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, hash.Sum(nil))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(signature)
}

// fofa base64加密接口
func FOFABaseEncode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

type FofaResult struct {
	Error   bool       `json:"error"`
	Errmsg  string     `json:"errmsg"`
	Mode    string     `json:"mode"`
	Page    int64      `json:"page"`
	Query   string     `json:"query"`
	Results [][]string `json:"results"`
	Size    int64      `json:"size"`
}

type FofaSearchResult struct {
	Status  bool
	Message string
	Total   int64
	Results []struct {
		URL      string
		Title    string
		IP       string
		Port     string
		Domain   string
		Protocol string
		Country  string
		Region   string
		City     string
		ICP      string
	}
}

func FofaApiSearch(search, pageSize, pageNum, addr, email, key string, fraud, cert bool) *FofaSearchResult {
	var fs FofaSearchResult
	address := addr + "api/v1/search/all?email=" + email + "&key=" + key + "&qbase64=" +
		FOFABaseEncode(search) + "&cert.is_valid" + fmt.Sprint(cert) + fmt.Sprintf("&is_fraud=%v&is_honeypot=%v", fmt.Sprint(fraud), fmt.Sprint(fraud)) +
		"&page=" + pageNum + "&size=" + pageSize + "&fields=host,title,ip,domain,port,protocol,country_name,region,city,icp"
	_, b, err := clients.NewSimpleGetRequest(address, http.DefaultClient)
	if err != nil {
		fs.Status = false
		fs.Message = "请求失败"
	}
	var fr FofaResult
	if err = json.Unmarshal(b, &fr); err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
	}
	fs.Total = fr.Size
	if fr.Error {
		fs.Status = false
		fs.Message = fr.Errmsg
	} else {
		if fr.Size == 0 {
			fs.Message = "未查询到数据"
		} else {
			var temp string
			fs.Status = true
			for i := 0; i < len(fr.Results); i++ {
				if !strings.Contains(fr.Results[i][0], "://") && (fr.Results[i][5] == "http" || fr.Results[i][5] == "https") {
					temp = fr.Results[i][5] + "://" + fr.Results[i][0]
				} else {
					temp = fr.Results[i][0]
				}
				fs.Results = append(fs.Results, struct {
					URL      string
					Title    string
					IP       string
					Port     string
					Domain   string
					Protocol string
					Country  string
					Region   string
					City     string
					ICP      string
				}{
					URL:      temp,
					Title:    fr.Results[i][1],
					IP:       fr.Results[i][2],
					Port:     fr.Results[i][4],
					Domain:   fr.Results[i][3],
					Protocol: fr.Results[i][5],
					Country:  fr.Results[i][6],
					Region:   fr.Results[i][7],
					City:     fr.Results[i][8],
					ICP:      fr.Results[i][9],
				})
			}
		}
	}
	time.Sleep(time.Second)
	return &fs
}
