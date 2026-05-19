# Design: Homepage Date Range Selector

## Context

`index.html` 的日期面板目前為自動偵測模式：`updateDataRange()` 從 `trend` 資料中找出最早與最晚的 `recorded_at`，以只讀方式顯示在 `<span>` 中。使用者無法自訂查詢期間。

所有資料（概況條、趨勢圖、單位圖、年度圖等）皆基於「今日」或全部歷史資料。

## Goals / Non-Goals

**Goals:**
- 在 date-panel 內加入兩個 `<input type="date">` 供使用者選擇區間
- 預設區間為今年 1/1 ~ 今天
- 區間選擇後觸發重新篩選與渲染（不重新 API 呼叫，使用已快取的 daily_stats 資料）
- 概況條、趨勢圖、單位圖、年度圖等依區間更新

**Non-Goals:**
- 不改變 daily_stats 的 API 呼叫邏輯（仍抓全部資料到記憶體快取）
- 不新增後端 API 端點
- 不改變單位篩選的運作方式
- 不支援月曆 UI（僅 date input）

## Decisions

1. **前端快取 + 客戶端篩選**：維持現有 `fetchAllSupabaseRows(daily_stats)` 抓取全部資料到 `cachedMetrics`，在客戶端依區間篩選。理由：daily_stats 資料量不大，且此模式已是現有架構。

2. **原生 `<input type="date">`**：使用瀏覽器原生 date input，不引入外部 date picker 套件。理由：零依賴、RWD 無虞、瀏覽器已優化。

3. **區間變更觸發 re-render**：`input` 事件監聽 → `applyDateRange()` → 篩選 `cachedMetrics` → `renderScopedDashboard()`。理由：不重新 API 呼叫，速度極快。

4. **圖表更新策略**：區間變更後，趨勢圖顯示該區間內的每日資料，年度圖顯示區間內每年的彙總資料。

## Risks / Trade-offs

| Risk | Mitigation |
|------|------|
| 大量 daily_stats 資料在客戶端計算慢 | daily_stats 為彙總表，單筆查詢量小；若未來資料增長，可改為 API 側過濾 |
| 結束日期 < 開始日期 | `applyDateRange()` 內加驗證，錯誤時自動修正或顯示提示 |
| 區間內無資料 | 圖表顯示空狀態提示，不 crash |

## Migration Plan

單一 HTML 檔案修改，無需資料庫或後端變更。直接部署即可。

## Open Questions

無。
