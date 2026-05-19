## Context

首頁的「每日總印量趨勢」圖表依賴以下資料鏈：
1. **`printer_metrics`** 表：Go 收集器透過 SNMP 定期寫入累積計數
2. **`daily_stats` SQL VIEW**：計算每日 delta（`daily_print`, `daily_copy`, `daily_fax`, `daily_total`）
3. **前端 `renderTrendChart`**：依選定日期區間聚合 `daily_total` 並渲染趨勢圖

使用者回報 4/15 之後的資料沒有顯示。問題可能出在：
- `printer_metrics` 在 4/15 之後無資料（資料收集中斷）
- `daily_stats` VIEW 的 delta 計算在邊界情況下產生 0
- 前端 `buildCorrectedTrend` 的 delta 計算邏輯與 VIEW 不一致

## Goals / Non-Goals

**Goals:**
- 確認 4/15 之後的資料是否存在
- 修復 `daily_stats` VIEW 或 `buildCorrectedTrend` 的計算問題
- 確保趨勢圖正確顯示所有日期區間的資料

**Non-Goals:**
- 修復資料收集流程（這是 deploy.yml 的責任）
- 新增新的資料收集機制

## Decisions

### 1. 優先調查 `daily_stats` SQL VIEW

`daily_stats` VIEW 使用 `LEFT JOIN LATERAL` 查找上一筆非零讀數來計算 delta。需要確認：
- `current_has_reading.has_value` 條件是否正確判斷非零讀數
- `prev` 的 WHERE 子句是否正確過濾出上一筆有效讀數
- `daily_total` 的 fallback 邏輯（`total_delta` vs `daily_print + daily_copy + daily_scan + daily_fax`）是否正確

### 2. 前端 `buildCorrectedTrend` 與 VIEW 保持同步

前端使用相同的 delta 計算邏輯（`counterDelta`）作為安全網。如果 VIEW 修正，前端也需同步更新。

### 3. 不修改 `printer_metrics` 的資料收集邏輯

資料收集（`collect_metrics.go`）不在此變更範圍內。如果 `printer_metrics` 確實無 4/15 之後的資料，需回報給 deploy 流程維護者。

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| 修改 VIEW 可能影響其他依賴 `daily_stats` 的功能 | 新增測試場景驗證所有欄位計算正確 |
| 前端與 VIEW 計算邏輯可能再次分歧 | 確保 `buildCorrectedTrend` 與 VIEW 邏輯一致 |

## Migration Plan

1. 驗證 `printer_metrics` 資料是否存在
2. 修正 `daily_stats` VIEW 的計算邏輯
3. 同步更新 `buildCorrectedTrend`
4. 部署後手動驗證趨勢圖顯示

## Open Questions

- `printer_metrics` 表在 4/15 之後是否有資料？
- `daily_stats` VIEW 是否影響其他儀表板功能？
