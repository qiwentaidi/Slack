package report

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slack/lib/result"
	"sync"
	"time"
)

type Report struct {
	sync.RWMutex
	Result     *result.Result
	of         *os.File
	ReportFile string
	Template   TemplateStyle
}

type TemplateStyle int

const (
	DefaultTemplate TemplateStyle = iota
)

const OutputDirectory = "./reports"

// fileName: the name of the report file
// template: the name of the template
func NewReport(template TemplateStyle) (*Report, error) {
	r := &Report{
		Result:   &result.Result{},
		Template: template,
	}

	if err := r.check(); err != nil {
		return nil, err
	}

	return r, nil
}

func (report *Report) check() error {
	if _, err := os.Stat(OutputDirectory); err != nil {
		os.Mkdir(OutputDirectory, 0700)
	}
	fileName := filepath.Join(OutputDirectory, time.Now().Format("20060102-150405")+".html")
	report.ReportFile = fileName
	if _, err := os.Stat(fileName); err != nil {
		file, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("unable to create output file: %v", err)
		}
		file.Close()
		time.Sleep(100 * time.Millisecond)
		return os.Remove(fileName)
	}
	return nil
}

func (report *Report) SetResult(result *result.Result) {
	report.Lock()
	defer report.Unlock()

	report.Result = result
}

func (report *Report) Append(number string) error {
	return report.Write(report.html(number))
}

func (report *Report) Write(data string) error {
	if len(data) == 0 {
		return nil
	}

	report.Lock()

	var f *os.File
	if report.of == nil {
		f, err := os.OpenFile(report.ReportFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			report.Unlock()
			return err
		}

		report.of = f

		header := report.header()
		wbuf := bufio.NewWriterSize(report.of, len(header))
		wbuf.WriteString(header)
		wbuf.Flush()
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(f *os.File) {
		defer report.Unlock()
		defer f.Close()
		defer wg.Done()

		wbuf := bufio.NewWriterSize(report.of, len(data))
		wbuf.WriteString(data)
		wbuf.Flush()
	}(f)
	wg.Wait()

	return nil
}

func (report *Report) header() string {
	switch report.Template {
	case DefaultTemplate:
		return defaultHeader()
	}
	return ""
}

func (report *Report) html(number string) string {
	switch report.Template {
	case DefaultTemplate:
		return report.defaultHmtl(number)
	}
	return ""
}
