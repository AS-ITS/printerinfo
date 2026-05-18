## 1. 日期選擇器 UI

- [x] 1.1 修改 date-panel 內 HTML：將兩個只讀 `.date-field` span 改為兩個 `<input type="date">` input
- [x] 1.2 為 date-panel 加入「重設」按鈕
- [x] 1.3 在 init / loadData 時設定預設日期（今年 1/1 ~ 今天）
- [x] 1.4 加入 input event listener 觸發區間更新

## 2. 區間邏輯實作

- [x] 2.1 新增 `getSelectedDateRange()` 函式讀取 input 值
- [x] 2.2 新增 `applyDateRange()` 函式驗證區間並觸發 re-render
- [x] 2.3 新增 `resetDateRange()` 函式恢復預設區間
- [x] 2.4 區間驗證：結束日期 < 開始日期時自動修正

## 3. 資料篩選更新

- [x] 3.1 新增 `filterMetricsByDateRange(metrics, range)` 依區間過濾 metrics
- [x] 3.2 修改 `computeSummary()` 依區間計算統計（非僅今日）
- [x] 3.3 修改 `buildPrinterData()` 的「今日」印量改為區間印量
- [x] 3.4 修改 `getUnitsWithTodayPrintData()` 改為依區間顯示資料點

## 4. UI 渲染更新

- [x] 4.1 修改 `renderSummary()` 的標籤從「今日」改為區間名稱
- [x] 4.2 修改 `updateDataRange()` 改為顯示 input 值而非自動推算
- [x] 4.3 修改 `renderTrendChart()` 依區間過濾資料
- [x] 4.4 修改 `renderYearlyChart()` 依區間過濾資料
- [x] 4.5 修改 `renderUnitComparisonChart()` 依區間過濾資料
- [x] 4.6 修改 `renderUnitFilter()` 的資料點標誌依區間顯示

## 5. 驗證

- [x] 5.1 驗證預設區間正確顯示
- [x] 5.2 驗證自訂區間後所有圖表同步更新
- [x] 5.3 驗證重設按鈕恢復預設值
- [x] 5.4 驗證區間內無資料時的空狀態處理
