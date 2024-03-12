import { reactive } from 'vue'
var space = reactive({
    fofaapi: 'https://fofa.info/',
    fofaemail: '',
    fofakey: '',
    hunterkey: '',
    quakekey: ''
})
var scan = reactive({
    dns1: "114.114.114.114",
    dns2: "223.5.5.5"
})
var proxy = reactive({
    enabled: false,
    mode: 'HTTP',
    address: '127.0.0.1',
    port: 8080,
    username: '',
    password: '',
})
export default {
    space,
    scan,
    proxy,
};