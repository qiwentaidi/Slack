<h4 align="center">一款安服集成化工具平台，希望能让你少开几个应用测试</h4>

<p align="center">
<img src="https://img.shields.io/github/go-mod/go-version/qiwentaidi/Slack?filename=go.mod">
<img src="https://img.shields.io/badge/wails-v2.9.1-blue">
<a href="https://github.com/qiwentaidi/Slack/releases/"><img src="https://img.shields.io/github/v/release/qiwentaidi/Slack"></a>
<a href="https://github.com/qiwentaidi/Slack/releases/"><img src="https://img.shields.io/github/downloads/qiwentaidi/Slack/total"></a>
</p>
<p align="center">
<a href="https://github.com/qiwentaidi/Slack/wiki/%E5%B8%B8%E8%A7%81%E9%97%AE%E9%A2%98">常见问题</a>
<a href="https://github.com/qiwentaidi/Slack/wiki/%E6%9B%B4%E6%96%B0%E6%97%A5%E5%BF%97">更新日志</a>
</p>



# 运行代码

## 安装环境

Linux/Debian

```sh
# need install Go 1.20+
sudo apt install golang-go

# need install nodejs 15+
sudo apt install nodejs npm

# install gcc && libpcap
sudo apt install build-essential libgtk-3-dev libpcap-dev 

# if apt can't be found 4.0-dev then install 4.1-dev
1、sudo apt install libwebkit2gtk-4.0-dev
2、sudo apt install libwebkit2gtk-4.1-dev

# install wails to run app
go install github.com/wailsapp/wails/v2/cmd/wails@latest

git clone https://github.com/qiwentaidi/Slack.git
```

Windows

```sh
# need install Go 1.20+ && nodejs 15+
# download gcc(https://github.com/niXman/mingw-builds-binaries/releases) and configure environment variables

# install wails to run app
go install github.com/wailsapp/wails/v2/cmd/wails@latest

git clone https://github.com/qiwentaidi/Slack.git
```

## 编译/调试

``````sh
cd Slack
# used for debugging while writing code
wails dev 

# build executable file on /build/bin
wails build

# open web console in executable file
wails build -debug -devtools
``````

`Mac build dmg`

``````sh
# need build Slack.app first

brew install create-dmg

create-dmg --volname "Slack" --window-pos 200 120 --window-size 800 400 --icon-size 100  --icon "Slack.app" 200 190 --app-drop-link 600 185 --hide-extension "Slack.app" --volicon build/bin/Slack.app/Contents/Resources/iconfile.icns  "Slack.dmg" build/bin/Slack.app
``````

# 模块介绍

## 首页

![image-20240801105341219](assets/image-20240801105341219.png)

## 渗透测试

### 网站扫描

![image-20240510214011206](assets/image-20240510214011206.png)

![image-20240510214144484](assets/image-20240510214144484.png)

### 主机扫描

![image-20240628171638340](assets/image-20240628171638340.png)

### 暴破与未授权检测

![image-20240628173000629](assets/image-20240628173000629.png)

### 目录扫描

![image-20240410135353422](assets/image-20240410135353422.png)

### JSFinder

![image-20240425150345629](assets/image-20240425150345629.png)

## 资产收集

### 公司名称查资产

![image-20240718104209501](assets/image-20240718104209501.png)

![image-20240412150938371](assets/image-20240412150938371.png)

### 子域名暴破

![image-20240222154726647](assets/image-20240222154726647.png)

### 域名信息查询

![image-20240629003953025](assets/image-20240629003953025.png)

## 空间引擎

### FOFA

![image-20240718103758587](assets/image-20240718103758587.png)

### 鹰图

![image-20240620234045723](assets/image-20240620234045723.png)

### Quake

![image-20240628172913333](assets/image-20240628172913333.png)

聚合搜索

## 小工具

### 加解密模块

![image-20240628172153383](assets/image-20240628172153383.png)

### 数据处理

![image-20240628172250510](assets/image-20240628172250510.png)

### 杀软/补丁识别

![image-20240628172737827](assets/image-20240628172737827.png)

### 反弹Shell生成

![image-20240628172826365](assets/image-20240628172826365.png)

### 备忘录

![image-20240628172702349](assets/image-20240628172702349.png)

### WeChat/DingDing

![image-20240512003433659](assets/image-20240512003433659.png)

# 联系方式

如果有问题或者好的提议可以Issue提问或者加我联系方式（请备注来意 进群或者问题交流）

![image-20231006124944803](assets/image-20231006124944803.png)

# 免责声明

本工具仅面向**合法授权**的企业安全建设行为，如您需要测试本工具的可用性，请自行搭建靶机环境。

为避免被恶意使用，本项目所有收录的poc均为漏洞的理论判断，不存在漏洞利用过程，不会对目标发起真实攻击和漏洞利用。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。**请勿对非授权目标进行扫描。**

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

在安装并使用本工具前，请您**务必审慎阅读、充分理解各条款内容**，限制、免责条款或者其他涉及您重大权益的条款可能会以加粗、加下划线等形式提示您重点注意。 除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。
