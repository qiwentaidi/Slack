import { ElMessage } from 'element-plus';
import { SpaceGetPort } from 'wailsjs/go/services/App'
import { ProcessTextAreaInput } from '@/util'
import { validateIp, isPrivateIP } from '@/stores/validate';
import { IPParse } from 'wailsjs/go/core/Tools';
import async from 'async'
import { form, spaceEngineConfig, shodanVisible, shodanIp, shodanPercentage, shodanThread, shodanRunningstatus } from './useWebscanState'
import { LinkFOFA, LinkHunter } from '@/linkage';

export const uncover = {
    fofa: async function () {
        spaceEngineConfig.fofaDialog = false
        let urls = await LinkFOFA(spaceEngineConfig.fofaQuery, spaceEngineConfig.fofaPageSize)
        if (urls) form.input = urls.join("\n")
    },
    hunter: async function () {
        spaceEngineConfig.hunterDialog = false
        let urls = await LinkHunter(spaceEngineConfig.hunterQuery, spaceEngineConfig.hunterPageSize)
        if (urls) form.input = urls.join("\n")
    },
    shodan: async function () {
        if (!validateIp(shodanIp.value)) {
            ElMessage.warning("目标输入格式不正确!")
            return
        }
        const lines = ProcessTextAreaInput(shodanIp.value)
        for (const line of lines) {
            if (isPrivateIP(line)) {
                ElMessage.warning(line + " 为内网地址不支持扫描!")
                return
            }
        }

        let id = 0
        form.input = ""
        shodanRunningstatus.value = true
        let ips = await IPParse(lines)
        async.eachLimit(ips, shodanThread.value, async (ip: string, callback: () => void) => {
            if (!shodanRunningstatus.value) {
                return
            }
            let ports = await SpaceGetPort(ip)
            id++
            if (ports == null) {
                return
            }
            for (const port of ports) {
                form.input += ip + ":" + port.toString() + "\n"
            }
        }, (err: any) => {
            shodanRunningstatus.value = false
            shodanVisible.value = false
        })
    }
}

