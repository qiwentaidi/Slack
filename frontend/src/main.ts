import { createApp, App } from "vue";
import AppComponent from "./App.vue";
import "./style.css";
import router from "./router";
import i18n from './i18n/index' //引入配置的语言
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
import "element-plus/theme-chalk/el-message.css";
import "element-plus/theme-chalk/el-message-box.css";
import "element-plus/theme-chalk/el-notification.css";

// 声明全局变量通过window.调用
declare global {
  var ActivePathPoc: string
  var AFGPathPoc: string
  var PocVersion: string
  var LocalPocVersionFile: string
  var ConfigPath: string
  var HomePath: string
  var PortBurstPath: string
}


export default (app: App<Element>) => {
  // 全局配置
  app.config.globalProperties.$ELEMENT = {};
  app.use(ElMessage);
  app.use(ElMessageBox);
  app.use(ElNotification);
};

createApp(AppComponent).use(router).use(i18n).mount("#app");
