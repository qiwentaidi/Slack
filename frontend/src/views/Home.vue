<script setup lang="ts">
import { InitConfigFile } from "../config";
import { onMounted } from "vue";
import MenuList from "../router/menu";
// 初始化时调用
onMounted(async () => {
  await InitConfigFile(1000);
});
</script>

<template>
  <el-scrollbar height="92vh">
    <template v-for="groups in MenuList">
      <div v-if="groups.children" class="card-group">
        <el-card style="margin-bottom: 5px;">
          <div style="margin-bottom: 10px">
            <span style="font-weight: bold">{{ $t(groups.name) }}</span>
          </div>
          <el-row>
            <el-col :lg="4" v-for="item in groups.children">
              <el-result :title="$t(item.name)" @click="$router.push(item.path)">
                <template #icon>
                  <el-image class="appNavIcon" :src="item.icon" />
                </template>
              </el-result>
            </el-col>
          </el-row>
        </el-card>
      </div>
    </template>
  </el-scrollbar>
</template>

<style>
.appNavIcon .el-image__inner {
  height: 30px;
  width: 30px;
}

.custom-tabs-label .el-icon {
  vertical-align: middle;
}

.custom-tabs-label span {
  vertical-align: middle;
  margin-left: 4px;
}
</style>
