import { shallowReactive } from 'vue';

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
    watchResultChange: (table: PaginationState<T>) => T[];
}

function usePagination<T>(initialPageSize: number): { table: PaginationState<T>, ctrl: PaginationController<T>, initTable: () => void } {
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
        watchResultChange: (table: PaginationState<T>) => {
            const start = (table.currentPage - 1) * table.pageSize;
            const end = table.currentPage * table.pageSize;
            return table.result.slice(start, end);
        },
    };
    const initTable = () => {
        table.result = []
        table.temp = []
        table.selectRows = []
        table.pageContent = ctrl.watchResultChange(table)
    };
    return { table, ctrl, initTable };
}

export default usePagination;