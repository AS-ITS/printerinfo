# Proposal: Homepage Date Range Selector

## Why

目前首頁 `index.html` 的日期面板完全依賴自動偵測資料範圍，使用者無法指定查詢區間，導致：
1. 無法針對特定期間（如上月、特定月份）檢視歷史資料
2. 概況條的「今日」統計與趨勢圖無法對應到使用者想看的期間
3. 趨勢圖（每日總印量）與統計圖僅展示所有歷史資料，缺乏時間篩選能力

## What Changes

- **日期選擇器 UI**：在 `hero-side` 的 `date-panel` 內加入兩個 `<input type="date">` 欄位（開始日期、結束日期），取代目前的只讀 `date-field`
- **預設日期**：自動填入今年的 1/1 到今天（如 2026/01/01 ~ 2026/05/18）
- **重設按鈕**：加入「重設」按鈕回到預設日期範圍
- **資料篩選**：選擇區間後，所有統計（概況條、趨勢圖、年度圖、單位圖、百分比圖）均基於該日期區間重算
- **今日統計模式**：概況條標籤從「今日」改為「區間」，圖表標題同步更新
- **URL state**：日期區間寫入 URL query params（`?dateStart=&dateEnd=`），方便分享與重整

## Capabilities

### New Capabilities
- `homepage-date-range`: 首頁日期區間選擇器，包含日期輸入、重設、區間篩選資料、預設區間

### Modified Capabilities
<!-- No existing spec-driven capabilities to modify -->

## Impact

- **受影響檔案**：`index.html`（單一 HTML 檔案）
- **受影響功能**：
  - `updateDataRange()` — 改為讀取 input 值而非自動推算
  - `computeSummary()` — 改為依區間計算（非僅今日）
  - `buildPrinterData()` — 今日印量改為區間印量
  - `renderSummary()` — 標籤從「今日」改為區間名稱
  - `renderTrendChart()` — 圖表標題與 x 軸標籤依區間更新
  - `renderYearlyChart()` — 資料過濾依區間
  - 資料載入流程：`fetchAllSupabaseRows(daily_stats)` 改為加 `.gte()` / `.lte()` 條件
  - 單位篩選按鈕的資料點標誌改為依區間顯示
