<script setup lang="ts">
import { ref, onMounted } from "vue";
import global from "./global";
import Sidebar from "./components/Sidebar.vue";
import { useRoute } from "vue-router";
import { EventsOn } from "../wailsjs/runtime/runtime";
import { CheckFileStat, UserHomeDir, Mkdir, WriteFile } from "../wailsjs/go/main/File";
import router from "./router";
import { ElLoading } from "element-plus";

const route = useRoute();
const showLogger = ref(false);

// 自动加载联动路由，避免内容丢失
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

function breadcrumbItems(label: string, separator: string) {
  switch (label) {
    case "/":
      return "Home";
    case "/settings":
      return "Settings";
    default:
      return label.slice(1).split('/').join(separator);
  }
}

interface MsgInfo {
  Level: string
  Msg: string
}

const logArray = [] as string[]

onMounted(() => {
  autoLoadRoutes(1000)
  InitDict()
  const levelClassMap: { [key: string]: string } = {
    "[INF]": "log-info",
    "[WRN]": "log-warning",
    "[ERR]": "log-error",
    "[DEB]": "log-debug",
    "[SUC]": "log-success"
  };
  EventsOn("gologger", (mi: MsgInfo) => {
    const logClass = levelClassMap[mi.Level];
    const logEntry = `<span class="${logClass}">${mi.Level}</span> ${mi.Msg}`;
    logArray.push(logEntry);

    if (logArray.length > global.Logger.length) {
      logArray.shift(); // 移除数组开头的元素
    }
    // 将数组内容拼接成字符串
    global.Logger.value = logArray.join('\n');
  });
})


// 初始化字典
async function InitDict() {
  global.PATH.PortBurstPath = await UserHomeDir() + global.PATH.PortBurstPath
  if (!(await CheckFileStat(global.PATH.PortBurstPath))) {
    if (await Mkdir(global.PATH.PortBurstPath)) {
      Mkdir(global.PATH.PortBurstPath + "/username")
      Mkdir(global.PATH.PortBurstPath + "/password")
      for (const item of global.dict.usernames) {
        WriteFile("txt", `${global.PATH.PortBurstPath}/username/${item.name}.txt`, item.dic.join("\n"))
      }
      WriteFile("txt", `${global.PATH.PortBurstPath}/password/password.txt`, global.dict.passwords.join("\n"))
    }
  }
}
</script>

<template>
  <el-container>
    <el-aside>
      <Sidebar></Sidebar>
    </el-aside>
    <el-container>
      <el-main>
        <!-- 一定要使用插槽否则keey-alive不会生效 -->
        <router-view v-slot="{ Component }">
          <keep-alive>
            <component :is="Component"></component>
          </keep-alive>
        </router-view>
      </el-main>
      <div class="console-log">
        <div>
          <span class="margin-25">
            {{ breadcrumbItems(route.path, ' > ') }}
          </span>
          <div class="margin-25">
            <el-button link @click="showLogger = true">
              <template #icon>
                <img src="/console.svg">
              </template>
              Console
            </el-button>
          </div>
        </div>
      </div>
    </el-container>
  </el-container>
  <!-- running logs -->
  <el-drawer v-model="showLogger" direction="ltr" size="50%">
    <template #header>
      <h4>运行日志</h4>
    </template>
    <div class="log-textarea" v-html="global.Logger.value"></div>
  </el-drawer>
</template>

<style>
.el-aside {
  width: 64px;
}

.margin-25 {
  margin: 2.5vh;
  font-size: 14px;
  font-weight: 500;
}
</style>
