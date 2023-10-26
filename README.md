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

# 更新日志

`2023年10月26日 v1.4.4`

```
新增:
1、企业微信公众号查询功能
2、表格新增排序功能
3、对已扫描出来的漏洞增加查看详情功能(即查看请求和响应包，多请求只查看第一个)
4、漏洞版本已同步0.0.3

优化:
1、漏洞利用(会写GO POC可联系我添加)、目录扫描模块布局以及功能

修复:
1、fofa数据导出时，默认选项未选定的bug
2、Ascii码不能正确解密
3、漏洞更新忘记更新workflow、webfinger文件
```

`2023年10月22日 v1.4.3`

```
新增:
1、客户端更新功能
2、完整的POC以及工作流更新
3、漏洞详情功能

优化: 
1、一些布局以及图标
2、原先字体无法正确显示粗体、目前已经更换字体修复

修复: 
1、目录扫描在内网环境下无法正确探活地址导致程序退出
2、致远、帆软OA密码解密输入字符串太短导致程序退出
```

`2023年10月10日 v1.4.2`

```
新增:
1、360夸克空间测绘
2、微信APPID校验功能
3、可在窗口聚焦时通过快捷键CTRL+,调出主题设置 通过CTRL+L调出日志

修复:
1、目录扫描处增加URL校验，防止程序异常退出
2、优化CVE-2020-1957、CVE-2020-11989 POC误报问题
```

`2023年9月25日 v1.4.1`

```
新增:
1、目录扫描功能、时间戳转换小工具
2、新增首页导航栏，可以快速定位需要使用的模块

优化:
1、修复子域名暴破缓慢以及日志刷新过于频繁导致的程序退出
2、修复域名备案&whois查询时因为查询不到结果导致的下标越界
3、优化了表格控件目前支持左右拉动
4、优化了POC更新、IP纯真库下载以及日志功能，目前均处于菜单栏处
5、优化了任务下发时的工作逻辑，减少工作量
```

`2023年9月20日 v1.4`

```
新增:
1、新增控制台(右下角)，在大部分运行进度等日志会在控制台输出
2、鹰图模块增加批量查询功能
3、其他工具模块增加 
	(1) systeminfo提权补丁查询
	(2) 反弹shell生成器
	(3) 联想字典生成器
4、资产收集模块增加
	(1) 端口暴破("ftp", "ssh", "telnet", "smb", "mssql", "oracle", "mysql", "postgresql", "vnc", "redis", "mongodb")、支持联想模式 
	(2) 域名备案&whois查询
5、网站扫描功能，新增任务下发按键 beta版(仅需要输入公司名称即可完全全量信息收集以及漏洞扫描、端口暴破)

优化:
1、其他工具->IP&CDN模块、域名CDN识别规则重写(目前根据CNAME)中出现的关键字判断
2、子域名暴破去除了HTTP请求，规则库改为使用IP纯真库以及CNAME关键字判断
3、增加部分poc
```

`v1.3`

```
1、网站扫描功能增加搜索引擎联动，目前支持鹰图和FOFA

2、网站扫描功能支持指定文件夹POC进行扫描

3、取消反连配置，改为内置CEYE，JNDI相关POC也都已经更改为反连探测

4、增加单独检测FastJson漏洞功能
```

`v1.2`

```
1、完善了FOFA查询功能，目前可以跟鹰图一样使用

2、网站扫描：
	(1)优化了界面布局
	(2)增加停止扫描功能
	(3)目标支持x.x.x.x:80形式
	(4)指纹匹配规则进行了优化，之前会出现指纹匹配不上的问题
	(5)仅指纹扫描功能取消敏感目录探测
	
3、优化了备忘录界面布局

4、UI库全新升级界面更美观，控件做了圆角、以及阴影处理，新版本中表格控件目前存在bug，所以目前使用的仍然是老版本的代码

5、带有文件logo的文本框，可以将txt文件移入读取路径
```

`v1.1`

```
1、在菜单栏中新增了主题配置的功能，可以调节字体大小以及主题颜色，为了防止黑色主题下图标不清晰情况，并已更改原先的图标(图标路径在gui\mytheme\static\下，未编译情况下可以替换文件自行更改)，切换主题会导致程序运行内存增加，重新启动即可

2、WEB扫描模块优化了代理检测，开始扫描时会优先检测代理连通性，再决定能否开始扫描

3、增加了漏洞验证模块，目前支持Redis未授权、Hadoop未授权RCE漏洞

4、优化了指纹拓展性，已将内置指纹更改为外部存放, 并重新进行整理(后续指纹等配置文件全部存储在config目录下，程序不便移动，需要放在对应目录下执行)

5、WEB扫描模块增加从指纹自动识别POC进行攻击

6、优化了一下&输出被转译等其他Bug
```
