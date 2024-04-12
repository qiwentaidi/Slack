<h4 align="center">由Go Wails实现的GUI工具，功能涵盖网站扫描、端口扫描、企业信息收集、子域名暴破、空间引擎搜索、CDN识别等多功能的工具</h4>

<p align="center">
<img src="https://img.shields.io/github/go-mod/go-version/qiwentaidi/Slack?filename=go.mod">
<img src="https://img.shields.io/badge/wails-v2.8.0-blue">
<a href="https://github.com/qiwentaidi/Slack/releases/"><img src="https://img.shields.io/github/v/release/qiwentaidi/Slack">
<a href="https://github.com/qiwentaidi/Slack/releases/"><img src="https://img.shields.io/github/downloads/qiwentaidi/Slack/total">
</p>






# 使用须知

Fyne版本已停止更新，可以通过fyne分支中获取到该模块代码。

# 运行	可以通过[wails](https://wails.io/zh-Hans/docs/gettingstarted/installation/)官网查看支持运行的客户端以及运行前条件，详细功能介绍及拓展等信息见[Wiki](https://github.com/qiwentaidi/Slack/wiki)

```shell
go install github.com/wailsapp/wails/v2/cmd/wails@latest

wails doctor // 检测是否安装成功

git clone https://github.com/qiwentaidi/Slack.git

cd Slack

wails dev // 调试

wails build // 编译

编译完的应用在 build/bin 目录下
```

# 渗透测试

## 网站扫描

![image-20240204164748870](assets/image-20240204164748870.png)

## 主机扫描

![image-20240410135219996](assets/image-20240410135219996.png)

## 目录扫描

![image-20240410135353422](assets/image-20240410135353422.png)

# 资产收集

## 公司名称查资产

![image-20231124155235223](assets/image-20231124155235223.png)

![image-20240412150938371](assets/image-20240412150938371.png)

## 子域名暴破

![image-20240222154726647](assets/image-20240222154726647.png)

## 域名信息查询

![image-20240410134646623](assets/image-20240410134646623.png)

# 空间引擎

## FOFA

![image-20231204000209352](assets/image-20231204000209352.png)

## 鹰图

![image-20231229143634981](assets/image-20231229143634981.png)

![image-20231229143654453](assets/image-20231229143654453.png)

# 小工具

<img src="assets/image-20240222155132016.png" alt="image-20240222155132016" style="zoom:55%;" />



## 常见问题

> Q：Windows grdp库无法打包，出现如下报错
>
> ```
> Wails CLI v2.7.1
> 
> Executing: go mod tidy
> • Generating bindings: 
> ERROR  
>        package slack-wails
>              imports slack-wails/core/portscan
>              imports github.com/tomatome/grdp/protocol/pdu
>              imports github.com/tomatome/grdp/protocol/t125/gcc
>              imports github.com/tomatome/grdp/plugin: build constraints exclude all Go files in C:\xx\go\pkg\mod\github.com\tomatome\grdp@v0.1.0\plugin
> ```
>
> A：
>
> ```
> 1、go env查看CGO_ENABLED是否为1，若不是则go env -w CGO_ENABLED=1 
> 2、需要安装GCC环境
> ```

# 联系方式

如果有问题或者好的提议可以Issue提问或者加我联系方式（请备注来意 进群或者问题交流）

![image-20231006124944803](assets/image-20231006124944803.png)

# 免责声明

本工具仅面向**合法授权**的企业安全建设行为，如您需要测试本工具的可用性，请自行搭建靶机环境。

为避免被恶意使用，本项目所有收录的poc均为漏洞的理论判断，不存在漏洞利用过程，不会对目标发起真实攻击和漏洞利用。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。**请勿对非授权目标进行扫描。**

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

在安装并使用本工具前，请您**务必审慎阅读、充分理解各条款内容**，限制、免责条款或者其他涉及您重大权益的条款可能会以加粗、加下划线等形式提示您重点注意。 除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。
