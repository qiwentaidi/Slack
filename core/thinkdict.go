package core

import (
	"slack-wails/lib/util"
	"strconv"
	"strings"
	"time"

	"github.com/mozillazg/go-pinyin"
)

var weakPasswordList = []string{"123456", "000000", "aa123456", "Aa123456", "Abc123!", "Abc123!@#", "abc123!", "abc1234!", "@bcd1234", "abc123!@#", "Abc123!@#", "666666", "888888", "88888888", "#EDC4rfv", "abcABC123", "1qaz!@#$", "admin@123", "Admin@123", "admin@1234", "Admin@1234", "QAZwsx123", "Pa$$w0rd", "P@ssw0rd", "P@$$word", "P@$$word123", "Abcd1234", "!QAZ2wsx", "!QAZ3edc", "2wsx#EDC", "1!qaz2@wsx", "1q2w3e4r", "1234abcd", "1234qwer", "1qaz!QAZ", "1qaz2wsx", "1qaz@WSX", "1qaz@WSX#EDC", "!q2w3e4r", "1234qwer", "1234QWER", "QWER!@#$", "Passwd@123", "Passwd12", "Passwd@123456", "P@ssw0rd", "1qaz@WSX#EDC", "p@ssw0rd", "qazasd123", "qazwsxedc123", "qweasdzxcqaz123", "asdf1234", "123456Aa", "Aa123456", "123456!Aa", "111111Aa", "111111"}

func GenerateDict(userNameCN, userNameEN, companyName, companyDomain, birthday, jobNumber, connectWord string, weakList []string) (dicts []string) {
	names := []string{Chinese2PinyinQuanPin(userNameCN), Chinese2PinyinFirstLetter(userNameCN), Chinese2PinyinHalfQuanPin(userNameCN), userNameEN}
	companyNames := []string{Chinese2PinyinQuanPin(companyName), Chinese2PinyinFirstLetter(companyName)}
	year := time.Now().Year()
	for i := 0; i <= 6; i++ {
		weakList = append(weakList, strconv.Itoa(year-i))
	}
	if birthday != "" {
		for _, name := range names {
			dicts = append(dicts, name+birthday)
			dicts = append(dicts, name+"@"+birthday)
			dicts = append(dicts, name+birthday[2:])
			dicts = append(dicts, name+"@"+birthday[2:])
		}
	}
	for _, weakpass := range weakList {
		for _, name := range names {
			dicts = append(dicts, name+weakpass)
			dicts = append(dicts, name+"@"+weakpass)
		}
		// baidu.com [0]  www.baidu.com [1] 1111.www.baidu.com [2]
		if strings.Contains(companyDomain, ".") {
			d := strings.Split(companyDomain, ".")
			cd := d[len(d)-2]
			dicts = append(dicts, cd+weakpass)
			dicts = append(dicts, cd+"@"+weakpass)
		} else {
			dicts = append(dicts, companyDomain+weakpass)
			dicts = append(dicts, companyDomain+"@"+weakpass)
		}
	}
	for _, name := range names {
		dicts = append(dicts, name+companyDomain)
		dicts = append(dicts, name+"@"+companyDomain)
		dicts = append(dicts, companyDomain+"@"+name)
		dicts = append(dicts, companyDomain+name)
		for _, company := range companyNames {
			dicts = append(dicts, name+company)
			dicts = append(dicts, name+"@"+company)
			dicts = append(dicts, company+"@"+name)
			dicts = append(dicts, company+name)
		}
	}
	for _, company := range companyNames {
		for _, weak := range weakList {
			dicts = append(dicts, company+weak)
			dicts = append(dicts, company+"@"+weak)
		}
		dicts = append(dicts, company+jobNumber)
		dicts = append(dicts, company+"@"+jobNumber)
	}
	dicts = append(dicts, weakPasswordList...)
	var connectWords []string
	if connectWord != "" {
		connectWords = append(connectWords, strings.Split(connectWord, ",")...)
	}
	for _, dict := range dicts {
		if strings.Contains(dict, "@") {
			for _, word := range connectWords {
				dicts = append(dicts, strings.ReplaceAll(dict, "@", word))
			}
		}
	}
	return util.RemoveDuplicates(dicts)
}

func Chinese2PinyinFirstLetter(str string) (fl string) {
	py := pinyin.NewArgs()
	py.Style = pinyin.FirstLetter              // 设置为获取拼音首字母（只包含首字母）
	for _, v := range pinyin.Pinyin(str, py) { // 将字符串转换为拼音首字母
		fl += v[0]
	}
	if fl == "" { // 说明不是中文
		return str
	}
	return fl
}

func Chinese2PinyinHalfQuanPin(str string) (fl string) {
	a := pinyin.NewArgs()
	dict := pinyin.Pinyin(str, a)
	for _, qpzm := range dict {
		fl += strings.Join(qpzm, "")
	}
	return fl
}

func Chinese2PinyinQuanPin(str string) (qp string) {
	// 默认
	a := pinyin.NewArgs()
	for _, qpzm := range pinyin.Pinyin(str, a) {
		for _, v := range qpzm {
			qp += v
		}
	}
	if qp == "" { // 说明不是中文
		return str
	}
	return qp
}

func FirstUpper(str string) string {
	bs := []byte(str)
	if len(bs) == 0 {
		return ""
	}
	bs[0] = byte(bs[0] - 32)
	return string(bs)
}
