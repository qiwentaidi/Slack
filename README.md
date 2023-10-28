<h4 align="center">一款Go Fyne实现的GUI扫描器，功能涵盖网站扫描、端口扫描、企业信息收集、子域名暴破、空间引擎搜索、CDN识别等多功能的工具</h4>

<p align="center">
<img src="https://img.shields.io/github/go-mod/go-version/qiwentaidi/Slack?filename=go.mod">
<img src="https://img.shields.io/badge/fyne-v2.4.1-blue">
<a href="https://github.com/qiwentaidi/Slack/releases/"><img src="https://img.shields.io/github/v/release/qiwentaidi/Slack">
<a href="https://github.com/qiwentaidi/Slack/releases/"><img src="https://img.shields.io/github/downloads/qiwentaidi/Slack/total">
</p>

# 使用说明

- 首页：所见即所得，点击对应的按钮会跳转到相应的模块

![image-20231026111120103](assets/image-20231026111120103.png)

## 渗透测试

- 网站扫描：网站扫描功能参考[AFROG](https://github.com/zan8in/afrog)，主动探测的漏洞或者指纹会写入到`reports`目录下的`html`文件中

![image-20231021224058649](assets/image-20231021224058649.png)

- 端口/指纹扫描：端口指纹协议调用的[gonmap](https://github.com/lcvvvv/gonmap)库

![image-20231026230025454](assets/image-20231026230025454.png)

- 端口暴破

![image-20231026230212953](assets/image-20231026230212953.png)

- 目录扫描

![image-20231026225023887](assets/image-20231026225023887.png)

- 漏洞利用

![image-20231026230317412](assets/image-20231026230317412.png)

- 漏洞详情

![image-20231022141004206](assets/image-20231022141004206.png)

## 资产收集模块

- 公司名称查资产：菜单栏-日志处可查看任务进度

![image-20231026233640475](assets/image-20231026233640475.png)

![image-20231023163359750](assets/image-20231023163359750.png)

- 子域名暴破：使用了IP纯真库以及CName关键字匹配进行CDN校验，注意如果网络环境差会导致暴破的很慢，这个自行考究

![image-20231006123920861](assets/image-20231006123920861.png)

- 域名信息收集

![image-20231006124136876](assets/image-20231006124136876.png)

## 空间引擎

- FOFA

- 360夸克

- 鹰图

三者整体风格和功能都跟`fofaview`类似只介绍一种，选项上都会备注是否需要会员等信息。

![image-20231026234307850](assets/image-20231026234307850.png)

![image-20230920110727981](assets/image-20230920110727981.png)

导出方式分为

- 全部导出，需要消耗积分，导出最大数量，以当前可用积分为准
- 导出当前数据，不消耗积分，仅保存当前查询结果

![image-20231010195929984](assets/image-20231010195929984.png)

## 小工具

- 编码转换

- 杀软识别&补丁检测

![image-20230920123610349](assets/image-20230920123610349.png)

- `Fscan`内容提取：除了常见的一些结果，另外添加了vcenter和海康摄像头的识别

![image-20231006124412626](assets/image-20231006124412626.png)

- 反弹`shell`生成器

- 联想字典生成器

- 备忘录

# 运行&中文显示问题&内存问题

由于该程序的`UI`库采用的`fyne`底层基于`C++`如果直接运行`main.go`文件则所以需要环境中先配置好`GCC`的环境才可以正常使用，还需要存在`go`环境，**本工具不适配Win7及以下系统**

[详情查看]: https://github.com/fyne-io/fyne/issues/1923

由于`Fyne`原生不支持中文字体的原因，所以需要引入外部字体文件解决中文乱码问题，但是直接引入完整的`Windows`字体文件要么（会导致引入中文字体后程序运行内存占用增加的问题）要么就是显示不清晰的问题，所以需要压缩`TTF(暂不支持TTC)`字体文件来到达减小程序内存占用的问题，经过测试加载中文字体内存占用增加`50MB`。

1、使用`python3 pip` 安装`fonttools`工具` pip install fonttools`or`python3 -m pip install fonttools`检查是否安装成功。

![image-20230525234444144](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230525234444144.png)

2、选择ttf字体文件对其进行压缩（即只保留txt中的字符文件）`fonttools subset ".\Dengb.ttf" --text-file=".\汉语字典.txt" --output-file="ysfonts.ttf"`，前往`gui\mytheme\fonts`目录下将`ysfonts.tff`字体替换即可。

![image-20230525234533367](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230525234533367.png)

# 打包

使用`fyne`自带的打包模块，体积更小还能打包`logo`（`ico、png、jpg`都支持），后续可以用`upx`进行压缩体积

```
[go>=1.16]
go install fyne.io/fyne/v2/cmd/fyne@latest 
[go<1.16]
go get fyne.io/fyne/v2/cmd/fyne@latest    
```

需要存在`FyneApp.toml`文件

- `windows`

  `fyne package -os windows .\main.go`

- 其他平台

  也可以使用`fyne package`或者`go build main.go`


# 一些模块的拓展规则

## 指纹拓展

规则如下所示，第一行为指纹名称，后续可以有3个键值对，分别为`title`匹配网页标题名，`iconhash`是`favicon`的`Mmh3Hash32`值可以通过`fofa模块进行hash计算`，`header`匹配响应头中的内容，`yaml`规则尽量最外层是''即可不用转译双引号以及冒号等符号干扰（所有文本仅先按照 || 条件进行的分割，因为没有用cel表达式）

```yaml
腾达路由器:
  title: '(Tenda && LOGIN) || (Tenda && Router) || (Tenda && Login) || (Tenda && Wi-Fi)  || (Tenda && 登录) || (Tenda && Web Master)) || 腾达无线路由器'
  iconhash: '-2145085239'
Nexus Repository Manager: 
  title: 'Nexus Repository Manager'
  iconhash: '1323738809 || -1546574541'
  header: 'NX-ANTI-CSRF-TOKEN'
Nginx:
  title: "Welcome to nginx!"
  header: 'Server:nginx'
```

## `POC`编写

`config\afrog-pocs\README.md`文件下有详细的规则

由于指纹`POC`根据`config/workflow.yaml`中的定义的指纹名称以及`POC`数组的值对进行攻击，后续需要拓展指纹要攻击的`POC`可以自行添加`workflow.yaml`文件中，如果是新增指纹同时加指纹`POC`，请保持指纹名称一致，但不区分大小写

![image-20230818001904001](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230818001904001.png)

## 杀软库拓展

匹配运行的进程名`finger/antivirues.yaml`

## 验证指纹是否可用

`example`目录下存放着一些测试文件，可以运行监测，特别是如果识别不到指纹，一定要检查`YAML`文件格式是否编写错误

# 联系方式

如果有问题可以加我联系方式（请备注来意）进工具交流群

![image-20231006124944803](assets/image-20231006124944803.png)

# 免责声明

本工具仅面向**合法授权**的企业安全建设行为，如您需要测试本工具的可用性，请自行搭建靶机环境。

为避免被恶意使用，本项目所有收录的poc均为漏洞的理论判断，不存在漏洞利用过程，不会对目标发起真实攻击和漏洞利用。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。**请勿对非授权目标进行扫描。**

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

在安装并使用本工具前，请您**务必审慎阅读、充分理解各条款内容**，限制、免责条款或者其他涉及您重大权益的条款可能会以加粗、加下划线等形式提示您重点注意。 除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。
