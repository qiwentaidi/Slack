package webscan

import (
	"container/list"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slack-wails/lib/gologger"
	"slack-wails/lib/util"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// 总和需要加载的模板文件夹
	TemplateFolders []string
	// 指纹规则文件
	// FingerprintRuleFile string
	// 主动探测的规则文件
	// ActiveRuleFile string
	// 默认模板文件夹
	// DefaultTemplateFolder string
}

func NewConfig(template []string) *Config {
	return &Config{
		TemplateFolders: template,
	}
}

type FingerPEntity struct {
	ProductName      string
	AllString        string
	Rule             []RuleData
	IsExposureDetect bool
}

type RuleData struct {
	Start int
	End   int
	Op    int16  // 0= 1!= 2== 3>= 4<= 5~=
	Key   string // body="123"中的body
	Value string // body="123"中的123
	All   string // body="123"
}

type ActiveFingerPEntity struct {
	Path []string
	Fpe  []FingerPEntity
}

var FingerprintDB []FingerPEntity
var ActiveFingerprintDB []ActiveFingerPEntity

func (config *Config) InitFingprintDB(ctx context.Context, fingerprintFile string) error {
	data, err := os.ReadFile(fingerprintFile)
	if err != nil {
		return err
	}

	fps := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &fps); err != nil {
		return err
	}

	m := make(map[string][]string)
	for productName, rulesInterface := range fps {
		rules, ok := rulesInterface.([]interface{})
		if !ok {
			fmt.Printf("Invalid fingerprint format for product [%s], rules [%v]\n", productName, rulesInterface)
			continue
		}

		for _, ruleInterface := range rules {
			rule, ok := ruleInterface.(string)
			if !ok {
				fmt.Printf("Invalid rule format for product [%s], rule [%v]\n", productName, ruleInterface)
				continue
			}

			if !util.ArrayContains(rule, m[productName]) {
				m[productName] = append(m[productName], rule)
			}
		}
	}

	for productName, ruleList := range m {
		for _, rule := range ruleList {
			FingerprintDB = append(FingerprintDB, FingerPEntity{
				ProductName: productName,
				Rule:        ParseRule(rule),
				AllString:   rule,
			})
		}
	}

	return nil
}
func (config *Config) InitActiveScanPath(activefingerFile string) error {
	data, err := os.ReadFile(activefingerFile)
	if err != nil {
		return err
	}
	sensitive := make(map[string][]string)
	err = yaml.Unmarshal(data, &sensitive)
	if err != nil {
		return err
	}
	for name, paths := range sensitive {
		var fpes []FingerPEntity
		for _, fpe := range FingerprintDB {
			if fpe.ProductName == name {
				fpes = append(fpes, fpe)
			}
		}
		if len(fpes) != 0 {
			ActiveFingerprintDB = append(ActiveFingerprintDB, ActiveFingerPEntity{
				Path: paths,
				Fpe:  fpes,
			})
		}
	}
	return nil
}

func ParseRule(rule string) []RuleData {
	var result []RuleData
	empty := RuleData{}

	for {
		data := getRuleData(rule)
		if data == empty {
			break
		}
		result = append(result, data)
		rule = rule[:data.Start] + "T" + rule[data.End:]
	}
	return result
}

func getRuleData(rule string) RuleData {
	if !strings.Contains(rule, "=\"") {
		return RuleData{}
	}
	pos := strings.Index(rule, "=\"")
	op := 0
	if rule[pos-1] == 33 {
		op = 1
	} else if rule[pos-1] == 61 {
		op = 2
	} else if rule[pos-1] == 62 {
		op = 3
	} else if rule[pos-1] == 60 {
		op = 4
	} else if rule[pos-1] == 126 {
		op = 5
	}

	start := 0
	ti := 0
	if op > 0 {
		ti = 1
	}
	for i := pos - 1 - ti; i >= 0; i-- {
		if (rule[i] > 122 || rule[i] < 97) && rule[i] != 95 {
			start = i + 1
			break
		}
	}
	key := rule[start : pos-ti]

	end := pos + 2
	for i := pos + 2; i < len(rule)-1; i++ {
		if rule[i] != 92 && rule[i+1] == 34 {
			end = i + 2
			break
		}
	}
	value := rule[pos+2 : end-1]
	all := rule[start:end]

	return RuleData{Start: start, End: end, Op: int16(op), Key: key, Value: value, All: all}
}

// 计算纯bool表达式，支持 ! && & || | ( )
func boolEval(ctx context.Context, expression string) bool {
	// 左右括号相等
	if strings.Count(expression, "(") != strings.Count(expression, ")") {
		gologger.Warning(ctx, fmt.Sprintf("纯布尔表达式 [%s] 左右括号不匹配", expression))
	}
	// 去除空格
	for strings.Contains(expression, " ") {
		expression = strings.ReplaceAll(expression, " ", "")
	}
	// 去除空表达式
	for strings.Contains(expression, "()") {
		expression = strings.ReplaceAll(expression, "()", "")
	}
	for strings.Contains(expression, "&&") {
		expression = strings.ReplaceAll(expression, "&&", "&")
	}
	for strings.Contains(expression, "||") {
		expression = strings.ReplaceAll(expression, "||", "|")
	}
	if !strings.Contains(expression, "T") && !strings.Contains(expression, "F") {
		return false
		// panic("纯布尔表达式错误，没有包含T/F")
	}

	expr := list.New()
	operator_stack := list.New()
	for _, ch := range expression {
		// ch 为 T或者F
		if ch == 84 || ch == 70 {
			expr.PushBack(int(ch))
		} else if advance(int(ch)) > 0 {
			if operator_stack.Len() == 0 {
				operator_stack.PushBack(int(ch))
				continue
			}
			// 两个!抵消
			if ch == 33 && operator_stack.Back().Value.(int) == 33 {
				operator_stack.Remove(operator_stack.Back())
				continue
			}
			for operator_stack.Len() != 0 && operator_stack.Back().Value.(int) != 40 && advance(operator_stack.Back().Value.(int)) >= advance(int(ch)) {
				e := operator_stack.Back()
				expr.PushBack(e.Value.(int))
				operator_stack.Remove(e)
			}
			operator_stack.PushBack(int(ch))

		} else if ch == 40 {
			operator_stack.PushBack(int(ch))
		} else if ch == 41 {
			for operator_stack.Back().Value.(int) != 40 {
				e := operator_stack.Back()
				expr.PushBack(e.Value.(int))
				operator_stack.Remove(e)
			}
			operator_stack.Remove(operator_stack.Back())
		}
	}
	for operator_stack.Len() != 0 {
		e := operator_stack.Back()
		expr.PushBack(e.Value.(int))
		operator_stack.Remove(e)
	}

	tf_stack := list.New()
	for expr.Len() != 0 {
		e := expr.Front()
		ch := e.Value.(int)
		expr.Remove(e)
		if ch == 84 || ch == 70 {
			tf_stack.PushBack(int(ch))
		}
		if ch == 38 { // &
			em := tf_stack.Back()
			a := em.Value.(int)
			tf_stack.Remove(em)
			em = tf_stack.Back()
			b := em.Value.(int)
			tf_stack.Remove(em)
			if a == 84 && b == 84 {
				tf_stack.PushBack(84)
			} else {
				tf_stack.PushBack(70)
			}
		}
		if ch == 124 { // |
			em := tf_stack.Back()
			a := em.Value.(int)
			tf_stack.Remove(em)
			em = tf_stack.Back()
			b := em.Value.(int)
			tf_stack.Remove(em)
			if a == 70 && b == 70 {
				tf_stack.PushBack(70)
			} else {
				tf_stack.PushBack(84)
			}
		}
		if ch == 33 { // !
			em := tf_stack.Back()
			a := em.Value.(int)
			tf_stack.Remove(em)
			if a == 70 {
				tf_stack.PushBack(84)
			} else if a == 84 {
				tf_stack.PushBack(70)
			}
		}
	}
	if tf_stack.Front().Value.(int) == 84 {
		return true
	} else {
		return false
	}

}

// 判断优先级 非运算符返回0
func advance(ch int) int {
	// !
	if ch == 33 {
		return 3
	}
	// &
	if ch == 38 {
		return 2
	}
	// |
	if ch == 124 {
		return 1
	}
	return 0
}

func regexMatch(pattern string, s string) (bool, error) {
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false, err
	}
	return matched, nil
}

// body="123"  op=0  dataSource为http.body dataRule=123
func dataCheckString(op int16, dataSource string, dataRule string) bool {
	dataSource = strings.ToLower(dataSource)

	dataRule = strings.ToLower(dataRule)
	dataRule = strings.ReplaceAll(dataRule, "\\\"", "\"")
	if op == 0 {
		if strings.Contains(dataSource, dataRule) {
			return true
		}
	} else if op == 1 {
		if !strings.Contains(dataSource, dataRule) {
			return true
		}
	} else if op == 2 {
		if dataSource == dataRule {
			return true
		}
	} else if op == 5 {
		rs, err := regexMatch(dataRule, dataSource)
		if err == nil && rs {
			return true
		}
	}
	return false
}

func dataCheckInt(op int16, dataSource int, dataRule int) bool {
	if op == 0 { // 数字相等
		if dataSource == dataRule {
			return true
		}
	} else if op == 1 { // 数字不相等
		if dataSource != dataRule {
			return true
		}
	} else if op == 3 { // 大于等于
		if dataSource >= dataRule {
			return true
		}
	} else if op == 4 {
		if dataSource <= dataRule {
			return true
		}
	}
	return false
}

var WorkFlowDB map[string][]string

func (config *Config) InitAll(ctx context.Context, webfinger, activefinger string) bool {
	FingerprintDB = nil
	if err := config.InitFingprintDB(ctx, webfinger); err != nil {
		gologger.Error(ctx, err)
		return false
	}
	if err := config.InitActiveScanPath(activefinger); err != nil {
		gologger.Error(ctx, err)
		return false
	}
	if err := GetTagsList(config.TemplateFolders); err != nil {
		gologger.Error(ctx, err)
		return false
	}
	return true
}

type TemplateInfo struct {
	Tags string `yaml:"tags"`
}

type Template struct {
	Info TemplateInfo `yaml:"info"`
}

var mutx sync.RWMutex

func GetTagsList(templateFolders []string) error {
	WorkFlowDB = make(map[string][]string)
	for _, folder := range templateFolders {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			continue
		}
		// 遍历所有模板文件
		filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(d.Name(), ".yaml") {
				file, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				var template Template
				err = yaml.Unmarshal(file, &template)
				if err != nil {
					return nil
				}

				if template.Info.Tags != "" {
					tags := strings.Split(template.Info.Tags, ",")
					poc := strings.TrimSuffix(d.Name(), ".yaml")
					mutx.Lock()
					WorkFlowDB[poc] = tags
					mutx.Unlock()
				}
			}
			return nil
		})
	}
	return nil
}
