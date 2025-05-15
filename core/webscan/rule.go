package webscan

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slack-wails/lib/gologger"
	"slack-wails/lib/util"
	"strings"
	"sync"

	"github.com/Knetic/govaluate"
	"gopkg.in/yaml.v2"
)

type Config struct {
	// 总和需要加载的模板文件夹
	TemplateFolders []string
	// 指纹规则文件
	FingerprintRuleFile string
	// 主动探测的规则文件
	ActiveRuleFile string
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
	Op    int16          // 0= 1!= 2== 3>= 4<= 5~=
	Key   string         // body="123"中的body
	Value string         // body="123"中的123
	All   string         // body="123"
	Regex *regexp.Regexp // 缓存编译后的正则表达式
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
	// 增加错误判断，防止切片越界
	if end-1 > len(rule) || pos+2 > len(rule) || end-1 < pos+2 {
		fmt.Printf("Error: rule [%s] pos [%d] end [%d] len [%d]\n", rule, pos, end, len(rule))
		return RuleData{}
	}
	value := rule[pos+2 : end-1]
	all := rule[start:end]

	return RuleData{Start: start, End: end, Op: int16(op), Key: key, Value: value, All: all}
}

// 将 T/F 替换为 true/false，并转换逻辑运算符符号
func normalizeExpression(expr string) string {
	expr = strings.ReplaceAll(expr, " ", "") // 去空格
	expr = strings.ReplaceAll(expr, "T", "true")
	expr = strings.ReplaceAll(expr, "F", "false")
	return expr
}

// 使用 govaluate 实现 boolEval
func boolEval(expression string) (bool, error) {
	// 检查是否有 T 或 F
	if !strings.Contains(expression, "T") && !strings.Contains(expression, "F") {
		return false, errors.New("纯布尔表达式错误，没有包含T/F")
	}

	// 标准化表达式
	exprStr := normalizeExpression(expression)

	// 创建表达式对象
	expr, err := govaluate.NewEvaluableExpression(exprStr)
	if err != nil {
		return false, fmt.Errorf("无法解析表达式 [%s]: %v", exprStr, err)

	}

	// 求值
	result, err := expr.Evaluate(nil) // 不需要变量
	if err != nil {
		return false, fmt.Errorf("执行表达式出错 [%s]: %v", exprStr, err)
	}

	// 类型断言并返回布尔结果
	if booleanResult, ok := result.(bool); ok {
		return booleanResult, nil
	}
	return false, errors.New("结果不是布尔类型")
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
	if dataSource == "" {
		return false
	}
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

func (config *Config) InitAll(ctx context.Context) bool {
	FingerprintDB = nil
	if err := config.InitFingprintDB(ctx, config.FingerprintRuleFile); err != nil {
		gologger.Error(ctx, err)
		return false
	}
	if err := config.InitActiveScanPath(config.ActiveRuleFile); err != nil {
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
