import bugIcon from '@/assets/icon/bug.svg'
import fingerprintIcon from '@/assets/icon/fingerprint.svg'
import deepscanIcon from '@/assets/icon/deepscan.svg'
import bulleyeIcon from '@/assets/icon/bulleye.svg'
import scanIcon from "@/assets/icon/scan.svg"
import themeIcon from "@/assets/icon/theme.svg"
import proxyIcon from "@/assets/icon/proxy.svg"
import dictmanagerIcon from "@/assets/icon/dict.svg"
import layersIcon from "@/assets/icon/layers.svg"
import aboutIcon from "@/assets/icon/about.svg"
import scriptIcon from "@/assets/icon/script.svg"
import maxmizeIcon from "@/assets/icon/maximize.svg"
import reductionIcon from "@/assets/icon/reduction.svg"
import pocIcon from '@/assets/icon/pocmanagement.svg'
import consoleIcon from '@/assets/icon/console.svg'
import htmlIcon from '@/assets/icon/html.svg'
import jsonIcon from '@/assets/icon/json.svg'
import excleIcon from '@/assets/icon/excle.svg'
import { Back, Right, RefreshRight, Minus, Close, Refresh, Setting, DataBoard } from '@element-plus/icons-vue';
import { WindowReload, WindowToggleMaximise, Quit, WindowMinimise } from "wailsjs/runtime/runtime";
import { computed } from "vue";
import global from "./index";

export const webscanOptions = [
    {
        label: "指纹扫描",
        value: 0,
        icon: fingerprintIcon,
    },
    {
        label: "全指纹扫描",
        value: 1,
        icon: deepscanIcon,
    },
    {
        label: "指纹漏洞扫描",
        value: 2,
        icon: bulleyeIcon,
    },
    {
        label: "专项扫描",
        value: 3,
        icon: bugIcon,
    },
]

export const webReportOptions = [
    {
        label: "HTML",
        icon: htmlIcon,
    },
    {
        label: "EXCEL",
        icon: excleIcon,
    },
    {
        label: "JSON",
        icon: jsonIcon,
    }
]

export const sortSeverityOptions = ["CRITICAL", "HIGH", "MEDIUM", "LOW", "INFO"]

export const quakeSyntaxOptions = [
    {
        title: "基本信息",
        data: [
            { key: 'ip:1.1.1.1/22', description: '基本信息' },
            { key: 'is_ipv6:true', description: '查询IPv6数据' },
            { key: 'port:80', description: '查询开放80端口的资产' },
            { key: 'port:[50 TO 60]', description: '查询开放端口在50至60之间的资产' },
            { key: 'transport:udp', description: '查询udp数据' },
            { key: 'domain:google.com', description: '查询资产域名' },
            { key: 'asn:12345', description: '查询自治域号码为"12345"的资产' },
            { key: 'org:No.31,Jin-rong Street', description: '查询自治域归属组织为"No.31,Jin-rong Street"的资产' },
            { key: 'hostname:unifiedlayer.com', description: '查询主机名包含"unifiedlayer.com"的资产' }
        ]
    },
    {
        title: "服务数据",
        data: [
            { key: 'service:http', description: '查询服务名称为"http"的资产' },
            { key: 'app:Apache', description: '查询应用名称为"Apache"的资产' },
            { key: 'response:奇虎科技', description: '查询端口原生返回数据中包含"奇虎科技"的资产' },
            { key: 'cert:奇虎科技', description: '查询包含"奇虎科技"的证书资产' },
            { key: 'catalog:IoT物联网', description: '查询应用类别为"IoT物联网"的资产', isVip: true },
            { key: 'type:VPN OR 防火墙', description: '查询应用类型为"VPN 或 防火墙"的资产', isVip: true },
            { key: 'level:硬件设备层', description: '查询应用层级为"硬件设备层"', isVip: true },
            { key: 'vendor:Sangfor OR 微软', description: '查询应用生产商为"Sangfor 或 微软"的资产', isVip: true }
        ]
    },
    {
        title: "IP归属与定位",
        data: [
            { key: 'country:China', description: '搜索国家，可使用简写CN' },
            { key: 'country_cn:中国', description: '搜索中文国家名称' },
            { key: 'province:Sichuan', description: '搜索英文省份名称' },
            { key: 'province_cn:四川', description: '搜索中文省份名称' },
            { key: 'city:Chengdu', description: '搜索英文城市名称' },
            { key: 'city_cn:成都', description: '搜索中文城市名称' },
            { key: 'owner:tencent.com', description: '搜索IP归属单位' },
            { key: 'isp:amazon.com', description: '搜索IP归属运营商' },
        ]
    },
    {
        title: "网页深度识别",
        data: [
            { key: 'icp:京ICP备08010314号', description: '查询ICP备案号', isVip: true },
            { key: 'icp_nature:企业', description: '查询备案主体性质' },
            { key: 'icp_keywords:奇虎', description: '查询备案网站中的关键词或域名', isVip: true },
            { key: 'link_url:/login', description: '查询网页url中包含"/login"的资产', isVip: true },
            { key: 'script_variable:csrfToken', description: '查询script标签变量中包含"csrfToken"的资产', isVip: true },
            { key: 'script_function:indexOf', description: '查询script标签函数中包含"indexOf"的资产', isVip: true },
            { key: 'css_class:login', description: '查询css标签class选择器中包含"login"的资产', isVip: true },
            { key: 'css_id:login', description: '查询css标签id选择器中包含"login"的资产', isVip: true },
            { key: 'iframe_title:企业微信登录', description: '查询iframe链接标题中包含"Enterprise wechat login"的资产', isVip: true },
            { key: 'iframe_url:/wwopen/sso/v1/', description: '查询iframe链接中包含"/wwopen/sso/v1/"的资产', isVip: true },
            { key: 'url_load:/login', description: '查询网页加载流中包含"/login"的资产', isVip: true },
            { key: 'page_type:门户网站', description: '查询类型是"门户网站"的资产', isVip: true }
        ]
    },
    {
        title: "布尔逻辑",
        data: [
            { key: 'port:"80" AND app:" 群晖NAS"', description: '查询开放80端口的"Synology NAS"应用资产' },
            { key: 'title:" 管理系统" OR body:" 登录页"', description: '查询网页标题包含"management system"或网页正文包含"登录页"的网站' },
            { key: 'title:"登录" AND(NOT title:"管理系统")', description: '查询网页标题包含"login"但不包含"管理系统"的网站' },
        ]
    }
]

export const hunterOptions = ({
    Server: [
        {
            value: '3',
            label: '全部资产',
        },
        {
            value: '1',
            label: 'WEB服务资产',
        },
        {
            value: '2',
            label: '非WEB服务资产',
        },
    ],
    Time: [
        {
            value: '0',
            label: '最近一个月',
        },
        {
            value: '1',
            label: '最近半年',
        },
        {
            value: '2',
            label: '最近一年',
        },
    ],
    Syntax: [
        {
            title: "IP",
            data: [
                { key: 'ip="1.1.1.1"', description: '搜索IP为 ”1.1.1.1”的资产' },
                { key: 'ip="220.181.111.0/24"', description: '搜索网段为"220.181.111.0"的C段资产' },
                { key: 'ip.port="80"', description: '	搜索开放端口为”80“的资产' },
                { key: 'ip.country="中国" 或 ip.country="CN"', description: '搜索IP对应主机所在国为”中国“的资产' },
                { key: 'ip.province="江苏"', description: '搜索IP对应主机在江苏省的资产' },
                { key: 'ip.city="北京"', description: '搜索IP对应主机所在城市为”北京“市的资产' },
                { key: 'ip.isp="电信"', description: '搜索运营商为”中国电信”的资产' },
                { key: 'ip.os="Windows"', description: '搜索操作系统标记为”Windows“的资产' },
                { key: 'app="Hikvision 海康威视 Firmware 5.0+" && ip.ports="8000"', description: '检索使用了Hikvision且ip开放8000端口的资产' },
                { key: 'ip.port_count>"2"', description: '搜索开放端口大于2的IP（支持等于、大于、小于）', hot: true },
                { key: 'ip.ports="80" && ip.ports="443"', description: '查询开放了80和443端口号的资产' },
                { key: 'ip.tag="CDN" ', description: '查询包含IP标签"CDN"的资产', hot: true, characteristic: true }
            ]
        },
        {
            title: "domain域名",
            data: [
                { key: 'is_domain=true', description: '搜索域名标记不为空的资产' },
                { key: 'domain="qianxin.com"', description: '搜索域名包含"qianxin.com"的网站' },
                { key: 'domain.suffix="qianxin.com"', description: '搜索主域为"qianxin.com"的网站', hot: true },
                { key: 'domain.status="clientDeleteProhibited"', description: '搜索域名状态为"client Delete Prohibited"的网站' },
                { key: 'domain.whois_server="whois.markmonitor.com"', description: '搜索whois服务器为"whois.markmonitor.com"的网站' },
                { key: 'domain.name_server="ns1.qq.com"', description: '搜索名称服务器为"ns1.qq.com"的网站' },
                { key: 'domain.created_date="2022-06-01"', description: '搜索域名创建时间为"2022-06-01"的网站' },
                { key: 'domain.expires_date="2022-06-01"', description: '搜索域名到期时间为"2022-06-01"的网站' },
                { key: 'domain.updated_date="2022-06-01"', description: '搜索域名更新时间为"2022-06-01"的网站' },
                { key: 'domain.cname="a6c56dbcc1f22283.qaxanyu.com"', description: '搜索cname包含“a6c56dbcc1f22283.qaxanyu.com”的网站' },
                { key: 'is_domain.cname=true', description: '搜索含有cname解析记录的网站' },
            ]
        },
        {
            title: "header响应头",
            data: [
                { key: 'header.server=="Microsoft-IIS/10"', description: '搜索server全名为“Microsoft-IIS/10”的服务器' },
                { key: 'header.content_length="691"', description: '搜索HTTP消息主体的大小为691的网站' },
                { key: 'header.status_code="402"', description: '搜索HTTP请求返回状态码为”402”的资产' },
                { key: 'header="elastic"', description: '搜索HTTP响应头中含有”elastic“的资产' },
            ]
        },
        {
            title: "web网站信息",
            data: [
                { key: 'is_web=true', description: '搜索web资产', hot: true },
                { key: 'web.title="北京"', description: '从网站标题中搜索“北京”' },
                { key: 'web.body="网络空间测绘"', description: '搜索网站正文包含”网络空间测绘“的资产' },
                { key: 'after="2021-01-01" && before="2021-12-31"', description: '搜索2021年的资产' },
                { key: 'web.similar="baidu.com:443"', description: '查询与baidu.com:443网站的特征相似的资产', hot: true, characteristic: true },
                { key: 'web.similar_icon=="17262739310191283300"', description: '查询网站icon与该icon相似的资产', hot: true, characteristic: true },
                { key: 'web.icon="22eeab765346f14faf564a4709f98548"', description: '查询网站icon与该icon相同的资产', hot: true },
                { key: 'web.similar_id="3322dfb483ea6fd250b29de488969b35"', description: '查询与该网页相似的资产', hot: true, characteristic: true },
                { key: 'web.tag="登录页面"', description: '查询包含资产标签"登录页面"的资产', hot: true, characteristic: true },
            ]
        },
        {
            title: "icp备案信息",
            data: [
                { key: 'icp.number="京ICP备16020626号-8"', description: '搜索通过域名关联的ICP备案号为”京ICP备16020626号-8”的网站资产' },
                { key: 'icp.web_name="奇安信"', description: '搜索ICP备案网站名中含有“奇安信”的资产' },
                { key: 'icp.name="奇安信"', description: '搜索ICP备案单位名中含有“奇安信”的资产' },
                { key: 'icp.type="企业"', description: '搜索ICP备案主体为“企业”的资产' },
                { key: 'icp.industry="软件和信息技术服务业"', description: '搜索ICP备案行业为“软件和信息技术服务业”的资产' },
                { key: 'icp.is_exception=true', description: '搜索含有ICP备案异常的资产', characteristic: true },
            ]
        },
        {
            title: "protocol协议/端口响应",
            data: [
                { key: 'protocol="http"', description: '搜索协议为”http“的资产' },
                { key: 'protocol.transport="udp"', description: '搜索传输层协议为”udp“的资产' },
                { key: 'protocol.banner="nginx"', description: '查询端口响应中包含"nginx"的资产' },
            ]
        },
        {
            title: "app组件信息",
            data: [
                { key: 'app.name="小米 Router"', description: '搜索标记为”小米 Router“的资产' },
                { key: 'app.type="开发与运维"', description: '查询包含组件分类为"开发与运维"的资产' },
                { key: 'app.vendor="PHP"', description: '查询包含组件厂商为"PHP"的资产' },
                { key: 'app.version="1.8.1"', description: '查询包含组件版本为"1.8.1"的资产' },
            ]
        },
        {
            title: "cert证书",
            data: [
                { key: 'cert="baidu"', description: '搜索证书中带有baidu的资产' },
                { key: 'cert.subject="qianxin.com"', description: '搜索证书使用者包含qianxin.com的资产' },
                { key: 'cert.subject.suffix="qianxin.com"', description: '搜索证书使用者为qianxin.com的资产' },
                { key: 'cert.subject_org="奇安信科技集团股份有限公司"', description: '搜索证书使用者组织是奇安信科技集团股份有限公司的资产' },
                { key: 'cert.issuer="Let\'s Encrypt Authority X3"', description: '搜索证书颁发者是Let\'s Encrypt Authority X3的资产' },
                { key: 'cert.issuer_org="Let\'s Encrypt"', description: '搜索证书颁发者组织是Let\'s Encrypt的资产' },
                { key: 'cert.sha-1="be7605a3b72b60fcaa6c58b6896b9e2e7442ec50"', description: '搜索证书签名哈希算法sha1为be7605a3b72b60fcaa6c58b6896b9e2e7442ec50的资产' },
                { key: 'cert.sha-256="4e529a65512029d77a28cbe694c7dad1e60f98b5cb89bf2aa329233acacc174e"', description: '搜索证书签名哈希算法sha256为4e529a65512029d77a28cbe694c7dad1e60f98b5cb89bf2aa329233acacc174e的资产' },
                { key: 'cert.sha-md5="aeedfb3c1c26b90d08537523bbb16bf1"', description: '搜索证书签名哈希算法shamd5为aeedfb3c1c26b90d08537523bbb16bf1的资产' },
                { key: 'cert.serial_number="35351242533515273557482149369"', description: '搜索证书序列号是35351242533515273557482149369的资产' },
                { key: 'cert.is_expired=true', description: '搜索证书已过期的资产' },
                { key: 'cert.is_trust=true', description: '搜索证书可信的资产', hot: true },
            ]
        },
    ],
})

export const fofaOptions = ({
    Advanced: [
        {
            key: '=',
            description: '匹配，=""时，可查询不存在字段或者值为空的情况。',
        },
        {
            key: '==',
            description: '完全匹配，==""时，可查询存在且值为空的情况。',
        },
        {
            key: '&&',
            description: '与',
        },
        {
            key: '||',
            description: '或',
        },
        {
            key: '!=',
            description: '不匹配，!=""时，可查询值不为空的情况。',
        },
        {
            key: '*=',
            description: '模糊匹配，使用*或者?进行搜索，比如banner*="mys??"。',
            level: "个人版"
        },
        {
            key: '()',
            description: '确认查询优先级，括号内容优先级最高。',
        }
    ],
    Syntax: [
        {
            title: "基础类",
            data: [
                { key: 'ip="1.1.1.1"', description: '通过单一IPv4地址进行查询', filed1: "✓", filed2: "✓", filed3: "-" },
                { key: 'port="6379"', description: '通过端口号进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'domain="qq.com"', description: '通过根域名进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'host=".fofa.info"', description: '通过主机名进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'os="centos"', description: '通过操作系统进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'server="Microsoft-IIS/10"', description: '通过服务器进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'asn="19551"', description: '通过自治系统号进行搜索', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'org="LLC Baxet"', description: '通过所属组织进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'is_domain=(true/false)', description: '筛选(拥有/没有)域名的资产', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'is_ipv6=(true/false)', description: '筛选是(ipv6/ipv4)的资产', filed1: "✓", filed2: "-", filed3: "-" },
            ]
        },
        {
            title: "标记类",
            data: [
                { key: 'app="Microsoft-Exchange"', description: '通过FOFA整理的规则进行查询', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'fid="sSXXGNUO2FefBTcCLIT/2Q=="', description: '通过FOFA聚合的站点指纹进行查询', filed1: "✓", filed2: "✓", filed3: "-" },
                { key: 'product="NGINX"', description: '通过FOFA标记的产品名进行查询', filed1: "✓", filed2: "✓", filed3: "-" },
                { key: 'category="服务"', description: '通过FOFA标记的分类进行查询', filed1: "✓", filed2: "✓", filed3: "-" },
                { key: 'type="subdomain"', description: '筛选服务（网站类）资产', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'cloud_name="Aliyundun"', description: '通过云服务商进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'is_cloud=(true/false)', description: '筛选(是/不是)云服务的资产', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'is_fraud=(true/false)', description: '筛选(是/不是)仿冒垃圾站群的资产', filed1: "✓", filed2: "-", filed3: "-", level: "专业版" },
                { key: 'is_honeypot=(true/false)', description: '筛选(是/不是)蜜罐的资产', filed1: "✓", filed2: "-", filed3: "-", level: "专业版" },
            ]
        },
        {
            title: "协议类",
            data: [
                { key: 'protocol="quic"', description: '通过协议名称进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'banner="users"', description: '通过协议返回信息进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'base_protocol="tcp"', description: '查询传输层为tcp协议的资产', filed1: "✓", filed2: "✓", filed3: "-" },
            ]
        },
        {
            title: "证书类",
            data: [
                { key: 'cert="baidu"', description: '通过证书进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.subject="Oracle Corporation"', description: '通过证书的持有者进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.issuer="DigiCert"', description: '通过证书的颁发者进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.subject.org="Oracle Corporation"', description: '通过证书持有者的组织进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.subject.cn="baidu.com"', description: '通过证书持有者的通用名称进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.issuer.org="cPanel, Inc."', description: '通过证书颁发者的组织进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.issuer.cn="Synology Inc. CA"', description: '通过证书颁发者的通用名称进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.is_valid=(true/false)"', description: '筛选证书(是/不是)有效证书的资产', filed1: "✓", filed2: "-", filed3: "-", level: "个人版" },
                { key: 'cert.is_match=(true/false)', description: '筛选证书和域名(匹配/不匹配)的资产', filed1: "✓", filed2: "-", filed3: "-", level: "个人版" },
                { key: 'cert.is_expired=(true/false)', description: '筛选证书(已过期/未过期)的资产', filed1: "✓", filed2: "-", filed3: "-", level: "个人版" },
                { key: 'tls.version="TLS 1.3"', description: '通过tls的协议版本进行查询', filed1: "✓", filed2: "✓", filed3: "-" },
            ]
        },
        {
            title: "时间类",
            data: [
                { key: 'after="2023-01-01"', description: '筛选某一时间之后有更新的资产', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'before="2023-12-01"', description: '筛选某一时间之前有更新的资产', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'after="2023-01-01" && before="2023-12-01"', description: '筛选某一时间区间有更新的资产', filed1: "✓", filed2: "-", filed3: "-" }
            ]
        },
        {
            title: "独立IP语法（不可和上面其他语法共用）",
            data: [
                { key: 'port_size="6"', description: '筛选开放端口数量等于6个的独立IP', filed1: "✓", filed2: "✓", filed3: "-", level: "个人版" },
                { key: 'port_size_gt="6"', description: '筛选开放端口数量大于6个的独立IP', filed1: "✓", filed2: "-", filed3: "-", level: "个人版" },
                { key: 'port_size_lt="12"', description: '筛选开放端口数量小于12个的独立IP', filed1: "✓", filed2: "-", filed3: "-", level: "个人版" },
                { key: 'ip_ports="80,161"', description: '筛选同时开放不同端口的独立IP', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'ip_country="CN"', description: '通过国家的简称代码进行查询独立IP', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'ip_region="Zhejiang"', description: '通过省份/地区英文名称进行查询独立IP', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'ip_city="Hangzhou"', description: '通过城市英文名称进行查询独立IP', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'ip_after="2021-03-18"', description: '筛选某一时间之后有更新的独立IP', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'ip_before="2019-09-09"', description: '筛选某一时间之前有更新的独立IP', filed1: "✓", filed2: "-", filed3: "-" }
            ]
        }
    ],
})

export var crackDict = ({
    // 可以进行漏洞检测的列表，包括未授权
    options: ["ftp", "ssh", "telnet", "smb", "oracle", "mssql", "mysql", "postgresql", "vnc", "redis", "memcached", "mongodb", "ldap", "mqtt", "socks5", "jdwp", "rmi"],
    usernames: [
        {
            name: "FTP",
            dic: [] as string[]
        },
        {
            name: "SSH",
            dic: [] as string[]
        },
        {
            name: "Telnet",
            dic: [] as string[]
        },
        {
            name: "LDAP",
            dic: [] as string[]
        },
        {
            name: "SMB",
            dic: [] as string[]
        },
        {
            name: "SOCKS5",
            dic: [] as string[]
        },
        {
            name: "MQTT",
            dic: [] as string[]
        },
        {
            name: "Mssql",
            dic: [] as string[]
        },
        {
            name: "Oracle",
            dic: [] as string[]
        },
        {
            name: "Mysql",
            dic: [] as string[]
        },
        {
            name: "RDP",
            dic: [] as string[]
        },
        {
            name: "Postgresql",
            dic: [] as string[]
        },
        {
            name: "VNC",
            dic: [] as string[]
        },
        {
            name: "Mongodb",
            dic: [] as string[]
        }
    ],
    passwords: [] as string[],
})

export const portGroupOptions = [
    {
        text: "数据库端口",
        value: "1433,1521,3306,5432,6379,9200,11211,27017",
    },
    {
        text: "企业端口",
        value: "21,22,80,81,135,139,443,445,1433,1521,3306,5432,6379,7001,8000,8080,8089,9000,9200,11211,27017,80,81,82,83,84,85,86,87,88,89,90,91,92,98,99,443,800,801,808,880,888,889,1000,1010,1080,1081,1082,1099,1118,1888,2008,2020,2100,2375,2379,3000,3008,3128,3505,5555,6080,6648,6868,7000,7001,7002,7003,7004,7005,7007,7008,7070,7071,7074,7078,7080,7088,7200,7680,7687,7688,7777,7890,8000,8001,8002,8003,8004,8006,8008,8009,8010,8011,8012,8016,8018,8020,8028,8030,8038,8042,8044,8046,8048,8053,8060,8069,8070,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8108,8118,8161,8172,8180,8181,8200,8222,8244,8258,8280,8288,8300,8360,8443,8448,8484,8800,8834,8838,8848,8858,8868,8879,8880,8881,8888,8899,8983,8989,9000,9001,9002,9008,9010,9043,9060,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9100,9200,9443,9448,9800,9981,9986,9988,9998,9999,10000,10001,10002,10004,10008,10010,10250,12018,12443,14000,16080,18000,18001,18002,18004,18008,18080,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018,20880",
    },
    {
        text: "高危端口",
        value: "21,22,23,53,80,443,8080,8000,139,445,3389,1521,3306,6379,7001,2375,27017,11211",
    },
    {
        text: "全端口",
        value: "1-65535"
    },
    {
        text: "两高一弱",
        value: "21,22,23,25,53,69,110,111,135,139,143,161,389,445,873,1025,1433,1158,1521,3306,3389,3690,5432,5900,5901,6379,7001,7002,9000,9043,9200,9300,27017,27018,28017,50030,50060,50070,1099,2049,2181,2222,2375,2379,2888,3128,3888,4000,4040,4440,4848,4899,5000,5005,5601,5631,5632,5984,6123,7051,7077,7180,7182,7848,8019,8020,8042,8048,8051,8069,8080,8081,8083,8086,8088,8161,8443,8649,8848,8880,8888,9001,9042,9083,9092,9100,9990,10000,11000,11111,11211,18080,19888,20880,25000,25010,50000,50090,60000,60010,60030"
    },
    {
        text: "自定义",
        value: ""
    },
]

export const aliveGroupOptions = [
    {
        value: "None",
        description: "不进行主机存活探测直接扫描端口"
    },
    {
        value: "ICMP",
        description: "使用ICMP进行主机存活探测，需要开启ROOT权限"
    },
    {
        value: "Ping",
        description: "使用Ping进行主机存活探测，ICMP权限不足时会自动降级为Ping"
    },
]

export const setupOptions = [
    {
        name: 'setting.scan',
        icon: scanIcon,
    },
    {
        name: 'setting.proxy',
        icon: proxyIcon,
    },
    {
        name: 'setting.mapping',
        icon: layersIcon,
    },
    {
        name: 'setting.display',
        icon: themeIcon,
    },
    {
        name: 'setting.dict',
        icon: dictmanagerIcon,
    },
    // {
    //   name: 'aside.script',
    //   icon: scriptIcon,
    // },
    {
        name: 'setting.about',
        icon: aboutIcon,
    },
]

export const wechatResponseDescription = [
    {
        code: "-1",
        describe: "system error",
        solution: "系统繁忙，此时请开发者稍候再试",
    },
    {
        code: "40001",
        describe: "invalid credential  access_token isinvalid or not latest",
        solution:
            "获取 access_token 时 AppSecret 错误，或者 access_token 无效。请开发者认真比对 AppSecret 的正确性，或查看是否正在为恰当的公众号调用接口",
    },
    {
        code: "40013",
        describe: "invalid appid",
        solution:
            "不合法的 AppID ，请开发者检查 AppID 的正确性，避免异常字符，注意大小写",
    },
    {
        code: "40002",
        describe: "invalid grant_type",
        solution: "不合法的凭证类型",
    },
    {
        code: "40125",
        describe: "不合法的 secret",
        solution: "请检查 secret 的正确性，避免异常字符，注意大小写",
    },
    {
        code: "40164",
        describe: "调用接口的IP地址不在白名单中",
        solution: "请在接口IP白名单中进行设置",
    },
    {
        code: "41004",
        describe: "appsecret missing",
        solution: "缺少 secret 参数",
    },
    {
        code: "50004",
        describe: "禁止使用 token 接口",
        solution: "",
    },
    {
        code: "50007",
        describe: "账号已冻结",
        solution: "",
    },
    {
        code: "61024",
        describe: "第三方平台 API 需要使用第三方平台专用 token",
        solution: "",
    },
    {
        code: "40243",
        describe: "AppSecret已被冻结，请登录小程序平台解冻后再次调用。",
        solution: "",
    },
]

export const dnsServerOptions = [
    {
        label: "谷歌",
        value: "8.8.8.8:53"
    },
    {
        label: "谷歌",
        value: "8.8.4.4:53"
    },
    {
        label: "‌电信",
        value: "114.114.114.114:53"
    },
    {
        label: "百度",
        value: "180.76.76.76:53"
    },
    {
        label: "腾讯",
        value: "119.29.29.29:53"
    },
    {
        label: "阿里",
        value: "223.6.6.6:53"
    }
]

export const subdomainRunnerOptions = [
    {
        label: "枚举模式",
        value: 0,
        tips: "不推荐使用"
    },
    {
        label: "查询模式",
        value: 1,
        tips: "通过API查询"
    },
    {
        label: "混合模式",
        value: 2,
        tips: "先查询后枚举"
    },
]

export const uncoverSyntaxOptions = [
    { filed: "IP", fofa: "✓", hunter: "✓", quake: "✓" },
    { filed: "域名", fofa: "✓", hunter: "✓", quake: "✓" },
    { filed: "标题", fofa: "✓", hunter: "✓", quake: "✓" },
    { filed: "Body", fofa: "✓", hunter: "✓", quake: "✓" },
    { filed: "备案名称", fofa: "✓", hunter: "✓", quake: "VIP" },
    { filed: "备案号", fofa: "✓", hunter: "✓", quake: "VIP" },
    { filed: "条件拼接数量", fofa: "100", hunter: "5", quake: "1000" },
    { filed: "单页API查询数", fofa: "10000", hunter: "100", quake: "500" }
]

export const databaseOptions = [
    {
        label: "MySQL",
        value: "mysql"
    },
    {
        label: "SQL Server",
        value: "mssql"
    },
    {
        label: "Oracle",
        value: "oracle"
    },
    {
        label: "PostgreSQL",
        value: "postgres"
    },
    {
        label: "Mongodb",
        value: "mongodb"
    }
]

export const pocdetailFilterOptions = [
    {
        label: '名称',
        value: 'Name'
    },
    {
        label: '关联指纹',
        value: 'Fingerprint'
    },
]

// 普通微信公众号
export const wechatApiOptions = [
    {
        name: "查询域名配置",
        method: "POST",
        url: "https://api.weixin.qq.com/wxa/getwxadevinfo?access_token=",
    },
    {
        name: "获取长期订阅用户",
        method: "POST",
        url: "https://api.weixin.qq.com/wxa/business/get_wxa_followers?access_token=",
    },
    {
        name: "获取用户列表(1w)",
        method: "POST",
        url: "https://api.weixin.qq.com/cgi-bin/user/get?count=10000&access_token=",
    },
    {
        name: "获取用户反馈列表",
        method: "GET",
        url: "https://api.weixin.qq.com/wxaapi/feedback/list?access_token=",
    },
]

// 企业微信

export const enterpriseWechatApiOptions = [
    {
        name: "获取成员ID列表(1000)",
        method: "POST",
        url: "https://qyapi.weixin.qq.com/cgi-bin/user/list_id?access_token=",
        body: `{"limit": 10000}`
    },
    {
        name: "获取部门列表",
        method: "GET",
        url: "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=",
        body: ``
    },
    {
        name: "获取部门的成员信息",
        method: "GET",
        url: "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist?department_id=1&access_token=",
        body: ``
    }
]

export const dingtalkApiOptions = [
    {
        name: "获取员工人数",
        method: "POST",
        url: "https://oapi.dingtalk.com/topapi/user/count?access_token=",
        body: `{"only_active":"true"}`,
    },
    {
        name: "获取管理员列表",
        method: "GET",
        url: "https://oapi.dingtalk.com/topapi/user/listadmin?access_token=",
        body: ``
    },
    {
        name: "获取部门用户完整信息(100)",
        method: "POST",
        url: "https://oapi.dingtalk.com/topapi/v2/user/list?access_token=",
        body: `{"dept_id":1,"cursor":0,"size":100}`
    }
]

// export const dgworkApiOptions = [

// ]

export const ApiDocsOptions = [
    {
        show: "wechat",
        label: "API文档",
        link: "https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html"
    },
    {
        show: "wechat",
        label: "官方调试工具",
        link: "https://mp.weixin.qq.com/debug/"
    },
    {
        show: "enterprise wechat",
        label: "API文档",
        link: "https://developer.work.weixin.qq.com/document/path/90664"
    },
    {
        show: "dingtalk",
        label: "API文档",
        link: "https://open.dingtalk.com/document/orgapp/api-overview"
    },
    {
        show: "dgwork",
        label: "API文档",
        link: "https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2674834"
    },
]

export const routerControl = [
    {
        label: "titlebar.back",
        icon: Back,
        action: () => {
            window.history.back();
        },
    },
    {
        label: "titlebar.forward",
        icon: Right,
        action: () => {
            window.history.forward();
        }
    },
    {
        label: "titlebar.reload",
        icon: RefreshRight,
        action: () => {
            WindowReload()
        }
    },
]

export const windowsControl = computed(() => [
    {
        icon: Minus,
        action: () => {
            WindowMinimise();
        },
    },
    {
        icon: global.temp.isMax ? reductionIcon : maxmizeIcon,
        action: () => {
            WindowToggleMaximise();
        },
    },
    {
        icon: Close,
        action: () => {
            Quit();
        },
        class: 'close',
    },
]);

export const sidebarBottomControl = [
    {
        label: "aside.update",
        icon: Refresh,
        action: () => {
            global.UPDATE.updateDialog = true
        }
    },
    {
        label: "aside.poc_manage",
        icon: pocIcon,
        path: "/PocManagement",
    },
    {
        label: "aside.setting",
        icon: Setting,
        path: "/Settings"
    },
]


export const JSFindOptions = [
    {
        label: "运行日志",
        value: 0,
        icon: consoleIcon,
    },
    {
        label: "数据展示",
        value: 1,
        icon: DataBoard
    }
]