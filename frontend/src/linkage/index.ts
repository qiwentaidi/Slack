import { HunterSearch, FofaSearch } from 'wailsjs/go/services/App'
import global from '@/stores'
import { ElMessage, ElNotification } from 'element-plus';

export async function LinkHunter(query: string, count: string) {
    if (!global.space.hunterkey) {
        ElNotification.warning("请在设置处填写Hunter Key")
        return
    }
    ElMessage.info("正在查询鹰图数据，请稍后...")
    let result = await HunterSearch(global.space.hunterkey, query, count, "1", "0", "3", false)
    if (result.code !== 200) {
        if (result.code == 40205) {
            ElMessage(result.message)
        } else {
            ElMessage.error(result.message);
            return
        }
    }
    return result.data.arr.map(item => item.url);
}

export async function LinkFOFA(query: string, count: number) {
    if (global.space.fofakey == "" && global.space.fofaemail) {
        ElNotification.warning("请在设置处填写FOFA Key && FOFA Email")
        return
    }
    ElMessage.info("正在查询FOFA数据，请稍后...")
    let result = await FofaSearch(query, count.toString(), "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, true, true)
    if (result.Error) {
        ElMessage.warning(result.Message)
        return
    }
    return result.Results.map(item => item.URL)
}
