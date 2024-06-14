<script setup lang="ts">
import { InitConfigFile } from "../config";
import global from "../global";
import { onMounted, ref, computed } from 'vue';
import Navigation from '../components/Navigation.vue'
import { Search } from "@element-plus/icons-vue";
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
// 初始化时调用
onMounted(async () => {
  await InitConfigFile(1000)
});

const searchFilter = ref('')
const activeCard = ref('local')

const filteredOptions = computed(() => {
  if (!searchFilter.value) {
    return global.onlineOptions;
  }
  return global.onlineOptions.map(group => ({
    ...group,
    value: group.value.filter(item => item.name.toLowerCase().includes(searchFilter.value.toLowerCase()))
  }));
});

</script>

<template>
  <div style="position: relative;">
    <el-tabs v-model="activeCard">
      <el-tab-pane :label="$t('aside.app_navigation')" name="local">
        <Navigation />
      </el-tab-pane>
      <el-tab-pane :label="$t('aside.navigation')" name="online">
        <el-scrollbar height="85vh">
          <div v-for="groups in filteredOptions" style="margin-top: 10px; width: 99%;">
            <el-card style="width: 100%">
              <div style="margin-bottom: 10px;"><span style="font-weight: bold;">{{ $t(groups.label) }}</span></div>
              <el-row>
                <el-col :lg="4" v-for="item in groups.value">
                  <el-tooltip :content="item.url" placement="top">
                    <el-result :title="item.name" @click="BrowserOpenURL(item.url)">
                      <template #icon>
                        <el-image class="onlineNavIcon" :src="item.icon" />
                      </template>
                    </el-result>
                  </el-tooltip>
                </el-col>
              </el-row>
            </el-card>
          </div>
        </el-scrollbar>

      </el-tab-pane>
    </el-tabs>
    <div class="custom_eltabs_titlebar" v-if="activeCard == 'online'">
      <el-input class="input-filter" :suffix-icon="Search" v-model="searchFilter" @input="filteredOptions"
        placeholder="Filter search"></el-input>
    </div>
  </div>
</template>

<style>
.input-filter {
  width: 240px;
  margin-right: 1%;
}
.onlineNavIcon .el-image__inner {
    height: 25px;
    width: 25px;
}
</style>