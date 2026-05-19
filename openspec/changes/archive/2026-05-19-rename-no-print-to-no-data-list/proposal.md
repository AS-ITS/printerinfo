# Proposal: Rename "今日未列印印表機" to "未收到資料列表"

## Why

當前 `no_print_today.html` 的命名與欄位名稱存在語意不清的問題：
1. 頁面標題「今日未列印印表機」容易誤導使用者，該頁面實際顯示的是「今日 total_count 與 printer_count 皆為 0」的印表機，包含根本未收到任何資料的印表機
2. 欄位「最近統計時間」名不副實，它顯示的是印表機最後有資料的時間，而非統計運算時間
3. 列表標題「今日未列印」同樣不夠精確

改名後能更準確反映頁面的真實用途：列出所有**未收到有效資料**的印表機，而非僅「未列印」的印表機。

## What Changes

- **頁面標題**：`no_print_today.html` 的 `<title>` 從「今日未列印印表機 | No Print Today」改為「未收到資料列表 | No Data Received」
- **登入頁面副標題**：從「今日未列印印表機列表」改為「未收到資料列表」
- **Hero 區塊標題**：從「今日未列印印表機」改為「未收到資料列表」
- **Hero 區塊副標說明**：從「total_count 與 printer_count 皆為 0 的印表機」改為「未收到有效資料的印表機」
- **篩選邏輯**：從「total_count = 0 且 printer_count = 0」改為「今日無任何紀錄 OR total_count 與 printer_count 皆為 0」（保持行為不變，僅改變名稱）
- **摘要欄位**：「今日未列印」改為「未收到資料」
- **表格欄位標題**：「最近統計時間」改為「上次抓取時間」
- **最後更新文字**：從「最後更新」改為「上次抓取」
- **變數名稱**：`last_recorded` 改為 `last_fetched`
- **JS 變數**：`setRefreshLoading` 相關文字從「更新中」改為「抓取中」

## Capabilities

### New Capabilities
<!-- No new capabilities needed. This is a naming/semantic change to existing page. -->

### Modified Capabilities
- `no-print-today-page`: 頁面標題、列表標題、表格欄位名稱從「未列印/統計」語意改為「未收到資料/抓取」語意

## Impact

- **受影響檔案**：`no_print_today.html`（單一 HTML 檔案）
- **受影響元素**：
  - HTML `<title>` 標籤
  - 登入畫面副標題 (`<p>`)
  - Hero 區塊標題 (`<h1>`) 與副標 (`<p>`)
  - 摘要面板文字（`renderSummary` 函式）
  - 表格 `<th>` 欄位標題
  - `renderTable` 中的 `recordedAt` 變數與顯示
  - `last-update` 元素顯示文字
  - `loadData` 中的 `last_recorded` 變數名稱
  - 頁面註解（篩選條件說明）
