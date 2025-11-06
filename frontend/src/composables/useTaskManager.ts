import { ElMessage, ElMessageBox } from 'element-plus';
import { nanoid as nano } from 'nanoid'
import {
    RemoveScanTask, ExportWebReportWithHtml, ExportWebReportWithJson, RetrieveAllScanTasks,
    RenameScanTask, RetrieveFingerscanResults, RetrievePocscanResults, UpdateScanTaskWithResults,
    AddFingerscanResult, AddPocscanResult, AddScanTask, ReadWebReportWithJson, ExportWebReportWithExcel
} from 'wailsjs/go/services/Database';
import { FileDialog, SaveFileDialog } from 'wailsjs/go/services/File';
import { Callgologger } from 'wailsjs/go/services/App';
import { dashboard, form, fp, vp, rp, reportOption, reportName, exportDialog, historyDialog } from './useWebscanState'

export const taskManager = {
    generateUniqueTaskName: (baseName: string): string => {
        let name = baseName.trim()
        let counter = 1
        while (rp.table.result.some(item => item.TaskName === name)) {
            name = `${baseName}-${counter}`
            counter++
        }
        return name
    },
    writeTask: async function () {
        const uniqueName = taskManager.generateUniqueTaskName(form.taskName)
        form.taskId = nano()
        let isSuccess = await AddScanTask(form.taskId, uniqueName, form.input, 0, 0)
        if (!isSuccess) {
            ElMessage.error("添加任务失败")
            return false
        }
        rp.table.result.push({
            TaskId: form.taskId,
            TaskName: uniqueName,
            Targets: form.input,
            Failed: 0,
            Vulnerability: 0,
        })
        rp.ctrl.watchResultChange(rp.table)
        return true
    },
    viewTask: async function (row: any) {
        historyDialog.value = false;
        form.input = row.Targets;
        form.taskName = row.TaskName;
        form.taskId = row.TaskId
        const [fingerResult, nucleiResult] = await Promise.all([
            RetrieveFingerscanResults(row.TaskId),
            RetrievePocscanResults(row.TaskId)
        ]);

        if (fingerResult) {
            fp.table.result = fingerResult;
            fp.ctrl.watchResultChange(fp.table);
        }
        // 初始化风险等级计数
        Object.keys(dashboard.riskLevel).forEach(key => {
            dashboard.riskLevel[key as keyof typeof dashboard.riskLevel] = 0;
        });
        if (!nucleiResult) {
            vp.table.result = []
            vp.ctrl.watchResultChange(vp.table);
            return
        }
        vp.table.result = nucleiResult;
        vp.ctrl.watchResultChange(vp.table);
        // 遍历结果，统计每个风险等级的数量
        vp.table.result.forEach(item => {
            const riskLevelKey = item.Severity as keyof typeof dashboard.riskLevel;
            if (dashboard.riskLevel[riskLevelKey] !== undefined) {
                dashboard.riskLevel[riskLevelKey]++;
            }
        });
    },
    deleteTask: function (alltaskids: string[]) {
        ElMessageBox.confirm(
            "确定删除所选中的任务记录",
            '警告',
            {
                type: 'warning',
            }
        )
            .then(async () => {
                for (const taskid of alltaskids) {
                    let isSuccess = await RemoveScanTask(taskid)
                    if (!isSuccess) {
                        ElMessage.error(`任务ID: ${taskid}, 删除失败`)
                        return
                    }
                    rp.table.result = rp.table.result.filter(item => item.TaskId != taskid)
                }
                ElMessage.success("删除成功")
                rp.ctrl.watchResultChange(rp.table)
            })
            .catch(() => {
                Callgologger("error", "[webscan] delete task failed")
            })
    },
    renameTask: async function (taskId: string) {
        ElMessageBox.prompt('重命名任务', '编辑', {})
            .then(async ({ value }) => {
                let isSuccess = await RenameScanTask(taskId, value)
                if (!isSuccess) {
                    ElMessage.error("修改失败")
                }
                const taskIndex = rp.table.result.findIndex(task => task.TaskId === taskId);
                if (taskIndex !== -1) {
                    rp.table.result[taskIndex] = { ...rp.table.result[taskIndex], TaskName: value };
                    rp.ctrl.watchResultChange(rp.table);
                }
                ElMessage.success("修改成功")
            })
    },
    importTask: async function () {
        const filepath = await FileDialog("*.json")
        if (!filepath) {
            return
        }
        const result = await ReadWebReportWithJson(filepath)
        if (result) {
            const id = nano()
            let isSuccess = await AddScanTask(id, id, result.Targets, 0, 0)
            if (!isSuccess) {
                ElMessage.error("添加任务失败")
                return false
            }
            result.Fingerprints.forEach(item => {
                item.TaskId = id
                AddFingerscanResult(item)
            })
            if (result.POCs) {
                result.POCs.forEach(item => {
                    item.TaskId = id
                    AddPocscanResult(item)
                })
                rp.table.result.push({
                    TaskId: id,
                    TaskName: id,
                    Targets: result.Targets,
                    Failed: 0,
                    Vulnerability: result.POCs.length,
                })
            } else {
                rp.table.result.push({
                    TaskId: id,
                    TaskName: id,
                    Targets: result.Targets,
                    Failed: 0,
                    Vulnerability: 0,
                })
            }
            rp.ctrl.watchResultChange(rp.table)
            ElMessage.success("添加成功")
        }
    },
    showExportDialog: function () {
        reportName.value = ""
        if (rp.table.selectRows.length >= 2) {
            reportName.value = "合并报告-"
        }
        reportName.value += rp.table.selectRows[0].TaskName
        exportDialog.value = true
    },
    exportTask: async function () {
        let isSuccess = false
        let taskids = rp.table.selectRows.map(item => item.TaskId)
        let filepath = await SaveFileDialog(reportName.value)
        if (!filepath) {
            return
        }
        switch (reportOption.value) {
            case "EXCEL":
                isSuccess = await ExportWebReportWithExcel(filepath + ".xlsx", rp.table.selectRows)
                isSuccess ? ElMessage.success("导出成功") : ElMessage.error("导出失败")
                break
            case "JSON":
                isSuccess = await ExportWebReportWithJson(filepath + ".json", rp.table.selectRows)
                isSuccess ? ElMessage.success("导出成功") : ElMessage.error("导出失败")
                break
            default:
                isSuccess = await ExportWebReportWithHtml(filepath + ".html", taskids)
                isSuccess ? ElMessage.success("导出成功") : ElMessage.error("导出失败")
        }
        exportDialog.value = false
    },
    updateTaskTable: function (taskId: string) {
        let vulcount = dashboard.riskLevel.CRITICAL + dashboard.riskLevel.HIGH + dashboard.riskLevel.MEDIUM + dashboard.riskLevel.LOW + dashboard.riskLevel.INFO
        UpdateScanTaskWithResults(taskId, 0, vulcount)
        // 更新对应taskId的表格列中的漏洞数量
        const taskIndex = rp.table.result.findIndex(task => task.TaskId === taskId);
        if (taskIndex !== -1) {
            rp.table.result[taskIndex].Vulnerability = vulcount;
            rp.ctrl.watchResultChange(rp.table);
        }
    },
}

