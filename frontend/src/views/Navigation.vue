<script lang="ts" setup>
var options = [
    {
        label: "aside.penetration",
        value: [
            {
                name: "aside.webscan",
                route: "/Permeation/Webscan",
                icon: "/webscan.png"
            },
            {
                name: "aside.portscan",
                route: "/Permeation/Portscan",
                icon: "/portscan.png"
            },
            {
                name: "aside.dirscan",
                route: "/Permeation/Dirsearch",
                icon: "/dirscan.png"
            },
            {
                name: "aside.jsfinder",
                route: "/Permeation/Jsfinder",
                icon: "/extract.png"
            },
            {
                name: "aside.exp",
                route: "/Permeation/Exploitation",
                icon: "/attack.png"
            },
        ],
    },
    {
        label: 'aside.asset_collection',
        value: [
            {
                name: "aside.asset_from_company",
                route: "/Asset/Company",
                icon: "/company.png",
            },
            {
                name: "aside.subdomain_brute_force",
                route: "/Asset/Subdomain",
                icon: "/bomb.png"
            },
            {
                name: "aside.search_domain_info",
                route: "/Asset/Ipdomain",
                icon: "/domain.png"
            },
        ],
    },
    {
        label: "aside.space_engine",
        value: [
            {
                name: "aside.fofa",
                route: "/SpaceEngine/FOFA",
                icon: "/fofa.ico"
            },
            {
                name: "aside.hunter",
                route: "/SpaceEngine/Hunter",
                icon: "/hunter.ico"
            },
            {
                name: "aside.agent_pool",
                route: "/SpaceEngine/AgentPool",
                icon: "/vpn.png"
            },
        ],
    },
    {
        label: "aside.tools",
        value: [
            {
                name: "aside.en_and_de",
                route: "/Tools/Codec",
                icon: "/encrypt.png"
            },
            {
                name: "aside.systeminfo",
                route: "/Tools/System",
                icon: "/virus.png"
            },
            {
                name:"aside.data_handing",
                route: "/Tools/DataHanding",
                icon: "/data.png"
            },
            {
                name: "aside.memorandum",
                route: "/Tools/Memo",
                icon: "/memo.png"
            },
            {
                name: "aside.associate_dictionary_generator",
                route: "/Tools/Thinkdict",
                icon: "/dict.png"
            },
            {
                name: "aside.aksk",
                route: "/Tools/AKSK",
                icon: "/wechat.png"
            },
        ],
    },
]

</script>

<template>
    <el-scrollbar height="85vh">
        <div v-for="groups in options" style="margin-top: 10px; width: 99%;">
            <el-card style="width: 100%">
                <div style="margin-bottom: 10px;"><span style="font-weight: bold;">{{ $t(groups.label) }}</span></div>
                <el-row>
                    <el-col :lg="4" v-for="item in groups.value">
                        <el-result :title="$t(item.name)" @click="$router.push(item.route)">
                            <template #icon>
                                <el-image class="appNavIcon" :src="item.icon" />
                            </template>
                        </el-result>
                    </el-col>
                </el-row>
            </el-card>
        </div>
    </el-scrollbar>
</template>

<style>
.el-result:hover {
    cursor: pointer;
    background-color: #f5f7fa;
}

.appNavIcon .el-image__inner {
    height: 30px;
    width: 30px;
}
</style>