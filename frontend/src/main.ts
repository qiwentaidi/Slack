import { createApp, App } from "vue";
import AppComponent from "./App.vue";
import "./style.css";
import router from "./router";
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
import "element-plus/theme-chalk/el-message.css";
import "element-plus/theme-chalk/el-message-box.css";
import "element-plus/theme-chalk/el-notification.css";

// 声明全局变量通过window.调用
declare global {
  var ActivePathPoc: string
  var AFGPathPoc: string
  var PocVersion: string
  var LocalPocVersion: string
}


export default (app: App<Element>) => {
  // 全局配置
  app.config.globalProperties.$ELEMENT = {};
  app.use(ElMessage);
  app.use(ElMessageBox);
  app.use(ElNotification);
};
createApp(AppComponent).use(router).mount("#app");
