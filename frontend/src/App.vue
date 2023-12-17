<script setup lang="ts">
import { ref, provide } from 'vue';
import Header from './components/Header.vue'
import Introduction from './components/Introduction.vue';
import Settings from './components/Settings.vue';
import Webscan from './components/Permeation/Webscan.vue';
import Ipdomain from './components/Asset/Ipdomain.vue';
import Wxappid from './components/Tools/Wxappid.vue';
import Fscan from './components/Tools/Fscan.vue';
import Reverse from './components/Tools/Reverse.vue';
import Codec from './components/Tools/Codec.vue';
import Systems from './components/Tools/System.vue';
import Thinkdict from './components/Tools/Thinkdict.vue';
import Subdomain from './components/Asset/Subdomain.vue';
import Asset from './components/Asset/Asset.vue';
import Dirsearch from './components/Permeation/Dirsearch.vue';
import Portscan from './components/Permeation/Portscan.vue';
import Pocdetail from './components/Permeation/Pocdetail.vue';
import Fofa from './components/SpaceEngine/Fofa.vue';

const currentComponent = ref('Introduction');
provide('currentComponent', currentComponent);

const components: { [key: string]: any } = {
  'Introduction': Introduction,
  'Settings': Settings,
  'Webscan': Webscan,
  'Wxappid': Wxappid,
  'Fscan': Fscan,
  'Reverse': Reverse,
  'Codec': Codec,
  'System': Systems,
  'Thinkdict': Thinkdict,
  'Ipdomain': Ipdomain,
  'Subdomain': Subdomain,
  'Asset': Asset,
  'Dirsearch': Dirsearch,
  'Portscan': Portscan,
  'Pocdetail': Pocdetail,
  'Fofa': Fofa,
};

</script>

<template>
  <el-container>
    <Vueroute></Vueroute>
    <!-- 菜单栏 -->
    <el-header>
      <Header></Header>
    </el-header>
    <!-- 主内容 -->
    <div class="main">
      <KeepAlive>
        <component :is="components[currentComponent]"></component>
      </KeepAlive>
    </div>
  </el-container>
</template>

<style scoped>
.main {
  margin-left: 2%;
  margin-right: 2%;
}
</style>
