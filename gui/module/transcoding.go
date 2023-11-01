package module

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html"
	"io"
	"net/url"
	"slack/gui/custom"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const cmd5 = "https://cmd5.com/"

func TranscodingUI() *fyne.Container {
	encode := widget.NewButtonWithIcon("encode", theme.MoveDownIcon(), nil)
	decode := widget.NewButtonWithIcon("decode", theme.MoveUpIcon(), nil)
	tobeEncry := custom.NewMultiLineEntryText("")                       // 待加密框
	tobeDecry := custom.NewMultiLineEntryPlaceHolder("需要解密的文本一律放在下文本框") // 待解密框
	auto := &widget.Check{Text: "自动", Checked: true}
	vbox := container.NewBorder(nil, container.NewBorder(nil, nil, auto, nil, container.NewGridWithColumns(2, encode, decode)), nil, nil, tobeEncry)
	list := container.NewVBox()
	add := widget.NewButtonWithIcon("添加规则(右键删除)", theme.ContentAddIcon(), func() {
		list.Add(custom.NewTappedSelect([]string{"Base64", "Base64 URL", "URLcode", "Unicode", "HEX", "HTML", "MD5", "Ascii", "FanruanOA", "SeeyonOA"}, list))
	})
	tobeEncry.OnChanged = func(s string) {
		if auto.Checked {
			var en []string
			for _, m := range list.Objects {
				en = append(en, m.(*custom.TappedSelect).Selected)
			}
			recursiveEncrypt(en, tobeEncry.Text, &tobeDecry.Text)
			tobeDecry.Refresh()
		}
	}
	tobeDecry.OnChanged = func(s string) {
		var de []string
		for _, m := range list.Objects {
			de = append(de, m.(*custom.TappedSelect).Selected)
		}
		recursiveDecrypt(de, tobeDecry.Text, &tobeEncry.Text)
		tobeEncry.Refresh()
	}
	encode.OnTapped = func() {
		var en []string
		for _, m := range list.Objects {
			en = append(en, m.(*custom.TappedSelect).Selected)
		}
		recursiveEncrypt(en, tobeEncry.Text, &tobeDecry.Text)
		tobeDecry.Refresh()
	}
	decode.OnTapped = func() {
		var de []string
		for _, m := range list.Objects {
			de = append(de, m.(*custom.TappedSelect).Selected)
		}
		recursiveDecrypt(de, tobeDecry.Text, &tobeEncry.Text)
		tobeEncry.Refresh()
	}
	u, _ := url.Parse(cmd5)
	link := widget.NewHyperlinkWithStyle("MD5解密", u, fyne.TextAlignCenter, fyne.TextStyle{})
	hbox := container.NewHSplit(container.NewBorder(add, link, nil, nil, container.NewVScroll(list)), container.NewVSplit(vbox, tobeDecry))
	hbox.Offset = 0.2
	return container.NewBorder(nil, nil, nil, nil, hbox)
}

func recursiveEncrypt(encrypt []string, input string, output *string) {
	if len(encrypt) == 0 {
		*output = input
		return
	}
	var encrypted string
	switch encrypt[0] {
	case "Base64":
		encrypted = base64.StdEncoding.EncodeToString([]byte(input))
	case "Base64 URL":
		encrypted = base64.URLEncoding.EncodeToString([]byte(input))
	case "URLcode":
		encrypted = url.QueryEscape(input)
	case "Unicode":
		temp := strconv.QuoteToASCII(input)
		encrypted = temp[1 : len(temp)-1]
	case "HEX":
		encrypted = hex.EncodeToString([]byte(input))
	case "HTML":
		encrypted = html.EscapeString(input)
	case "MD5":
		if input != "" {
			hash := md5.New()
			io.WriteString(hash, input)
			encrypted = hex.EncodeToString(hash.Sum(nil))
		}
	case "Ascii":
		var temp string
		for _, v := range input {
			temp += fmt.Sprintf("%v ", v)
		}
		encrypted = temp
	default:
		encrypted = input
	}
	recursiveEncrypt(encrypt[1:], encrypted, output)
}

func recursiveDecrypt(decrypt []string, input string, output *string) {
	if len(decrypt) == 0 {
		*output = input
		return
	}
	var decrypted string
	switch decrypt[0] {
	case "Base64":
		data, _ := base64.StdEncoding.DecodeString(input)
		decrypted = string(data)
	case "Base64 URL":
		data, _ := base64.URLEncoding.DecodeString(input)
		decrypted = string(data)
	case "URLcode":
		data, _ := url.QueryUnescape(input)
		decrypted = data
	case "Unicode":
		data, _ := strconv.Unquote("\"" + input + "\"")
		decrypted = data
	case "HEX":
		data, _ := hex.DecodeString(input)
		decrypted = string(data)
	case "HTML":
		html.UnescapeString(input)
	case "Ascii":
		temp := ""
		if strings.Contains(input, " ") {
			for _, s := range strings.Split(input, " ") {
				if i, err := strconv.Atoi(s); err == nil {
					temp += fmt.Sprintf("%c", i)
				} else {
					temp += fmt.Sprintf("what?(%v)", s)
				}
			}
		} else {
			if i, err := strconv.Atoi(input); err == nil {
				temp = fmt.Sprintf("%c", i)
			}
		}
		decrypted = temp
	case "FanruanOA":
		PASSWORD_MASK_ARRAY := [8]int{19, 78, 10, 15, 100, 213, 43, 23} // 掩码
		Password := ""
		if len(input) > 3 {
			input = input[3:] // 截断三位后
			for i := 0; i < len(input)/4; i++ {
				c1, _ := strconv.ParseInt(input[i*4:(i+1)*4], 16, 32)
				c2 := int(c1) ^ PASSWORD_MASK_ARRAY[i%8]
				Password = Password + string(rune(c2))
			}
		}
		decrypted = Password
	case "SeeyonOA-Database":
		pass := strings.ReplaceAll(input, "/", "")
		p := strings.Split(pass, ".0")
		if len(p) > 1 {
			iv, _ := strconv.Atoi(p[0])
			password, _ := base64.StdEncoding.DecodeString(p[1])
			var builder strings.Builder
			for _, char := range password {
				builder.WriteString(string(char - byte(iv)))
			}
			decrypted = builder.String()
		}
	default:
		decrypted = input
	}
	recursiveDecrypt(decrypt[1:], decrypted, output)
}
