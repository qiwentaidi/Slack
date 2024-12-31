import { shallowReactive } from 'vue';
import { sortSeverityOptions } from './stores/options';

interface PaginationState<T> {
    currentPage: number;
    pageSize: number;
    result: T[];
    pageContent: T[];
    selectRows: T[];
    sortTemp: T[]; // 排序时的临时数据
    filterTemp: T[]; // 过滤时的临时数据
    isSorted: boolean;
}

interface PaginationController<T> {
    handleSizeChange: (val: number) => void;
    handleCurrentChange: (val: number) => void;
    handleSelectChange: (rows: any[]) => void;
    watchResultChange: (table: PaginationState<T>) => void;
    sortChange: (data: { column: any, prop: string, order: any }, isvultable: boolean) => void;
    getColumnFilters: (prop: string) => Array<{ text: string, value: string }>;
    filterChange: (newFilters: any) => void; // 由element-plus的filter-change事件触发
    inputFilter: (prop: string, value: string) => void, // 与filterChange事件不能同时使用在一个表格中
}

function usePagination<T>(initialPageSize: number): { table: PaginationState<T>, ctrl: PaginationController<T>, initTable: () => void } {
    const table = shallowReactive<PaginationState<T>>({
        currentPage: 1,
        pageSize: initialPageSize,
        result: [] as T[],
        pageContent: [] as T[],
        selectRows: [],
        sortTemp: [],
        filterTemp: [],
        isSorted: false,
    });
    const ctrl: PaginationController<T> = {
        handleSizeChange: (val: number) => {
            table.pageSize = val;
            table.currentPage = 1;
            if (table.isSorted) {
                table.pageContent = table.sortTemp.slice(0, val);
            } else {
                table.pageContent = table.result.slice(0, val);
            }
        },
        handleCurrentChange: (val: number) => {
            table.currentPage = val;
            if (table.isSorted) {
                table.pageContent = table.sortTemp.slice((val - 1) * table.pageSize, val * table.pageSize);
            } else {
                table.pageContent = table.result.slice((val - 1) * table.pageSize, val * table.pageSize);
            }
        },
        handleSelectChange: (rows) => {
            table.selectRows = rows
        },
        watchResultChange: (table: PaginationState<T>) => {
            const start = (table.currentPage - 1) * table.pageSize;
            const end = table.currentPage * table.pageSize;
            table.pageContent = table.result.slice(start, end);
        },
        // 通过判断缺省参数是否为true，来判断是否为漏洞表格排序
        sortChange(data: { column: any, prop: string, order: any }, isvultable: boolean = false): void {
            if (!data.prop || !data.order) {
                table.isSorted = false
                // 使用原数据进行输出，即可恢复数据状态
                ctrl.watchResultChange(table)
            } else {
                table.isSorted = true
                // 排序的数据都由temp提供
                table.sortTemp = [...table.result];

                // 根据排序规则对结果进行排序
                table.sortTemp.sort((a, b) => {
                    var valA: any
                    var valB: any
                    if (isvultable) {
                        valA = sortSeverityOptions.indexOf(a["Severity"]);
                        valB = sortSeverityOptions.indexOf(b["Severity"]);
                    } else {
                        valA = a[data.prop];
                        valB = b[data.prop];
                    }
                    // 比较逻辑，支持字符串、数字、日期等
                    if (valA < valB) {
                        return data.order === 'ascending' ? -1 : 1;
                    } else if (valA > valB) {
                        return data.order === 'ascending' ? 1 : -1;
                    } else {
                        return 0;
                    }
                });
                const start = (table.currentPage - 1) * table.pageSize;
                const end = table.currentPage * table.pageSize;
                table.pageContent = table.sortTemp.slice(start, end);
            }
        },
        getColumnFilters(prop: string): Array<{ text: string, value: string }> {
            let values = table.result.map(item => item[prop])
            let uniqueValues = Array.from(new Set(values));
            return uniqueValues.map(value => ({ text: value, value }));
        },
        // column-key 一定要和prop的值一致
        filterChange (newFilters: any) {
            // 获取 newFilters 中的第一个，也是唯一的属性名
            const filterKey = Object.keys(newFilters)[0];
            const selectedFilters = newFilters[filterKey] ? [...newFilters[filterKey]] : [];
            // 头一次进行筛选数据
            if (table.filterTemp.length == 0) {
                table.filterTemp = [...table.result];
            }
            if (selectedFilters.length == 0) {
                table.result = [...table.filterTemp];
            } else {
                table.result = table.filterTemp.filter(item => selectedFilters.includes(item[filterKey]));
            }
            ctrl.watchResultChange(table)
        },
        inputFilter(prop: string, value: string) {
            // 判断是不是首次进行数据筛选
            if (table.filterTemp.length == 0) {
                table.filterTemp = [...table.result];
            }
            if (value === '') {
                table.result = [...table.filterTemp];
            } else {
                table.result = table.filterTemp.filter(item => item[prop].toLowerCase().includes(value.toLowerCase()));
            }
            table.currentPage = 1
            ctrl.watchResultChange(table)
        }
    };
    const initTable = () => {
        table.result = []
        table.sortTemp = []
        table.filterTemp = []
        table.selectRows = []
        ctrl.watchResultChange(table)
    };
    return { table, ctrl, initTable };
}

export default usePagination;