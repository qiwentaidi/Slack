<script setup lang="ts">
import MenuList from "@/router/menu";

</script>

<template>
  <el-scrollbar height="92vh">
    <div v-for="groups in MenuList">
      <div v-if="groups.children" class="card-group">
        <el-divider>{{ $t(groups.name) }}</el-divider>
        <!-- <el-card style="margin-bottom: 5px;"> -->
          <!-- <div style="margin-bottom: 10px">
            <span style="font-weight: bold">{{ $t(groups.name) }}</span>
          </div> -->
          <el-row>
            <el-col :lg="6" v-for="item in groups.children">
              <el-result :title="$t(item.name)" @click="$router.push(item.path)">
                <template #icon>
                  <el-image class="appNavIcon" :src="item.icon" />
                </template>
              </el-result>
            </el-col>
          </el-row>
        <!-- </el-card> -->
      </div>
    </div>
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
