import { shallowReactive } from 'vue';

interface PaginationState<T> {
    currentPage: number;
    pageSize: number;
    result: T[];
    pageContent: T[];
    selectRows: T[];
    temp: T[];
    isSorted: boolean;
}

interface PaginationController<T> {
    handleSizeChange: (val: number) => void;
    handleCurrentChange: (val: number) => void;
    handleSelectChange: (rows: any[]) => void;
    watchResultChange: (table: PaginationState<T>) => void;
    sortChange: (data: { column: any, prop: string, order: any }) => void;
}

function usePagination<T>(initialPageSize: number): { table: PaginationState<T>, ctrl: PaginationController<T>, initTable: () => void } {
    const table = shallowReactive<PaginationState<T>>({
        currentPage: 1,
        pageSize: initialPageSize,
        result: [] as T[],
        pageContent: [] as T[],
        selectRows: [],
        temp: [],
        isSorted: false,
    });
    const ctrl: PaginationController<T> = {
        handleSizeChange: (val: number) => {
            table.pageSize = val;
            table.currentPage = 1;
            if (table.isSorted) {
                table.pageContent = table.temp.slice(0, val);
            } else {
                table.pageContent = table.result.slice(0, val);
            }
        },
        handleCurrentChange: (val: number) => {
            table.currentPage = val;
            if (table.isSorted) {
                table.pageContent = table.temp.slice((val - 1) * table.pageSize, val * table.pageSize);
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
        sortChange(data: { column: any, prop: string, order: any }): void {
            if (!data.prop || !data.order) {
                table.isSorted = false
                // 使用原数据进行输出，即可恢复数据状态
                ctrl.watchResultChange(table)
            } else {
                table.isSorted = true
                // 排序的数据都由temp提供
                table.temp = [...table.result];

                // 根据排序规则对结果进行排序
                table.temp.sort((a, b) => {
                    const valA = a[data.prop];
                    const valB = b[data.prop];

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
                table.pageContent = table.temp.slice(start, end);
            }
        }
    };
    const initTable = () => {
        table.result = []
        table.temp = []
        table.selectRows = []
        ctrl.watchResultChange(table)
    };
    return { table, ctrl, initTable };
}

export default usePagination;