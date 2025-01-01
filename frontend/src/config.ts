import { CheckFileStat, InitConfig, ReadFile, RemoveOldConfig, SaveDataToFile, ReadLocalStore } from 'wailsjs/go/services/File';
import { ElLoading, ElNotification } from 'element-plus';
import global from "./stores";
import { compareVersion, ReadLine, sleep } from './util';
import router from "./router";
import { CreateTable } from 'wailsjs/go/services/Database';
import { crackDict } from './stores/options';

function catchError(result: boolean, loading: any) {
  if (result) {
    loading.setText("配置文件初始化成功!");
  } else {
    ElNotification({
      message: "无法下载配置文件,请自行到https://gitee.com/the-temperature-is-too-low/slack-poc/releases/下载config.zip并解压到用户根目录下的/slack/文件夹下!",
      position: "bottom-right"
    })
  }
}

export async function InitConfigFile(timeout: number) {
  LoadConfig();
  CreateTable()
  let backgroundColor = global.Theme.value ? "rgba(0, 0, 0)" : "rgba(255, 255, 255)"
  const loading = ElLoading.service({
    lock: true,
    background: backgroundColor,
    spinner: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 48 48" id="Slack-Logo--Streamline-Core.svg" height="48" width="48"><desc>Slack Logo Streamline Icon: https://streamlinehq.com</desc><g id="slack"><path id="Union" fill="#8fbffa" fill-rule="evenodd" d="M30.11753142857143 1.7142857142857142C27.46285714285714 1.7142857142857142 25.310811428571427 3.8663314285714283 25.310811428571427 6.521005714285714v11.36136c0 2.6546742857142855 2.1520457142857143 4.80672 4.80672 4.80672 2.6547085714285714 0 4.806582857142857 -2.1520457142857143 4.806582857142857 -4.80672V6.521005714285714C34.92411428571428 3.8663314285714283 32.77224 1.7142857142857142 30.11753142857143 1.7142857142857142ZM1.7142857142857142 17.882434285714282c0 2.6547085714285714 2.1520457142857143 4.806754285714285 4.80672 4.806754285714285h11.36136c2.6546742857142855 0 4.80672 -2.1520457142857143 4.80672 -4.806754285714285 0 -2.6546742857142855 -2.1520457142857143 -4.80672 -4.80672 -4.80672l-11.36136 0C3.8663314285714283 13.075714285714286 1.7142857142857142 15.22776 1.7142857142857142 17.882434285714282ZM22.68918857142857 41.478857142857144c0 2.654742857142857 -2.1520457142857143 4.806857142857142 -4.806754285714285 4.806857142857142 -2.6546742857142855 0 -4.80672 -2.1521142857142856 -4.80672 -4.806857142857142l0 -11.361222857142856c0 -2.6546742857142855 2.1520457142857143 -4.80672 4.80672 -4.80672 2.6547085714285714 0 4.806754285714285 2.1520457142857143 4.806754285714285 4.80672l0 11.361222857142856ZM46.285714285714285 30.991645714285713c0 -2.6546742857142855 -2.1521142857142856 -4.80672 -4.806857142857142 -4.80672H30.11763428571428c-2.6546742857142855 0 -4.80672 2.1520457142857143 -4.80672 4.80672 0 2.6546742857142855 2.1520457142857143 4.806754285714285 4.80672 4.806754285714285l11.361222857142856 0c2.654742857142857 0 4.806857142857142 -2.1520799999999998 4.806857142857142 -4.806754285714285Z" clip-rule="evenodd" stroke-width="1"></path><path id="Union_2" fill="#2859c5" fill-rule="evenodd" d="M18.320091428571427 1.7142857142857142c-2.4133714285714283 0 -4.3697485714285715 1.9564114285714285 -4.3697485714285715 4.3697485714285715s1.9563771428571426 4.3697485714285715 4.3697485714285715 4.3697485714285715h4.3697485714285715V6.084034285714285C22.68984 3.670697142857142 20.733428571428572 1.7142857142857142 18.320091428571427 1.7142857142857142ZM46.285714285714285 18.319234285714284c0 -2.413337142857143 -1.956342857142857 -4.3697485714285715 -4.369714285714285 -4.3697485714285715s-4.369714285714285 1.9564114285714285 -4.369714285714285 4.3697485714285715v4.3697485714285715h4.369714285714285c2.4133714285714283 0 4.369714285714285 -1.9563771428571426 4.369714285714285 -4.3697485714285715ZM34.04965714285714 41.916c0 2.4133714285714283 -1.9564114285714285 4.369714285714285 -4.3697485714285715 4.369714285714285s-4.3697485714285715 -1.956342857142857 -4.3697485714285715 -4.369714285714285V37.546285714285716h4.3697485714285715c2.413337142857143 0 4.3697485714285715 1.956342857142857 4.3697485714285715 4.369714285714285ZM1.7142857142857142 29.680765714285716c0 2.413337142857143 1.9564114285714285 4.3697485714285715 4.3697485714285715 4.3697485714285715s4.3697485714285715 -1.9564114285714285 4.3697485714285715 -4.3697485714285715l0 -4.3697485714285715H6.084034285714285C3.670697142857142 25.311017142857143 1.7142857142857142 27.267394285714282 1.7142857142857142 29.680765714285716Z" clip-rule="evenodd" stroke-width="1"></path></g></svg>`,
  })
  let cfgStat = await CheckFileStat(global.PATH.homedir + "/slack")
  if (!cfgStat) {
    loading.setText('未检测到配置文件，正在下载...')
    let result = await InitConfig();
    catchError(result, loading)
  } else {
    let result = await ReadFile(global.PATH.homedir + global.PATH.LocalPocVersionFile)
    global.UPDATE.LocalPocVersion = result.Content!
    if (!global.UPDATE.LocalPocVersion || compareVersion(global.UPDATE.LocalPocVersion, "0.0.4") != 1) {
      await RemoveOldConfig();
      loading.setText("正在下载新配置文件...");
      let result = await InitConfig();
      catchError(result, loading)
    }
  }
  // 联动
  loading.setText('正在初始化资源中...')
  const waitRouter = ["/Permeation/Webscan", "/"]
  for (const route of waitRouter) {
    await sleep(timeout)
    router.push(route);
  }
  loading.close();
}

// 加载本地配置信息
async function LoadConfig() {
  let stat = await CheckFileStat(global.PATH.homedir + "/slack/config.json")
  if (!stat) {
    var data = { proxy: global.proxy, space: global.space, jsfind: global.jsfind, webscan: global.webscan, database: global.database };
    await SaveDataToFile(data);
  } else {
    let result = await ReadLocalStore()
    Object.assign(global.proxy, result["proxy"])
    Object.assign(global.space, result["space"])
    Object.assign(global.jsfind, result["jsfind"])
    Object.assign(global.webscan, result["webscan"])
    Object.assign(global.database, result["database"])
    if (!result["webscan"]["highlight_fingerprints"] || global.database.columnsNameKeywords == "") {
      SaveConfig()
    }
  }
  await ReadCrackDict()
}

export function SaveConfig() {
  // 获取space的所有value值
  let list = Object.entries(global.space).map(([key, value]) => value);
  // 去除不可见字符
  list = list.map(item => item.replace(/[\r\n\s]/g, ''));
  var data = { proxy: global.proxy, space: global.space, jsfind: global.jsfind, webscan: global.webscan, database: global.database };
  SaveDataToFile(data).then(result => {
    if (result) {
      ElNotification.success({
        message: 'Save successful',
        position: 'bottom-right'
      })
    }
  })
};

async function ReadCrackDict() {
  for (var item of crackDict.usernames) {
    item.dic = (await ReadLine(global.PATH.homedir + global.PATH.PortBurstPath + "/username/" + item.name.toLowerCase() + ".txt"))!
  }
  crackDict.passwords = (await ReadLine(global.PATH.homedir + global.PATH.PortBurstPath + "/password/password.txt"))!
  crackDict.passwords.push("")
}