<script lang="ts" setup>
import { Picture } from '@element-plus/icons-vue'
import { highlightFingerprints } from '@/stores/style'
import { fp } from '@/composables/useWebscanState'
import { handleWebscanContextMenu } from '@/linkage/contextMenu'

function pictrueSRC(filepath: string): string {
    if (filepath == '') return ''
    const filename = filepath.split(/[/\\]/).pop();
    return `http://127.0.0.1:8732/screenhost/${filename}`;
}
</script>

<template>
    <el-table :data="fp.table.pageContent" 
        stripe height="100vh" 
        :highlight-current-row="true"
        :cell-style="{ textAlign: 'center' }"
        :header-cell-style="{ 'text-align': 'center' }"
        @row-contextmenu="handleWebscanContextMenu"
        @sort-change="fp.ctrl.sortChange">
        <el-table-column fixed prop="URL" label="Link" width="350px" />
        <el-table-column width="170px" label="Port & Protocol" :show-overflow-tooltip="true">
            <template #default="scope">
                <el-tag type="primary" round effect="plain">{{ scope.row.Port }}</el-tag>
                <el-tag type="primary" round effect="plain" class="ml-5px">{{ scope.row.Scheme }}</el-tag>
            </template>
        </el-table-column>
        <el-table-column prop="StatusCode" width="100px" label="Code" sortable="custom" />
        <el-table-column prop="Length" width="100px" label="Length" sortable="custom" />
        <el-table-column prop="Title" label="Title" width="220px" />
        <el-table-column prop="Fingerprints" label="Technologies" :min-width="300">
            <template #default="scope">
                <div class="finger-container">
                    <el-tag v-for="finger in scope.row.Fingerprints" :key="finger"
                        :effect="scope.row.Detect === 'Default' ? 'light' : 'dark'"
                        :type="highlightFingerprints(finger)">
                        {{ finger }}
                    </el-tag>
                    <el-tag type="warning" v-if="scope.row.IsWAF">{{ scope.row.WAF }}</el-tag>
                </div>
            </template>
        </el-table-column>
        <el-table-column label="Screen" width="150">
            <template #default="scope">
                <el-image :src="pictrueSRC(scope.row.Screenshot)"
                    :preview-src-list="[pictrueSRC(scope.row.Screenshot)]" :initial-index="0"
                    preview-teleported :max-scale="1" v-if="scope.row.Screenshot != ''">
                    <template #error>
                        <div class="image-slot">
                            <el-icon :size="16">
                                <Picture />
                            </el-icon>
                        </div>
                    </template>
                </el-image>
                <span v-else>-</span>
            </template>
        </el-table-column>
        <template #empty>
            <el-empty />
        </template>
    </el-table>
    <div class="flex-between mt-5px">
        <div></div>
        <el-pagination size="small" background @size-change="fp.ctrl.handleSizeChange"
            @current-change="fp.ctrl.handleCurrentChange" :pager-count="5"
            :current-page="fp.table.currentPage" :page-sizes="[10, 20, 50, 100, 200]"
            :page-size="fp.table.pageSize" layout="total, sizes, prev, pager, next"
            :total="fp.table.result.length">
        </el-pagination>
    </div>
</template>

