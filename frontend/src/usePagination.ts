import { shallowReactive, watch } from 'vue';

interface PaginationState<T> {
    currentPage: number;
    pageSize: number;
    result: T[];
    pageContent: T[];
    selectRows: T[];
}

interface PaginationController {
    handleSizeChange: (val: number) => void;
    handleCurrentChange: (val: number) => void;
    handleSelectChange: (rows:any[]) => void;
    getColumnData: (prop: string) => any[];
    watchResultChange: (result: any[], currentPage: number, pageSize: number) => any[];
    initTable: () => void; // 初始化表格数据
}

function usePagination<T>(data: T[], initialPageSize: number): { table: PaginationState<T>, ctrl: PaginationController } {
    const table = shallowReactive<PaginationState<T>>({
        currentPage: 1,
        pageSize: initialPageSize,
        result: data,
        pageContent: data.slice(0, initialPageSize),
        selectRows: [],
    });
    watch(() => table.result, (newResult: T[]) => {
        const start = (table.currentPage - 1) * table.pageSize;
        const end = table.currentPage * table.pageSize;
        table.pageContent = newResult.slice(start, end);
    });
    const ctrl: PaginationController = {
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
        watchResultChange: (result: any[], currentPage: number, pageSize: number) => {
            return result.slice((currentPage - 1) * pageSize, (currentPage - 1) * pageSize + pageSize)
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