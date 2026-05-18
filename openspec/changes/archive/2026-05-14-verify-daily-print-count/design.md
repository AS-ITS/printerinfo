## Context

目前的 `daily_audit.html` 只檢查 `daily_stats` 視圖的 `total_count` 間隔異常、counter 重置與日均 > 500 三種情況。`daily_stats` 是 Supabase 上的 PostgreSQL 視圖，透過 `printer_metrics` 表中相鄰兩筆記錄的差值計算每日增量。`collect_metrics.go` 負責透過 SNMP 採集印表機計數器並寫入 `printer_metrics` 表。

## Goals / Non-Goals

**Goals:**
- 在 `daily_audit.html` 新增多項檢查以驗證增量計算正確性
- 提供 Ground Truth 比對模式：直接從 `printer_metrics` 重新計算並與 `daily_stats` 比對
- 日均異常門檻可設定

**Non-Goals:**
- 修正 `daily_stats` 視圖的 SQL 邏輯（此 change 不包含資料庫端修改）
- 新增後端 API（所有檢查由前端 client-side 完成）
- 自動修復異常數據

## Decisions

1. **所有檢查在前端完成**：`daily_audit.html` 已可直接讀取 `printer_metrics` 和 `daily_stats`，無需新增後端 API。
2. **Ground Truth 使用 client-side 重算**：前端拉取所有 `printer_metrics` 資料後，在 JavaScript 中重算 delta 並與 `daily_stats` 比對。
3. **日均門檻改為 URL query param**：維持簡單設計，用 `?daily_avg_threshold=XXX` 參數控制，預設 500。
4. **CSV 匯出使用 Blob URL**：前端原生 `Blob` + `URL.createObjectURL` 產生下載，無需後端支援。

## Risks / Trade-offs

- [大量資料客戶端性能] → 若 `printer_metrics` 資料超過數萬筆，前端 JavaScript 處理可能變慢。Mitigation：加入分頁或限制最近 90 天的預設範圍。
- [Ground Truth 與 daily_stats 差異可能源自兩者計算基準不同] → 需清楚標示差異原因（基準筆選擇邏輯不同等）。
- [supabase.sql 的 daily_stats 視圖有複雜的 prev 選擇邏輯] → 若需修正 SQL，需另開 change。
