import { CheckFileStat, InitConfig, UserHomeDir } from '../wailsjs/go/main/File';
import { ElNotification } from 'element-plus';
import Loading from './components/Loading.vue';
import global from "./global";

export async function InitConfigFile() {
  LoadConfig() // 加载本地配置信息
  let cfg = await CheckFileStat(await UserHomeDir() + "/slack")
  if (!cfg) {
    ElNotification({
      duration: 0,
      message: '未检测到配置文件，正在初始化...',
      icon: Loading,
    });
    if (await InitConfig()) {
      ElNotification.closeAll()
      ElNotification({
        message: "配置文件初始化成功!",
        type: "success",
      });
    } else {
      ElNotification.closeAll()
      ElNotification({
        title: "无法下载配置文件",
        message: "请自行到https://gitee.com/the-temperature-is-too-low/slack-poc/releases/下载config.zip并解压到用户根目录下的/slack/文件夹下!",
        type: "warning",
      });
    }
  }
}

function LoadConfig() {
    const allLocaolStorage = [
      {
        key: "proxy",
        value: global.proxy,
      },
      {
        key: "space",
        value: global.space,
      },
      {
        key: "jsfind",
        value: global.jsfind,
      },
      {
        key: "webscan",
        value: global.webscan,
      }
    ];
    allLocaolStorage.forEach(item => {
      const v = localStorage.getItem(item.key)
      if (v) {
        Object.assign(item.value, JSON.parse(v));
      }
    });
  }