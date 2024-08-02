package dirsearch

import (
	"bytes"
	"context"
	"net/http"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/util"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	ExitFunc      = false
	bodyLengthMap map[string]int
	mutex         = sync.Mutex{}
)

type Result struct {
	Status   int
	URL      string
	Location string
	Length   int
	Message  string
}

type Options struct {
	Method                 string
	URL                    string
	Paths                  []string
	Workers                int
	Timeout                int
	BodyExclude            string
	BodyLengthExcludeTimes int // 过滤响应包长度数据次数出现了多少次，就过滤相同的
	StatusCodeExclude      []int
	FailedCounts           int32
	Redirect               bool
	Interval               int
	CustomHeader           string
}

// method 请求类型
func NewScanner(ctx context.Context, o Options) {
	var urls []string
	if util.IsURL(o.URL) {
		urls = append(urls, o.URL)
	} else {
		urls = util.ReadLine(o.URL)
	}
	runtime.EventsEmit(ctx, "dirsearchCounts", len(urls)*len(o.Paths))
	bodyLengthMap = make(map[string]int)
	// 初始化请求信息
	if o.Timeout == 0 {
		o.Timeout = 8
	}
	client := clients.NotFollowClient()
	if o.Redirect {
		client = clients.DefaultClient()
	}
	var header = http.Header{}
	if o.CustomHeader != "" {
		for _, single := range strings.Split(o.CustomHeader, "\n") {
			temp := strings.Split(single, ":")
			header.Set(temp[0], temp[1])
		}
	}
	var id int32
	count := len(o.Paths) * len(urls)
	single := make(chan struct{})
	retChan := make(chan Result, count)
	var wg sync.WaitGroup
	go func() {
		for pr := range retChan {
			runtime.EventsEmit(ctx, "dirsearchLoading", pr)
		}
		close(single)
		runtime.EventsEmit(ctx, "dirsearchComplete", "done")
	}()

	dirScan := func(url string) {
		r := Scan(ctx, url, header, o, client)
		atomic.AddInt32(&id, 1)
		runtime.EventsEmit(ctx, "dirsearchProgressID", id)
		retChan <- r
	}
	threadPool, _ := ants.NewPoolWithFunc(o.Workers, func(p interface{}) {
		path := p.(string)
		dirScan(path)
		wg.Done()
	})
	defer threadPool.Release()
	for _, url := range urls {
		if !strings.HasSuffix(url, "/") {
			url = url + "/"
		}
		for _, path := range o.Paths {
			if ExitFunc {
				return
			}
			if strings.HasPrefix(path, "/") {
				path = strings.TrimLeft(path, "/")
			}
			wg.Add(1)
			threadPool.Invoke(url + path)
			if o.Interval != 0 {
				time.Sleep(time.Second * time.Duration(o.Interval))
			}
		}
	}
	wg.Wait()
	close(retChan)
	<-single
}

// status 1 表示被排除显示在外，不计入前端ERROR请求中
func Scan(ctx context.Context, url string, header http.Header, o Options, client *http.Client) Result {
	var result Result
	resp, body, err := clients.NewRequest(o.Method, url, header, nil, o.Timeout, true, client)
	if err != nil {
		gologger.IntervalError(ctx, err)
		return result
	}
	if o.BodyExclude != "" && bytes.Contains(body, []byte(o.BodyExclude)) {
		result.Status = 1
		return result
	} else {
		result.Status = resp.StatusCode
	}
	if util.ArrayContains(result.Status, o.StatusCodeExclude) {
		result.Status = 1
		return result
	}
	result.Length = len(body)
	// 记录同一状态码下长度出现次数，当次数超过o.BodyLengthExcludeTimes时，将状态码设置为1，过滤显示
	mutex.Lock()
	bodyLengthMap[resp.Status]++
	if bodyLengthMap[resp.Status] > o.BodyLengthExcludeTimes {
		result.Status = 1
	}
	mutex.Unlock()
	result.Location = resp.Header.Get("Location")
	result.URL = url
	return result
}
