declare module "echarts" {
  export type ECharts = any;
  export function init(dom: HTMLElement): ECharts;
  export const graphic: any;
  const echarts: {
    init: typeof init;
    graphic: typeof graphic;
  };
  export default echarts;
}
