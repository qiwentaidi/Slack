<h4 align="center">一款安服集成化工具平台，希望能让你少开几个应用测试</h4>

<p align="center">
<img src="https://img.shields.io/github/go-mod/go-version/qiwentaidi/Slack?filename=go.mod">
<img src="https://img.shields.io/badge/wails-v2.9.1-blue">
<a href="https://github.com/qiwentaidi/Slack/releases/"><img src="https://img.shields.io/github/v/release/qiwentaidi/Slack"></a>
<a href="https://github.com/qiwentaidi/Slack/releases/"><img src="https://img.shields.io/github/downloads/qiwentaidi/Slack/total"></a>
</p>
<p align="center">
<a href="https://github.com/qiwentaidi/Slack/wiki/%E5%B8%B8%E8%A7%81%E9%97%AE%E9%A2%98">常见问题</a>
<a href="https://github.com/qiwentaidi/Slack/wiki/%E7%BD%91%E7%AB%99%E6%89%AB%E6%8F%8F%E8%A7%84%E5%88%99%E4%BA%8C%E6%AC%A1%E6%8B%93%E5%B1%95">规则拓展</a>
<a href="https://github.com/qiwentaidi/Slack/wiki/%E4%BA%8C%E6%AC%A1%E5%BC%80%E5%8F%91">二次开发</a>
<a href="https://github.com/qiwentaidi/Slack/wiki/%E6%9B%B4%E6%96%B0%E6%97%A5%E5%BF%97">更新日志</a>
</p>


# 支持的平台

- Windows 10/11 AMD64/ARM64
- MacOS 10.13+ AMD64
- MacOS 11.0+ ARM64
- Linux AMD64/ARM64

# 运行代码

## 安装依赖

- Go 1.20+
- Node.js 16+
- Wails 2.9+  `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Linux

- build-essential
- libgtk-3-dev
- libpcap-dev
- libwebkit2gtk-4.0-dev or libwebkit2gtk-4.1-dev

> [!NOTE]
>
> 新版`Linux`安装`libwebkit2gtk-4.1-dev`编译应用时需要增加 `-tags webkit2_41`

### Windows

- [gcc](https://github.com/niXman/mingw-builds-binaries/releases)

### Mac

- xcode-select（已默认安装）

## 编译/调试

``````sh
git clone https://github.com/qiwentaidi/Slack.git && cd Slack

wails dev # 调试模式运行

wails build # 编译应用 文件存在于build/bin路径下

wails build -debug -devtools # 编译可开启调试模式应用
``````

`Mac build dmg`

``````sh
# 需要先编译完Slack.app
brew install create-dmg

create-dmg --volname "Slack" --window-pos 200 120 --window-size 800 400 --icon-size 100  --icon "Slack.app" 200 190 --app-drop-link 600 185 --hide-extension "Slack.app" --volicon build/bin/Slack.app/Contents/Resources/iconfile.icns  "Slack.dmg" build/bin/Slack.app
``````

# 首页

![image-20240907162449101](assets/image-20240907162449101.png)

# 特色功能介绍

## 端口扫描

可以联动网站扫描以及协议爆破

![image-20240928220948717](assets/image-20240928220948717.png)

## 网站扫描

目前内置8600+指纹，2700+POC，易用可扩展

![image-20241006193512657](assets/image-20241006193512657.png)

![image-20240929152018714](assets/image-20240929152018714.png)

## 公司信息查询

可以通过输入公司名称一步完成IP、域名的收集

![image-20240907180729539](assets/image-20240907180729539.png)

![image-20240907180805041](assets/image-20240907180805041.png)

![image-20240907202829701](assets/image-20240907202829701.png)

## 空间搜索

`FOFA`、`Hunter`、`Quake`查询功能，保留搜索提示、语法收藏以及数据可视性的同时，增加特色功能区，减少数据导出操作。

![image-20240907172923247](assets/image-20240907172923247.png)

![image-20240907172244822](assets/image-20240907172244822.png)

![image-20240907172356305](assets/image-20240907172356305.png)

![image-20240907172750605](assets/image-20240907172750605.png)

![image-20240907173203880](assets/image-20240907173203880.png)

## 数据库自动取样

对敏感列名以及数据内容进行匹配，内容支持身份证、手机号、AKSK信息匹配。

![image-20241006194025842](assets/image-20241006194025842.png)

![image-20241006193952875](assets/image-20241006193952875.png)

## 加解密模块

可通过下载`CyberChef`集成环境实现本地调用

![image-20240907172017116](assets/image-20240907172017116.png)

## 数据处理

针对日常工作中一些常见的数据进行处理，例如提取Fscan结果、提取IP、数据去重等

![image-20240907203934509](assets/image-20240907203934509.png)

## 应用启动器

用于管理繁琐的脚本，可以自定义启动命令，支持`cmd`打开文件所在命令行、`java`以`java -jar`命令启动`java GUI`应用、`App`打开`exe GUI`

![image-20240907234517306](assets/image-20240907234517306.png)

可视化日志（右上角终端符号打开）

![image-20240907181024885](assets/image-20240907181024885.png)

# 联系方式

如果有问题或者好的提议可以Issue提问或者加我联系方式（请备注来意 进群或者问题交流）

![image-20231006124944803](assets/image-20231006124944803.png)

# 免责声明

本工具仅面向**合法授权**的企业安全建设行为，如您需要测试本工具的可用性，请自行搭建靶机环境。

为避免被恶意使用，本项目所有收录的poc均为漏洞的理论判断，不存在漏洞利用过程，不会对目标发起真实攻击和漏洞利用。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。**请勿对非授权目标进行扫描。**

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

在安装并使用本工具前，请您**务必审慎阅读、充分理解各条款内容**，限制、免责条款或者其他涉及您重大权益的条款可能会以加粗、加下划线等形式提示您重点注意。 除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。
