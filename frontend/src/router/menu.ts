import homeIcon from "@/assets/icon/home.svg"
import toolsIcon from "@/assets/icon/tools.svg"
import mappingIcon from "@/assets/icon/mapping.svg"

const MenuList = [
    {
        name: "aside.home",
        path: "/",
        icon: homeIcon,
    },
    {
        name: "aside.penetration",
        path: "/Permeation", // 为 views 文件夹下的一级目录名称
        icon: "Smoking", // 如果为字符串，使用是 element-plus 的图标名称，由于注册了全图标会自动加载
        children: [
            {
                name: "aside.webscan",
                path: "/Webscan", // 二级文件 .vue 的名称
            },
            {
                name: "aside.dirscan",
                path: "/Dirsearch",
            },
            {
                name: "aside.jsfinder",
                path: "/JSFinder",
            },
            {
                name: "aside.exp",
                path: "/Exploitation",
            },
            {
                name: "aside.dumpall",
                path: "/Dumpall",
            },
            {
                name: "aside.request",
                path: "/Request",
            },
        ]
    },
    {
        name: "aside.asset_collection",
        path: "/Asset",
        icon: "OfficeBuilding",
        children: [
            {
                name: "aside.asset_from_company",
                path: "/Company",
            },
            {
                name: "aside.search_domain_info",
                path: "/Domain",
            },
            {
                name: "aside.isic",
                path: "/ISICollection",
            },
        ]
    },
    {
        name: "aside.space_engine",
        path: "/SpaceEngine",
        icon: mappingIcon,
        children: [
            {
                name: "aside.fofa",
                path: "/FOFA",
            },
            {
                name: "aside.hunter",
                path: "/Hunter",
            },
            {
                name: "aside.360quake",
                path: "/Quake",
            },
            {
                name: "aside.polymerization",
                path: "/Polymerization",
            },
            {
                name: "aside.agent_pool",
                path: "/AgentPool",
            },
        ]
    },
    {
        name: "aside.tools",
        path: "/Tools",
        icon: toolsIcon,
        children: [
            {
                name: "aside.cyberchef",
                path: "/CyberChef",
            },
            {
                name: "aside.systeminfo",
                path: "/System",
            },
            {
                name: "aside.data_handing",
                path: "/DataHanding",
            },
            {
                name: "aside.fscan",
                path: "/Fscan",
            },
            {
                name: "aside.aksk",
                path: "/AKSK",
            },
            {
                name: "aside.extract_database_info",
                path: "/ExtractDbInfo",
            },
            {
                name: "aside.fileinfo",
                path: "/FileContentRetrieval",
            },
            {
                name: "aside.data_comparison",
                path: "/DataComparison",
            },
            {
                name: "aside.timestamp",
                path: "/Timestamp",
            },
            {
                name: "aside.memorandum",
                path: "/Memo",
            },
        ]
    },
]

export default MenuList
