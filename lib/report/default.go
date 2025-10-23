package report

import "strings"

func xssfilter(s string) string {
	s = strings.ReplaceAll(s, "<", "%3C")
	s = strings.ReplaceAll(s, ">", "%3E")
	return s
}

func defaultHeader() string {
	return `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width,initial-scale=1" />
<title>Vulnerability Report</title>
<style>
  :root{
    --bg:#f7f9fc;
    --panel:#ffffff;
    --panel-alt:#f1f3f5;
    --text:#1a1a1a;
    --muted:#6b7280;
    --line:#d1d5db;
    --ok:#16a34a;
    --warn:#f59e0b;
    --bad:#ef4444;
    --info:#3b82f6;
  }
  html,body{margin:0;padding:0;background:var(--bg);color:var(--text);font-family:ui-sans-serif,system-ui,-apple-system,Segoe UI,Roboto,"Helvetica Neue",Arial,"Noto Sans",sans-serif;font-size:14px;line-height:1.6;}
  a{color:var(--info);text-decoration:none;}
  a:hover{text-decoration:underline;}
  .container{max-width:1200px;margin:24px auto;padding:0 16px;}
  .hdr{display:flex;gap:12px;justify-content:space-between;align-items:center;margin-bottom:16px;}
  .hdr h1{margin:0;font-size:24px;}
  .toolbar{display:flex;gap:8px;flex-wrap:wrap;}
  .btn{background:#ffffff;border:1px solid var(--line);color:var(--text);padding:6px 12px;border-radius:8px;cursor:pointer;transition:0.2s;}
  .btn:hover{background:#f3f4f6;}
  .tabs{display:flex;border-bottom:1px solid var(--line);margin-bottom:16px;}
  .tab{padding:10px 16px;cursor:pointer;}
  .tab.active{font-weight:600;border-bottom:3px solid var(--info);color:var(--info);}
  .tab-content{display:none;}
  .tab-content.active{display:block;}
  .fp-row{padding:8px;border-bottom:1px solid var(--line);display:flex;gap:8px;align-items:center;flex-wrap:wrap;}
  .fp-id{width:40px;text-align:right;font-weight:600;}
  .badge{display:inline-block;padding:2px 6px;border-radius:8px;font-size:12px;}
  .badge.ok{color:var(--ok);border:1px solid var(--ok);}
  .badge.bad{color:var(--bad);border:1px solid var(--bad);}
  .badge.warn{color:var(--warn);border:1px solid var(--warn);}
  table{width:100%;border-collapse:collapse;margin-bottom:12px;background:var(--panel);border-radius:8px;overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,0.05);}
  .vuln-head{background:#e5e7eb;cursor:pointer;}
  .vuln-head td{padding:10px;font-weight:600;}
  tbody td{padding:10px;border-bottom:1px solid var(--line);}
  td.vuln{width:35%;}
  td.severity{width:15%;text-align:center;}
  td.url{width:50%;word-break:break-all;}
  .security { display:inline-block;padding:2px 8px;font-size:12px;font-weight:600;border-radius:6px;border:1px solid; }
  .security.low{ color:#3b82f6; border-color:#3b82f6; background:#eff6ff; }
  .security.medium{ color:#f59e0b; border-color:#f59e0b; background:#fffbeb; }
  .security.high{ color:#ef4444; border-color:#ef4444; background:#fef2f2; }
  .security.critical{ color:#b91c1c; border-color:#b91c1c; background:#fef2f2; }
  .clr{display:flex;gap:10px;position:relative;flex-wrap:wrap;}
  .request,.response{position:relative;border:1px solid var(--line);background:var(--panel-alt);border-radius:8px;overflow:hidden;flex:1 1 48%;}
  .w100{flex:1 1 100%;}
  .reqresp-header {position:absolute;top:0; left:0; right:0; display:flex;justify-content:space-between; align-items:center; padding:4px 8px; background:rgba(107,114,128,0.85); color:white; font-size:12px;}
  .reqresp-header .actions{ display:flex; gap:6px; }
  .reqresp-header .btn-mini{padding:2px 6px;border-radius:4px;background:#10b981;cursor:pointer;}
  .reqresp-header .time{ font-size:12px; opacity:0.9; }
  	xmp {
		white-space: pre-wrap; /* 保留空格，但允许换行 */
		word-break: break-all; /* 强制长单词换行，避免撑开容器 */
		margin: 0;
		padding: 32px 12px 12px 12px;
		font-family: ui-monospace, Consolas, monospace;
		font-size: 12px;
		line-height: 1.45;
		max-height: 400px; /* 最大高度，根据需求调整 */
		overflow: auto; /* 超出内容显示滚动条 */
		display: block; /* 确保块级显示 */
	}
	.fp-tags {
		display: flex;
		flex-wrap: wrap;
		gap: 6px;
	}

	.fp-tag {
		display: inline-block;
		border: 1px solid #3b82f6; /* primary 蓝色边框 */
		color: #3b82f6;             /* 同边框色文字 */
		background: #f0f6ff;        /* 浅蓝背景，或者用 transparent */
		font-size: 12px;
		padding: 2px 6px;
		border-radius: 6px;
	}
</style>
</head>
<body>
<div class="container">
  <div class="hdr">
    <h1>Vulnerability Report</h1>
    <div class="toolbar">
      <button class="btn" onclick="expandAll()">展开全部</button>
      <button class="btn" onclick="collapseAll()">折叠全部</button>
      <button class="btn" onclick="window.print()">打印/导出 PDF</button>
    </div>
  </div>
  <div class="tabs">
    <div class="tab active" data-tab="fp">Fingerprints</div>
    <div class="tab" data-tab="vuln">Vulnerabilities</div>
  </div>
  <div id="fp" class="tab-content active">`
}

func defaultFooter() string {
	return `</div> <!-- vuln tab -->
<div id="vuln" class="tab-content"></div>
</div> <!-- container -->
<script src="https://cdn.jsdelivr.net/npm/jquery@3.7.1/dist/jquery.min.js"></script>
<script>
$(".tab").on("click", function(){
  $(".tab").removeClass("active");
  $(this).addClass("active");
  $(".tab-content").removeClass("active");
  $("#" + $(this).data("tab")).addClass("active");
});
$(document).on('click', '.vuln-head', function(){
  $(this).next('tbody').slideToggle(150);
});
function copyXmpContent($btn){
  try{
    var $box = $btn.closest('.request, .response');
    var txt = $box.find('xmp').text();
    navigator.clipboard.writeText(txt).then(function(){ showToast('Copied'); }).catch(function(){ alert('复制失败'); });
  }catch(e){ alert('复制失败：'+e); }
}
function showToast(msg){
  var n = document.createElement('div');
  n.textContent = msg;
  n.style.cssText = "position:fixed;right:16px;bottom:16px;background:#10b981;color:white;padding:8px 12px;border-radius:8px;z-index:9999";
  document.body.appendChild(n);
  setTimeout(()=>{ n.remove(); },1500);
}
function expandAll(){ $('table tbody').show(); }
function collapseAll(){ $('table tbody').hide(); }
</script>
</body>
</html>`
}
