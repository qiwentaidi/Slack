import { computed } from 'vue'
import global from '@/global';

export const titleStyle = computed(() => {
  return global.Theme.value ? {
      backgroundColor: '#333333',
  } : {
      backgroundColor: '#eee',
  };
})

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
    return global.temp.isMacOS ? { marginLeft: '5.5%' } : {};
})

export const defaultIconSize = {
    width: '16px',
    height: '16px',
}