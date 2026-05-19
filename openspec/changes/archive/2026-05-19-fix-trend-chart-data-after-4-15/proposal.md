## Why

首頁的「每日總印量趨勢」圖表沒有顯示 4/15 之後的資料。使用者要求重新計算趨勢資料以恢復正常顯示。

## What Changes

- 調查並修復 `daily_stats` SQL VIEW 的 daily delta 計算邏輯
- 確認 `printer_metrics` 在 4/15 之後是否有資料
- 如有需要，修正 `buildCorrectedTrend` 前端 delta 計算以避免與 VIEW 結果不一致
- 確保趨勢圖正確顯示所有日期區間的資料

## Capabilities

### Modified Capabilities
- `dashboard-unit-navigation`: daily_stats 的 daily_total 計算邏輯有問題，導致 4/15 之後的資料無法正確顯示

## Impact

- `supabase.sql`: daily_stats VIEW 的 delta 計算邏輯
- `index.html`: `buildCorrectedTrend` 和 `renderTrendChart` 的資料聚合邏輯
- `scripts/collect_metrics.go`: 資料收集是否持續運作
