## 1. HTML 文字替換

- [x] 1.1 修改 `<title>` 從「今日未列印印表機 \| No Print Today」改為「未收到資料列表 \| No Data Received」
- [x] 1.2 修改登入畫面副標題從「今日未列印印表機列表」改為「未收到資料列表」
- [x] 1.3 修改 Hero `<h1>` 從「今日未列印印表機」改為「未收到資料列表」
- [x] 1.4 修改 Hero `<p>` 副標從「total_count 與 printer_count 皆為 0 的印表機」改為「未收到有效資料的印表機」
- [x] 1.5 修改摘要面板文字從「今日未列印」改為「未收到資料」

## 2. 表格與欄位名稱替換

- [x] 2.1 修改表格 `<th>` 欄位標題從「最近統計時間」改為「上次抓取時間」
- [x] 2.2 修改 `renderTable` 中的 `recordedAt` 變數名為 `fetchedAt`
- [x] 2.3 修改 `renderTable` 中的 `recordedAt` 引用為 `fetchedAt`

## 3. JavaScript 變數與文字更新

- [x] 3.1 修改 `loadData` 中的 `last_recorded` 變數名為 `last_fetched`
- [x] 3.2 修改 `document.getElementById('last-update').innerText` 從「最後更新」改為「上次抓取」
- [x] 3.3 修改 `setRefreshLoading` 按鈕文字從「更新中」改為「抓取中」

## 4. 驗證

- [x] 4.1 檢查所有「未列印」文字是否已替換
- [x] 4.2 檢查所有「統計時間」文字是否已替換
- [x] 4.3 確認頁面功能不受影響
