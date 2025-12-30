import { computed, markRaw } from 'vue'
import global from '@/stores';
import { webReportOptions } from './options';
import { ActivityItem } from './interface';
import { CircleCheck, Notification,  Warning, CloseBold, MoreFilled } from '@element-plus/icons-vue';

export const titlebarStyle = computed(() => {
    return global.Theme.value ? {
        backgroundColor: '#333333',
        borderBottom: '1px solid #3B3B3B'
    } : {
        backgroundColor: '#F9F9F9',
        borderBottom: '1px solid #E6E6E6'
    };
})

export const defaultIconSize = {
    width: '16px',
    height: '16px',
}

export const appStartStyle = computed(() => {
    return global.temp.isGrid ? {
        display: 'grid',
        gridTemplateColumns: 'repeat(10, 1fr)',
        gap: '10px'
    } : {
        display: 'flex',
        FlexWrap: 'wrap',
        gap: '10px'
    };
})

// v1.7.3 更新漏洞风险等级显示样式
export function getTagTypeBySeverity(severity: string) {
    switch (severity) {
        // v1.7.4 CRITICAL 类型由动态类实现，不然会出现控制台类型报错
        // case 'CRITICAL':
        //     return 'success';
        case 'HIGH':
            return 'danger';
        case 'MEDIUM':
            return 'warning';
        case 'LOW':
            return 'primary';
        default:
            return 'info';
    }
}

export function getBadgeClass(level: string) {
    switch (level) {
        case 'CRITICAL':
            return 'risk-badge-critical'
        case 'HIGH':
            return 'risk-badge-red'
        case 'MEDIUM':
            return 'risk-badge-yellow'
        case 'LOW':
            return 'risk-badge-green'
        case 'INFO':
        default:
            return 'risk-badge-gray'
    }
}

export function highlightFingerprints(fingerprint: string) {
    if (fingerprint == "疑似蜜罐") {
        return "warning"
    }
    if (global.webscan.highlight_fingerprints.includes(fingerprint)) {
        return "danger"
    }
    return "primary"
}

// 获取选中的 `icon`
export function getSelectedIcon(selectedLabel: string) {
    const selectedItem = webReportOptions.find(item => item.label === selectedLabel);
    return selectedItem ? selectedItem.icon : null;
};

// 手动指定5种类型对应的图标（都 markRaw 了）
export const typeToIcon: Record<ActivityItem["type"], any> = {
    primary: markRaw(MoreFilled),
    success: markRaw(CircleCheck),
    warning: markRaw(Warning),
    danger:  markRaw(CloseBold),
    info:    markRaw(Notification),
}