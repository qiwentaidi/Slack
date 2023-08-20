# 使用说明

## 配置

由于有些`POC`需要反连平台验证，所以配置`ceye`这些还是比较重要的，`jndi`则是用来验证一些类似`Fastjson`的漏洞，我也不知道为什么`afrog`作者不用反连探测![image-20230814060208008](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230814060208008.png)

## Web扫描

`Web`扫描功能缝了`afrog`[项目地址](https://github.com/zan8in/afrog)，基本就是差不多把该扫描器的功能`UI`化了，主动探测的漏洞或者指纹会写入到`report`目录下的`html`文件中，并没有将`afrog`命令的输出进行删除（所以你在`go run main.go`运行工具时，依然能看到命令行存在`afrog`的输出内容）

指纹来源`Ehole`的全部指纹以及部分`Goby`指纹

### 主动指纹探测/指纹`POC`扫描

运行优先级，如果你勾选了仅指纹扫描，那么就不会打`POC`了，但是还是会进行敏感目录探测

- 主动指纹探测，用于发现一些hw中常见危害较高的敏感目录

- 指纹`POC`扫描避，免发送无用数据包，导致扫的很慢

![image-20230820045813847](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230820045813847.png)

![image-20230820065026410](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230820065026410.png)

### 自定义`POC`扫描

- 关键字（搜关键字）：

根据`yaml poc`的`id`值进行检测关键字

- 风险等级：

根据`yaml poc`的`severity`值进行检测风险等级

![image-20230820064124101](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230820064124101.png)

![image-20230814060624593](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230814060624593.png)

点击保存按钮，会将上半部分的结果写入到`report`文件夹下的`CSV`文件结果中（Q:为什么不保存下半部分的？A:因为会扫的时候自动生成）

![image-20230814064335713](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230814064335713.png)

双击`URL`链接也可以实现打开链接功能

**由于存在`afrog`的更新模块，他的`poc`有更新的时候会自动下载的用户根目录下，请自行查看更新以及添加**







## 端口扫描模块

勾选`ICMP`会先进行主机探活，端口指纹协议调用的`lcvvvv`大佬的`gonmap`库并添加了`JDWP`指纹识别，`Send HTTP To WebScan`按钮可以将扫描到的`HTTP or HTTPS`服务发送到Web扫描模块目标中

```
目标支持换行分割,IP支持如下格式:
192.168.1.1
192.168.1.1/8
192.168.1.1/16
192.168.1.1/24
192.168.1.1,192.168.1.2
192.168.1.1-192.168.255.255
192.168.1.1-255

仅支持下面格式排除IP:
!192.168.1.1

如果端口遗漏多请在配置中调高端口超时时长
```

![image-20230805001537793](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230805001537793.png)

## 资产收集模块

### 从企业名称进行信息收集

查看控股企业获得1级全资子公司的域名或者股权等信息，勾选查询`HUNTER`资产数量会查询资产数量，运行完毕后运行日志处会出现提示

![image-20230810152036444](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230810152036444.png)

### 子域名暴破

![image-20230820062824013](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230820062824013.png)

## 空间引擎

### 鹰图

整体风格和功能都跟`fofaview`类似，注意勾选上数据去重后需要消耗权益积分，没有权益积分的话别勾选。

`HUNTER`模块标签右键功能与`Web`扫描一致

![image-20230805000337553](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230805000337553.png)

![image-20230805000525983](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230805000525983.png)

![image-20230805000754835](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230805000754835.png)

导出方式分为

- 全部导出，需要消耗积分，导出最大数量，以当前可用积分为准
- 导出当前数据，不消耗积分，仅保存当前查询结果

![image-20230805001156867](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230805001156867.png)

### `FOFA`

`FOFA`是之前做的还没来得及改，后续会做成跟`HUNTER`模块一样（现在跟答辩一样可以不用）

![image-20230814070534358](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230814070534358.png)

## 其他工具

### 编码转换

从上到下进行模式加解密（比如我先添加了`Base64`，再添加`Urlencode`那程序会按添加顺序进行操作）。

- 加密：选择要添加的模式进行加解密，选择好加密模式后在上方文本框输入文本内容（无需点击`encode`按钮，本身就是个摆设按钮）即可自动完成加密功能。
- 解密：把需要解密的内容复制到下方文本框，点击`decode`按钮进行解密。

![GIF 2023-5-31 22-31-09](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/GIF 2023-5-31 22-31-09.gif)

### 杀软识别

将`tasklist or tasklist /svc`的内容复制到上方，点击识别 ，可以通过修改`/finger/antivirues.yaml`文件增加杀软指纹。

![image-20230805002333636](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230805002333636.png)

### `IP&CDN`

```
IP提取:
请输入任意内容，如果有IP地址会进行提取并统计C段数量

CDN识别:
输入任意内容，如果有域名会提取域名进行判断

CDN识别规则根据域名解析的IP进行判断,解析到非同C段IP判断为CDN
```

![image-20230810143306954](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230810143306954.png)

### `IP`转换

就支持下面两种转换格式

![image-20230805002634901](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230805002634901.png)

### `Fscan`内容提取

![image-20230810143000846](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230810143000846.png)

## 漏洞验证

目前支持验证漏洞

```
Redis未授权访问
Apache Hadoop Yarn RPC RCE
```

对于一些无回显的`POC`会提示执行反弹`shell`命令

![image-20230816185538003](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230816185538003.png)

![image-20230815191332939](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230815191332939.png)

![image-20230816184945520](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230816184945520.png)

## 备忘录

想写啥写啥，有需要的别忘了点保存

![image-20230805002856586](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230805002856586.png)

# 运行

由于该程序的`UI`库采用的`fyne`底层基于`C++`如果直接运行`main.go`文件则所以需要环境中先配置好`GCC`的环境才可以正常使用，下面这样就算配置完成了

![image-20230814065112447](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230814065112447.png)

# 中文显示问题&内存问题

由于`Fyne`内置不支持中文字体的原因，所以需要引入外部字体文件解决中文乱码问题，但是由于原生的`Windows`字体文件要么过大（会导致引入中文字体后程序运行内存占用增加的问题）要么就是显示不清晰的问题，所以需要压缩`TTF`字体文件来到达减小程序内存占用的问题，经过测试加载中文字体内存占用增加`50MB`。

1、使用`python3 pip` 安装`fonttools`工具` pip install fonttools`or`python3 -m pip install fonttools`检查是否安装成功。

![image-20230525234444144](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230525234444144.png)

2、选择ttf字体文件对其进行压缩（即只保留txt中的字符文件）`fonttools subset ".\Dengb.ttf" --text-file=".\汉语字典.txt" --output-file="ysfonts.ttf"`，前往`gui\mytheme\fonts`目录下将`ysfonts.tff`字体替换即可。

![image-20230525234533367](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230525234533367.png)

# 打包

使用`fyne`自带的打包模块，体积更小还能打包`logo`（`ico、png、jpg`都支持）

```
[go>=1.16]
go install fyne.io/fyne/v2/cmd/fyne@latest 
[go<1.16]
go get fyne.io/fyne/v2/cmd/fyne@latest    
```

- `windows`

  `fyne package -os windows -icon .\gui\mytheme\static\logo.ico .\main.go`

- 其他平台

  也可以使用`fyne package`
  
  `go build main.go`

# 一些模块的拓展规则

## `POC`编写

`config\afrog-pocs\README.md`文件下有详细的规则

## 指纹拓展

从`slack v1.1`开始舍弃了`iconmd5`指纹，规则如下所示，第一行为指纹名称，后续可以有3个键值对，分别为`title`匹配网页标题名，`iconhash`是`favicon`的`Mmh3Hash32`值可以通过`example/iconhash/main.go`文件进行计算，`header`匹配响应头中的内容，`yaml`规则尽量最外层是''即可不用转译双引号以及冒号等符号干扰...自行阅读`yaml`规则，不懂再问我

```yaml
Nexus Repository Manager: 
  title: 'Nexus Repository Manager'
  iconhash: '1323738809'
  header: 'NX-ANTI-CSRF-TOKEN'
```

## 杀软库拓展

匹配运行的进程名`finger/antivirues.yaml`

## `GoPOC`扩展

例如下面`Redis`未授权访问漏洞回显，需要传入`address`地址以及`cmd`需要执行的命令，执行成功后`common.SetTextRefresh(response, global.AttackResult)`刷新内容即可，第一个参数传输出的内容

```
func Redis(address, cmd string) {
	conn, err := net.DialTimeout("tcp4", address, time.Second*time.Duration(10))
	if err != nil {
		// 连接失败
		common.SetTextRefresh("连接失败\n", global.AttackResult)
		return
	}
	defer conn.Close()
	request := strings.TrimSpace(cmd) + "\r\n"
	if _, err = conn.Write([]byte(request)); err != nil {
		common.SetTextRefresh("发送数据包失败\n", global.AttackResult)
		return
	}
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		common.SetTextRefresh("读取数据包失败\n", global.AttackResult)
		return
	}
	response := string(buffer[:n])
	common.SetTextRefresh(response, global.AttackResult)
}

```

将写好的`go`文件放到`plugins/vulattack/`路径下，然后在`gui/module/attack.go`文件17行中找到`vul := []string`为其添加一个可以显示的漏洞名，如果无回显的漏洞则可以如文件26所示，增加一个`if`或的判断，给用户提示需要输入反弹`shell`命令

![image-20230817044008145](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230817044008145.png)

## 验证指纹是否可用

`example`目录下存放着一些测试文件，可以运行监测，特别是如果识别不到指纹，一定要检查`YAML`文件格式是否编写错误

# FAQ

- 指纹`POC`扫描如何实现

根据`config/workflow.yaml`中的定义的指纹名称以及`POC`数组的键值对进行攻击，先进行指纹扫描，在把主动和被动的指纹聚合到一个目标以及指纹数组的键值对中，等后续再从指纹数组中匹配要攻击的`POC`路径去进行扫描，后续需要拓展指纹要攻击的`POC`可以自行添加`workflow.yaml`文件中，如果是新增指纹同时加指纹`POC`，请保持指纹名称一致，但不区分大小写

![image-20230818001904001](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230818001904001.png)

其实整体逻辑是调用了两次`webscan.webscan`这个函数，通过改变参数让程序判断是进行判断是进行`POC`扫描还是指纹探测

![image-20230820063133233](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230820063133233.png)

![image-20230820063003172](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230820063003172.png)

- 如果识别不到指纹如何打选择`POC`进行攻击

利用关键字搜索`POC`进行扫描

- 为什么不做`SYN`扫描


![image-20230815144652592](https://qwtd-image.oss-cn-hangzhou.aliyuncs.com/img/image-20230815144652592.png)

# 更新日志

`v1.1`

```
1、在菜单栏中新增了主题配置的功能，可以调节字体大小以及主题颜色，为了防止黑色主题下图标不清晰情况，并已更改原先的图标(图标路径在gui\mytheme\static\下，未编译情况下可以替换文件自行更改)，切换主题会导致程序运行内存增加，重新启动即可

2、WEB扫描模块优化了代理检测，开始扫描时会优先检测代理连通性，再决定能否开始扫描

3、增加了漏洞验证模块，目前支持Redis未授权、Hadoop未授权RCE漏洞

4、优化了指纹拓展性，已将内置指纹更改为外部存放, 并重新进行整理(后续指纹等配置文件全部存储在config目录下，程序不便移动，需要放在对应目录下执行)

5、WEB扫描模块增加从指纹自动识别POC进行攻击

6、优化了一下&输出被转译等其他Bug
```

