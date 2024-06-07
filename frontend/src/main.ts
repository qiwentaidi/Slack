import { createApp, App } from "vue";
import AppComponent from "./App.vue";
import "./style.css";
import router from "./router";
import i18n from './i18n/index' //引入配置的语言
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
import "element-plus/theme-chalk/el-message.css";
import "element-plus/theme-chalk/el-message-box.css";
import "element-plus/theme-chalk/el-notification.css";
import * as ElementPlusIconsVue from "@element-plus/icons-vue";
import { Splitpanes, Pane } from 'splitpanes'
import 'splitpanes/dist/splitpanes.css'

export default (app: App<Element>) => {
  // 全局配置
  app.config.globalProperties.$ELEMENT = {};
  app.use(ElMessage);
  app.use(ElMessageBox);
  app.use(ElNotification);
};

const app = createApp(AppComponent)
// 使得comonent 可以正确渲染el-icon
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.component('Splitpanes', Splitpanes);
app.component('Pane', Pane);

app.use(router).use(i18n).mount("#app");
