package runner

import (
	"slack-wails/core/webscan/poc"
	"slack-wails/lib/report"

	"sync"
	"time"
)

var CheckerPool = sync.Pool{
	New: func() any {
		return &Checker{
			Options: &Options{},
			// OriginalRequest: &http.Request{},
			VariableMap: make(map[string]any),
			Result:      &report.Result{},
			CustomLib:   NewCustomLib(),
		}
	},
}

func (e *Engine) AcquireChecker() *Checker {
	c := CheckerPool.Get().(*Checker)
	c.Options = e.options
	//c.Result.Output = e.options.Output
	return c
}

func (e *Engine) ReleaseChecker(c *Checker) {
	c.VariableMap = make(map[string]any)
	c.Result = &report.Result{}
	c.CustomLib = NewCustomLib()
	CheckerPool.Put(c)
}

type Engine struct {
	options *Options
	ticker  *time.Ticker
}

// func (runner *Runner) Execute(t string, pocpathList []string) {
// 	options := runner.options
// 	pocSlice := options.CreatePocList(pocpathList) // 执行指纹扫描逻辑就是让pocSilce只返回指纹对应的poc路径即可
// 	reversePocs, otherPocs := options.ReversePoCs(pocSlice)
// 	rwg := sync.WaitGroup{}

// 	rwg.Add(1)
// 	go func() {
// 		defer rwg.Done()

// 		runner.engine.ticker = time.NewTicker(time.Second / time.Duration(options.ReverseRateLimit))
// 		var wg sync.WaitGroup
// 		p, _ := ants.NewPoolWithFunc(options.ReverseConcurrency, func(p any) {

// 			defer wg.Done()
// 			<-runner.engine.ticker.C

// 			tap := p.(*TransData)
// 			runner.exec(tap)

// 		})
// 		defer p.Release()

// 		for _, poc := range reversePocs {
// 			wg.Add(1)
// 			p.Invoke(&TransData{Target: t, Poc: poc}) // 这里执行反弹poc
// 		}
// 		wg.Wait()
// 	}()

// 	rwg.Add(1)
// 	go func() {
// 		defer rwg.Done()

// 		runner.engine.ticker = time.NewTicker(time.Second / time.Duration(options.RateLimit))
// 		var wg sync.WaitGroup

// 		p, _ := ants.NewPoolWithFunc(options.Concurrency, func(p any) {

// 			defer wg.Done()
// 			<-runner.engine.ticker.C

// 			tap := p.(*TransData)
// 			runner.exec(tap)

// 		})
// 		defer p.Release()
// 		for _, poc := range otherPocs {
// 			wg.Add(1)
// 			p.Invoke(&TransData{Target: t, Poc: poc}) // 这里执行正常poc
// 		}

// 		wg.Wait()
// 	}()

//		rwg.Wait()
//	}
func (runner *Runner) Execute(t string, pocpathList []string) {
	options := runner.options
	pocSlice := options.CreatePocList(pocpathList)
	reversePocs, otherPocs := options.ReversePoCs(pocSlice)
	rwg := &sync.WaitGroup{} // Use pointer to WaitGroup

	// Handling reversePocs
	rwg.Add(1)
	go func() {
		defer rwg.Done()

		runner.engine.ticker = time.NewTicker(time.Second / time.Duration(options.ReverseRateLimit))
		var wg sync.WaitGroup

		for _, p := range reversePocs {
			wg.Add(1)
			go func(poc poc.Poc) {
				defer wg.Done()
				<-runner.engine.ticker.C
				runner.exec(&TransData{Target: t, Poc: poc})
			}(p)
		}
		wg.Wait()
	}()

	// Handling otherPocs
	rwg.Add(1)
	go func() {
		defer rwg.Done()

		runner.engine.ticker = time.NewTicker(time.Second / time.Duration(options.RateLimit))
		var wg sync.WaitGroup

		for _, p := range otherPocs {
			wg.Add(1)
			go func(poc poc.Poc) {
				defer wg.Done()
				<-runner.engine.ticker.C
				runner.exec(&TransData{Target: t, Poc: poc})
			}(p)
		}
		wg.Wait()
	}()

	rwg.Wait()
}

func (runner *Runner) exec(tap *TransData) {
	if len(tap.Target) > 0 && len(tap.Poc.Id) > 0 {
		runner.executeExpression(tap.Target, &tap.Poc)
	}
}

func (runner *Runner) executeExpression(target string, poc *poc.Poc) {
	c := runner.engine.AcquireChecker()
	defer runner.engine.ReleaseChecker(c)

	defer func() {
		if r := recover(); r != nil {
			c.Result.IsVul = false
			runner.OnResult(c.Result)
		}
	}()

	c.Check(target, poc)
	runner.OnResult(c.Result)
}

type TransData struct {
	Target string
	Poc    poc.Poc
}
