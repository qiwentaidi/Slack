import { shallowReactive, watch } from 'vue';

interface PaginationState<T> {
    currentPage: number;
    pageSize: number;
    result: T[];
    pageContent: T[];
    selectRows: T[];
    temp: T[];
}

interface PaginationController<T> {
    handleSizeChange: (val: number) => void;
    handleCurrentChange: (val: number) => void;
    handleSelectChange: (rows:any[]) => void;
    getColumnData: (prop: string) => any[];
    watchResultChange: (table: PaginationState<T>) => T[];
    initTable: () => void; // 初始化表格数据
}

function usePagination<T>(initialPageSize: number): { table: PaginationState<T>, ctrl: PaginationController<T> } {
    const table = shallowReactive<PaginationState<T>>({
        currentPage: 1,
        pageSize: initialPageSize,
        result: [] as T[],
        pageContent: [] as T[],
        selectRows: [],
        temp: [],
    });
    const ctrl: PaginationController<T> = {
        handleSizeChange: (val: number) => {
            table.pageSize = val;
            table.currentPage = 1;
            table.pageContent = table.result.slice(0, val);
        },
        handleCurrentChange: (val: number) => {
            table.currentPage = val;
            table.pageContent = table.result.slice((val - 1) * table.pageSize, val * table.pageSize);
        },
        handleSelectChange: (rows) => {
            table.selectRows = rows
        },
        getColumnData: (prop: string) => {
            return table.result.map((item: any) => item[prop]);
        },
        watchResultChange: (table: PaginationState<T>) => {
            const start = (table.currentPage - 1) * table.pageSize;
            const end = table.currentPage * table.pageSize;
            return table.result.slice(start, end);
        },
        initTable: () => {
            table.result = []
            table.pageContent = []
            table.selectRows = []
        }
    };
    return { table, ctrl };
}

export default usePagination;