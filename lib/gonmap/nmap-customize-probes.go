package gonmap

var nmapCustomizeProbes = `
Probe TCP SMB_NEGOTIATE q|\x00\x00\x00\xc0\xfeSMB@\x00\x00\x00\x00\x00\x00\x00\x00\x00\x1f\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00$\x00\b\x00\x01\x00\x00\x00\u007f\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00x\x00\x00\x00\x02\x00\x00\x00\x02\x02\x10\x02"\x02$\x02\x00\x03\x02\x03\x10\x03\x11\x03\x00\x00\x00\x00\x01\x00&\x00\x00\x00\x00\x00\x01\x00 \x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x03\x00\x0e\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00|
rarity 1
ports 445

match microsoft-ds m|^\0\0...SMB.*|s

Probe TCP JSON_RPC q|{"id":1,"jsonrpc":"2.0","method":"login","params":{}}\r\n|
rarity 4
ports 443,80,8443,8080

match jsonrpc m|^{"jsonrpc":"([\d.]+)".*"height":(\d+),"seed_hash".*|s v/$1/ p/ETH/ i/height:$2/
match jsonrpc m|^{"jsonrpc":"([\d.]+)".*|s v/$1/
# Java RMI Registry probe
# RMI协议的JRMP握手消息
Probe TCP JavaRMI q|\x4a\x52\x4d\x49\x00\x02\x4b|
rarity 7
ports 1098,1099,1090,10999,4444,9010,8999

# 标准RMI registry响应
match rmi m|^\x4e\x00\x09\x31\x32\x37\x2e\x30\x2e\x30\x2e\x31\x00\x00| p/Java RMI Registry/ i/Standard RMI service/

# 通用RMI响应匹配
match rmi m|^\x4e| p/Java RMI Registry/
`
