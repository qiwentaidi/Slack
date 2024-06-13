import { CheckFileStat, InitConfig, UserHomeDir, GetFileContent, RemoveOldConfig } from '../wailsjs/go/main/File';
import { ElLoading, ElNotification } from 'element-plus';
import Loading from './components/Loading.vue';
import global from "./global";
import { compareVersion } from './util';
import router from "./router";

var needJump = true

async function autoLoadRoutes(timeout: number) {
  const loading = ElLoading.service({
    lock: true,
    background: "rgba(255, 255, 255)",
    spinner: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 48 48" id="Slack-Logo--Streamline-Core.svg" height="48" width="48"><desc>Slack Logo Streamline Icon: https://streamlinehq.com</desc><g id="slack"><path id="Union" fill="#8fbffa" fill-rule="evenodd" d="M30.11753142857143 1.7142857142857142C27.46285714285714 1.7142857142857142 25.310811428571427 3.8663314285714283 25.310811428571427 6.521005714285714v11.36136c0 2.6546742857142855 2.1520457142857143 4.80672 4.80672 4.80672 2.6547085714285714 0 4.806582857142857 -2.1520457142857143 4.806582857142857 -4.80672V6.521005714285714C34.92411428571428 3.8663314285714283 32.77224 1.7142857142857142 30.11753142857143 1.7142857142857142ZM1.7142857142857142 17.882434285714282c0 2.6547085714285714 2.1520457142857143 4.806754285714285 4.80672 4.806754285714285h11.36136c2.6546742857142855 0 4.80672 -2.1520457142857143 4.80672 -4.806754285714285 0 -2.6546742857142855 -2.1520457142857143 -4.80672 -4.80672 -4.80672l-11.36136 0C3.8663314285714283 13.075714285714286 1.7142857142857142 15.22776 1.7142857142857142 17.882434285714282ZM22.68918857142857 41.478857142857144c0 2.654742857142857 -2.1520457142857143 4.806857142857142 -4.806754285714285 4.806857142857142 -2.6546742857142855 0 -4.80672 -2.1521142857142856 -4.80672 -4.806857142857142l0 -11.361222857142856c0 -2.6546742857142855 2.1520457142857143 -4.80672 4.80672 -4.80672 2.6547085714285714 0 4.806754285714285 2.1520457142857143 4.806754285714285 4.80672l0 11.361222857142856ZM46.285714285714285 30.991645714285713c0 -2.6546742857142855 -2.1521142857142856 -4.80672 -4.806857142857142 -4.80672H30.11763428571428c-2.6546742857142855 0 -4.80672 2.1520457142857143 -4.80672 4.80672 0 2.6546742857142855 2.1520457142857143 4.806754285714285 4.80672 4.806754285714285l11.361222857142856 0c2.654742857142857 0 4.806857142857142 -2.1520799999999998 4.806857142857142 -4.806754285714285Z" clip-rule="evenodd" stroke-width="1"></path><path id="Union_2" fill="#2859c5" fill-rule="evenodd" d="M18.320091428571427 1.7142857142857142c-2.4133714285714283 0 -4.3697485714285715 1.9564114285714285 -4.3697485714285715 4.3697485714285715s1.9563771428571426 4.3697485714285715 4.3697485714285715 4.3697485714285715h4.3697485714285715V6.084034285714285C22.68984 3.670697142857142 20.733428571428572 1.7142857142857142 18.320091428571427 1.7142857142857142ZM46.285714285714285 18.319234285714284c0 -2.413337142857143 -1.956342857142857 -4.3697485714285715 -4.369714285714285 -4.3697485714285715s-4.369714285714285 1.9564114285714285 -4.369714285714285 4.3697485714285715v4.3697485714285715h4.369714285714285c2.4133714285714283 0 4.369714285714285 -1.9563771428571426 4.369714285714285 -4.3697485714285715ZM34.04965714285714 41.916c0 2.4133714285714283 -1.9564114285714285 4.369714285714285 -4.3697485714285715 4.369714285714285s-4.3697485714285715 -1.956342857142857 -4.3697485714285715 -4.369714285714285V37.546285714285716h4.3697485714285715c2.413337142857143 0 4.3697485714285715 1.956342857142857 4.3697485714285715 4.369714285714285ZM1.7142857142857142 29.680765714285716c0 2.413337142857143 1.9564114285714285 4.3697485714285715 4.3697485714285715 4.3697485714285715s4.3697485714285715 -1.9564114285714285 4.3697485714285715 -4.3697485714285715l0 -4.3697485714285715H6.084034285714285C3.670697142857142 25.311017142857143 1.7142857142857142 27.267394285714282 1.7142857142857142 29.680765714285716Z" clip-rule="evenodd" stroke-width="1"></path></g></svg>`,
    text: '正在初始化联动模块中...',
  })
  await new Promise(resolve => setTimeout(resolve, timeout));
  router.push("/Permeation/Crack");

  await new Promise(resolve => setTimeout(resolve, timeout));
  router.push("/Permeation/Webscan");

  await new Promise(resolve => setTimeout(resolve, timeout));
  router.push("/");
  loading.close();
}

export async function InitConfigFile() {
  try {
    LoadConfig()
    // 判断配置文件是否存在
    let cfgStat = await CheckFileStat(await UserHomeDir() + "/slack")
    // 先进行判断配置文件版本是否存在且小于等于0.0.4
    let pcfgStat = await CheckFileStat(await UserHomeDir() + global.PATH.LocalPocVersionFile)
    if (pcfgStat) {
      global.UPDATE.LocalPocVersion = await GetFileContent(await UserHomeDir() + global.PATH.LocalPocVersionFile)
      if (compareVersion(global.UPDATE.LocalPocVersion, "0.0.4") != 1) {
        ElNotification("检测到旧版本配置文件，正在移除...")
        await RemoveOldConfig()
      }
    }else {
      ElNotification("未检测到version文件，正在移除...")
      await RemoveOldConfig()
    }
    if (!cfgStat) {
      ElNotification.closeAll()
      ElNotification({
        duration: 0,
        message: '未检测到配置文件，正在初始化...',
        icon: Loading,
      });
      if (await InitConfig()) {
        showNotification(true, "配置文件初始化成功!", "success")
      } else {
        showNotification(true, "无法下载配置文件,请自行到https://gitee.com/the-temperature-is-too-low/slack-poc/releases/下载config.zip并解压到用户根目录下的/slack/文件夹下!", "warning")
      }
    }
  } catch (err: any) {
    showNotification(true, "初始化配置文件失败，跳过初始化联动模块", "error")
    needJump = false
  } finally {
    if (needJump) {
      autoLoadRoutes(1000)
    }
  }
}

function showNotification(isClose: boolean, message: string, notiType: any){
  if (isClose) {
    ElNotification.closeAll()
  }
  ElNotification({
    message: message,
    type: notiType,
  })
}

// 加载本地配置信息
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