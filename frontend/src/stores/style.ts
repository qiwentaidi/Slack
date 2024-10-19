import { computed } from 'vue'
import global from '@/global';

export const titlebarStyle = computed(() => {
    return global.Theme.value ? {
        backgroundColor: '#333333',
        borderBottom: '1px solid #3B3B3B'
    } : {
        backgroundColor: '#F9F9F9',
        borderBottom: '1px solid #E6E6E6'
    };
})

export const rightStyle = computed(() => {
    return global.temp.isMacOS ? { marginRight: '3.5px' } : {};
})

export const leftStyle = computed(() => {
    return !global.temp.isMacOS ? { marginLeft: '3.5px' } : {};
})

export const macStyle = computed(() => {
    return global.temp.isMacOS ? { marginLeft: '80px' } : {};
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

export const eldividerStyle = computed(() => {
    if (global.Theme.value) {
        switch (global.temp.goos) {
            case "windows":
                return {
                    backgroundColor: '#121212',
                }
            case "darwin":
                return {
                    backgroundColor: '#1E1E1E',
                }
            case "linux":
                return {
                    backgroundColor: '#1D1E1F',
                }
        }
    }
});

export function getClassBySeverity(severity: string) {
    switch (severity) {
        case 'CRITICAL':
            return 'severity-critical';
        case 'HIGH':
            return 'severity-high';
        case 'MEDIUM':
            return 'severity-medium';
        case 'LOW':
            return 'severity-low';
        default:
            return 'severity-info';
    }
}
