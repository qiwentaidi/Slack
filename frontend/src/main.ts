import { createApp, App } from 'vue';
import AppComponent from './App.vue';
import './style.css';
import { ElMessage, ElMessageBox } from 'element-plus';
import 'element-plus/theme-chalk/el-message.css';
import 'element-plus/theme-chalk/el-message-box.css';

export default (app: App<Element>) => {
  // 全局配置
  app.config.globalProperties.$ELEMENT = {};
  app.use(ElMessage);
  app.use(ElMessageBox);
};
createApp(AppComponent).mount('#app')

