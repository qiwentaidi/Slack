<script lang="ts" setup>
import { ElMessage, ElMessageBox } from "element-plus";
import { reactive, ref } from "vue";
import { LocalOpitons, Child } from "../interface";
import global from "../global";
import { DeleteFilled, Edit, FolderOpened, Document } from "@element-plus/icons-vue";
import { onMounted } from "vue";
import { OnFileDrop, EventsOn } from "../../wailsjs/runtime/runtime";
import { Path, GetLocalNaConfig, InsetGroupNavigation, InsetItemNavigation, OpenFolder, SaveNavigation, RunApp, FileDialog } from "../../wailsjs/go/main/File";
import { GOOS } from "../../wailsjs/go/main/App";
import { useI18n } from 'vue-i18n';
const { t } = useI18n();

onMounted(async () => {
    let os = await GOOS()
    // windows
    OnFileDrop((x, y, paths) => {
        if (os == "windows") {
            let card = localGroup.options.value.find(item => item.Name === config.mouseOnGroupName)
            paths.forEach(p => {
                Path(p).then((pathinfo: any) => {
                    if (card) {
                        let c = {
                            Name: pathinfo.Name,
                            Type: localGroup.getExtType(pathinfo.Ext),
                            Path: p
                        }
                        if (!card.Children) {
                            card.Children = [c];
                        } else {
                            card.Children.push(c)
                        }
                        InsetItemNavigation(config.mouseOnGroupName, c)
                    }
                })
            })
        }
    }, true)
    GetLocalNaConfig().then((result: LocalOpitons[]) => {
        if (result) {
            localGroup.options.value.push(...result)
        }
    })
    // macos
    EventsOn("wails-drop", (paths: string[]) => {
        let card = localGroup.options.value.find(item => item.Name === config.mouseOnGroupName)
        paths.forEach(p => {
            Path(p).then((pathinfo: any) => {
                if (card) {
                    let c = {
                        Name: pathinfo.Name,
                        Type: localGroup.getExtType(pathinfo.Ext),
                        Path: p
                    }
                    if (!card.Children) {
                        card.Children = [c];
                    } else {
                        card.Children.push(c)
                    }
                    InsetItemNavigation(config.mouseOnGroupName, c)
                }
            })
        })
    })
})

const config = reactive({
    defaultType: "CMD",
    defualtGroupName: "",
    name: "",
    path: "",
    mouseOnGroupName: "", // 鼠标移入时的组名
    editDialog: false,
    editName: "",
    editPath: "",
    editType: "",
    editChild: {} as Child,
    editGroupName: "", // 正在被编辑的组名
})

const localGroup = ({
    options: ref([] as LocalOpitons[]),
    openGroup: ["CMD", "APP", "JAR"],
    getGroupNames: function () {
        return localGroup.options.value.map(item => item.Name)
    },
    getSvgPath: function (type: string) {
        switch (type) {
            case "JAR":
                return "/navigation/java.svg"
            case "APP":
                return "/navigation/app.svg"
            default:
                return "/navigation/iterm.svg"
        }
    },
    getExtType: function (type: string) {
        if (type == "JAR") {
            return type
        } else if (type == "LNK" || type == "EXE") {
            return "APP"
        }
        return "CMD"
    },
    addGroup: function () {
        ElMessageBox.prompt('Name', t("navigator.add_group"), {
            confirmButtonText: 'OK',
            cancelButtonText: 'Cancel',
            inputPattern: /.+/,
            inputErrorMessage: "Group name can't be empty",
        })
            .then(({ value }) => {
                let existingGroup = localGroup.options.value.find(item => item.Name == value)?.Name
                if (existingGroup) {
                    ElMessage({
                        type: "warning",
                        message: "A group name with the same name already exists, please rename it"
                    })
                    return
                } else {
                    localGroup.options.value.push({
                        Name: value,
                        Children: null
                    })
                    InsetGroupNavigation({
                        Name: value,
                        Children: null
                    })
                }
            })
    },
    addItem: async function () {
        let card = localGroup.options.value.find(item => item.Name === config.defualtGroupName)
        if (card) {
            let child = {
                Name: config.name,
                Type: config.defaultType,
                Path: config.path
            }
            if (!card.Children) {
                card.Children = [child];
            } else {
                card.Children.push(child)
            }
            InsetItemNavigation(config.defualtGroupName, child)
        }
    },
    deleteGroup: function (name: string) {
        localGroup.options.value = localGroup.options.value.filter(item => item.Name !== name)
        SaveNavigation(localGroup.options.value)
    },
    deleteItem: function (groupName: string, child: Child) {
        localGroup.options.value = localGroup.options.value.map(group => {
            if (group.Name === groupName) {
                return {
                    ...group,
                    Children: group.Children?.filter(c => c.Name !== child.Name || c.Type !== child.Type || c.Path !== child.Path) || null
                };
            }
            return group;
        });
        SaveNavigation(localGroup.options.value)
    },
    editItem: function (groupName: string, child: Child) {
        config.editDialog = true
        config.editName = child.Name
        config.editType = child.Type
        config.editPath = child.Path
        config.editChild = child
        config.editGroupName = groupName
    },
    saveEdit: function () {
        const groupIndex = localGroup.options.value.findIndex(group =>
            group.Name == config.editGroupName
        );
        if (groupIndex !== -1) { // If the group is found
            localGroup.options.value[groupIndex].Children!.forEach((item, index) => {
                if (item.Name == config.editChild.Name && item.Type == config.editChild.Type && item.Path == config.editChild.Path) {
                    localGroup.options.value[groupIndex].Children![index].Name = config.editName;
                    localGroup.options.value[groupIndex].Children![index].Type = config.editType;
                    localGroup.options.value[groupIndex].Children![index].Path = config.editPath;
                    SaveNavigation(localGroup.options.value);
                    config.editDialog = false;
                }
            });
        }
    },
    handleDrop: function (event: DragEvent, groupName: string) {
        event.preventDefault();
        config.mouseOnGroupName = groupName
    },
    handleOpenFolder: async function (filepath: string) {
        let result = await OpenFolder(filepath)
        if (result != "") {
            ElMessage({
                type: "error",
                message: result
            })
        }
    },
    selectFile: async function () {
        let result = await FileDialog("")
        config.path = result
    },
})

// 暴露addGroup方法供其他地方调用
defineExpose({
    addGroup: localGroup.addGroup
})

const showCardPopover = ref(false)
const xRef = ref(0)
const yRef = ref(0)
function handleCardContextMenu(event: MouseEvent) {
    event.preventDefault()
    showCardPopover.value = true
    xRef.value = event.clientX
    yRef.value = event.clientY
}
</script>


<template>
    <div v-for="groups in localGroup.options.value" class="card-group">
        <el-card @drop="(event: any) => localGroup.handleDrop(event, groups.Name)" class="drop-enable"
            @contextmenu.prevent="handleCardContextMenu">
            <div class="my-header" style="margin-bottom: 10px">
                <span style="font-weight: bold">{{ groups.Name }}</span>
                <el-popconfirm title="Are you sure to delete this?" @confirm="localGroup.deleteGroup(groups.Name)">
                    <template #reference>
                        <el-button :icon="DeleteFilled" link></el-button>
                    </template>
                </el-popconfirm>
            </div>
            <div v-if="groups.Children" class="button-grid">
                <div v-for="item in groups.Children">
                    <el-popover ref="popover" placement="right" :width="200" trigger="contextmenu">
                        <el-menu class="right-click">
                            <el-menu-item index="open" @click="localGroup.handleOpenFolder(item.Path)">
                                <el-icon>
                                    <FolderOpened />
                                </el-icon>
                                {{ $t("navigator.open_folder") }}</el-menu-item>
                            <el-menu-item index="edit" @click="localGroup.editItem(groups.Name, {
                                Name: item.Name,
                                Type: item.Type,
                                Path: item.Path,
                            })">
                                <el-icon>
                                    <Edit />
                                </el-icon>
                                {{ $t("navigator.edit") }}</el-menu-item>
                            <el-menu-item index="delete" @click="localGroup.deleteItem(groups.Name, {
                                Name: item.Name,
                                Type: item.Type,
                                Path: item.Path,
                            })">
                                <el-icon>
                                    <DeleteFilled />
                                </el-icon>
                                {{ $t("navigator.delete") }}</el-menu-item>
                        </el-menu>
                        <template #reference>
                            <el-button bg text @click="RunApp(item.Type, item.Path)">
                                <template #icon>
                                    <img :src="localGroup.getSvgPath(item.Type)" style="width: 16px;">
                                </template>
                                {{ item.Name }}
                            </el-button>
                        </template>
                    </el-popover>
                </div>
            </div>
        </el-card>
    </div>
    <el-dialog v-model="global.temp.localAddItem" :title="$t('navigator.add_item')" width="500">
        <el-form label-width="auto">
            <el-form-item label="Group">
                <el-select v-model="config.defualtGroupName">
                    <el-option v-for="name in localGroup.getGroupNames()" :value="name">{{ name }}</el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="Name">
                <el-input v-model="config.name" />
            </el-form-item>
            <el-form-item label="Type">
                <el-radio-group v-model="config.defaultType">
                    <el-radio border v-for="item in localGroup.openGroup" :key="item" :value="item">{{ item
                        }}</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item label="Path">
                <el-input v-model="config.path">
                    <template #suffix>
                        <el-button link @click="localGroup.selectFile()">
                            <template #icon>
                                <Document />
                            </template>
                        </el-button>
                    </template>
                </el-input>
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="dialog-footer">
                <el-button type="primary" @click="localGroup.addItem">
                    OK
                </el-button>
            </div>
        </template>
    </el-dialog>
    <el-dialog v-model="config.editDialog" :title="$t('navigator.edit_item')" width="500">
        <el-form label-width="auto">
            <el-form-item label="Name">
                <el-input v-model="config.editName" />
            </el-form-item>
            <el-form-item label="Type">
                <el-radio-group v-model="config.editType">
                    <el-radio border v-for="item in localGroup.openGroup" :key="item" :value="item">{{ item
                        }}</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item label="Path">
                <el-input v-model="config.editPath"></el-input>
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="dialog-footer">
                <el-button type="primary" @click="localGroup.saveEdit">
                    Save
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>


<style scoped>
.button-grid {
    display: flex;
    flex-wrap: wrap;
    margin-left: -10px;
    /* 调整左边距以抵消子元素的左边距 */
}

.button-grid .el-button {
    margin-left: 10px;
    /* 设置按钮的左边距 */
    margin-bottom: 10px;
    /* 设置按钮的下边距 */
    flex: 0 1 auto;
    /* 按钮根据内容自动调整宽度 */
}

.drop-enable {
    --wails-drop-target: drop;
    width: 100%;
}
</style>