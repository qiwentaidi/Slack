package dirsearch

import (
	"bytes"
	"context"
	"slack-wails/lib/gologger"
	"slack-wails/lib/utils/arrayutil"
	"slack-wails/lib/utils/httputil"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qiwentaidi/clients"

	"github.com/go-resty/resty/v2"
	"github.com/panjf2000/ants/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var maxReponseLength = 1024 * 100 // 最大响应包长度，100kb

type Dirsearch struct {
	ctx           context.Context
	options       Options
	bodyLengthMap map[int]int
	mutex         sync.Mutex
	client        *resty.Client
	headers       map[string]string
}

type Result struct {
	Status      int
	URL         string
	Title       string
	Location    string
	ContentType string
	Length      int
	Body        string
	Recursion   int
}

type Options struct {
	Method                 string
	URLs                   []string
	Paths                  []string
	Workers                int
	Timeout                int
	BodyExclude            string
	BodyLengthExcludeTimes int // 过滤响应包长度数据次数出现了多少次，就过滤相同的
	StatusCodeExclude      []int
	Redirect               bool
	Interval               int
	CustomHeader           string
	Recursion              int
	Backupscan             bool
}

func NewDirsearchEngine(ctx, ctrlCtx context.Context, o Options) *Dirsearch {
	headers := clients.Str2HeadersMap(o.CustomHeader)

	return &Dirsearch{
		ctx:           ctx,
		options:       o,
		bodyLengthMap: make(map[int]int),
		client:        clients.NewRestyClient(nil, o.Redirect),
		headers:       headers,
	}
}

func (s *Dirsearch) Runner(ctrlCtx context.Context) {
	runtime.EventsEmit(s.ctx, "dirsearchCounts", len(s.options.URLs)*len(s.options.Paths))
	// 初始化请求信息
	if s.options.Timeout == 0 {
		s.options.Timeout = 8
	}
	var id int32
	single := make(chan struct{})
	retChan := make(chan Result)
	var wg sync.WaitGroup
	go func() {
		for pr := range retChan {
			pr.Recursion = s.options.Recursion
			runtime.EventsEmit(s.ctx, "dirsearchLoading", pr)
		}
		close(single)
		runtime.EventsEmit(s.ctx, "dirsearchComplete", "done")
	}()

	dirScan := func(url string) {
		r := s.Scan(s.ctx, url)
		atomic.AddInt32(&id, 1)
		runtime.EventsEmit(s.ctx, "dirsearchProgressID", id)
		retChan <- r
	}
	threadPool, _ := ants.NewPoolWithFunc(s.options.Workers, func(p interface{}) {
		path := p.(string)
		dirScan(path)
		wg.Done()
	})
	defer threadPool.Release()
	for _, url := range s.options.URLs {
		url = httputil.PrettyURL(url)
		for _, path := range s.options.Paths {
			if ctrlCtx.Err() != nil {
				return
			}
			path = httputil.PrettyPath(path)
			wg.Add(1)
			threadPool.Invoke(url + path)
			if s.options.Interval != 0 {
				time.Sleep(time.Second * time.Duration(s.options.Interval))
			}
		}
	}
	wg.Wait()
	close(retChan)
	<-single
}

// status 1 表示被排除显示在外，不计入前端ERROR请求中
func (s *Dirsearch) Scan(ctx context.Context, url string) Result {
	var result Result
	result.URL = url
	resp, err := clients.DoRequest(s.options.Method, url, s.headers, nil, s.options.Timeout, s.client)
	if err != nil {
		gologger.IntervalError(ctx, err)
		return result
	}
	if s.options.BodyExclude != "" && bytes.Contains(resp.Body(), []byte(s.options.BodyExclude)) {
		result.Status = 1
		return result
	} else {
		result.Status = resp.StatusCode()
	}
	if arrayutil.ArrayContains(result.Status, s.options.StatusCodeExclude) {
		result.Status = 1
		return result
	}
	result.Length = len(resp.Body())
	// 记录同一状态码下长度出现次数，当次数超过o.BodyLengthExcludeTimes时，将状态码设置为1，过滤显示
	s.mutex.Lock()
	s.bodyLengthMap[result.Length]++
	if s.bodyLengthMap[result.Length] > s.options.BodyLengthExcludeTimes {
		result.Status = 1
	}
	s.mutex.Unlock()
	result.Body = httputil.LimitResponse(string(resp.Body()), maxReponseLength, "响应包过大, 请自行打开连接查看")
	result.Location = resp.Header().Get("Location")
	result.Title = clients.GetTitle(resp.Body())
	return result
}

type ScanTask struct {
	URL  string
	Path string
}

func (s *Dirsearch) BackupRunner(ctrlCtx context.Context) {
	var tasks []ScanTask
	for _, url := range s.options.URLs {
		prettyURL := httputil.PrettyURL(url)
		paths := generateBackupPaths(prettyURL)
		for _, path := range paths {
			prettyPath := httputil.PrettyPath(path)
			tasks = append(tasks, ScanTask{
				URL:  prettyURL,
				Path: prettyPath,
			})
		}
	}

	runtime.EventsEmit(s.ctx, "dirsearchCounts", len(tasks))

	if s.options.Timeout == 0 {
		s.options.Timeout = 8
	}

	var id int32
	var wg sync.WaitGroup
	retChan := make(chan Result)
	single := make(chan struct{})

	go func() {
		for pr := range retChan {
			if strings.Contains(pr.ContentType, "html") ||
				strings.Contains(pr.ContentType, "image") ||
				strings.Contains(pr.ContentType, "xml") ||
				strings.Contains(pr.ContentType, "text") ||
				strings.Contains(pr.ContentType, "json") ||
				strings.Contains(pr.ContentType, "javascript") {
				continue
			}
			pr.Recursion = s.options.Recursion
			runtime.EventsEmit(s.ctx, "dirsearchLoading", pr)
		}
		close(single)
		runtime.EventsEmit(s.ctx, "dirsearchComplete", "done")
	}()

	threadPool, _ := ants.NewPoolWithFunc(s.options.Workers, func(p interface{}) {
		task := p.(ScanTask)
		r := s.Scan(s.ctx, task.URL+task.Path)
		atomic.AddInt32(&id, 1)
		runtime.EventsEmit(s.ctx, "dirsearchProgressID", id)
		retChan <- r
		wg.Done()
	})
	defer threadPool.Release()

	for _, task := range tasks {
		if ctrlCtx.Err() != nil {
			return
		}
		wg.Add(1)
		threadPool.Invoke(task)
		if s.options.Interval != 0 {
			time.Sleep(time.Second * time.Duration(s.options.Interval))
		}
	}
	wg.Wait()
	close(retChan)
	<-single
}
func generateBackupPaths(baseURL string) []string {
	suffixFormat := []string{
		".zip", ".rar", ".tar.gz", ".tgz", ".tar.bz2", ".tar", ".jar", ".war",
		".7z", ".bak", ".sql", ".gz", ".sql.gz", ".tar.tgz", ".backup",
	}
	infoDict := []string{
		"1", "127.0.0.1", "admin", "backup", "backups", "db", "sql", "root", "www", "website",
		"2020", "2021", "2022", "2023", "2024", "2025", "old", "new", "test", "data", "dump",
		"web", "wordpress", "wp", "localhost", "users", "store", "auth",
	}

	// 提取主域名用于组合（如：www.example.com -> example）
	host := strings.TrimPrefix(strings.TrimPrefix(baseURL, "http://"), "https://")
	if strings.Contains(host, "/") {
		host = strings.Split(host, "/")[0]
	}
	domainParts := strings.Split(host, ".")

	domainVariants := map[string]struct{}{
		host:                               {},
		strings.ReplaceAll(host, ".", ""):  {},
		strings.ReplaceAll(host, ".", "_"): {},
	}

	if len(domainParts) > 1 {
		domainVariants[strings.Join(domainParts[1:], "")] = struct{}{}
		domainVariants[strings.Join(domainParts[1:], "_")] = struct{}{}
		domainVariants[domainParts[0]] = struct{}{}
	}

	for d := range domainVariants {
		infoDict = append(infoDict, d)
	}

	paths := []string{}
	for _, info := range infoDict {
		for _, suffix := range suffixFormat {
			paths = append(paths, info+suffix)
		}
	}
	return paths
}
