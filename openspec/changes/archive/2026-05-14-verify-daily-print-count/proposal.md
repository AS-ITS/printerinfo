## Why

目前的 `daily_audit.html` 只檢查 `daily_stats` 視圖的 `total_count` 間隔異常、counter 重置與日均 > 500 三種情況，無法驗證 daily_stats 的增量計算是否正確（例如 print/copy/scan/fax 各別欄位的 delta 是否有誤、daily_total 是否等於各別欄位總和、是否有負值等）。需要增強驗證能力來確認每日列印量計算邏輯的正確性。

## What Changes

- 增強 `daily_audit.html` 的檢查項目：
  - 各別欄位（print, copy, scan, fax）的 delta 為負值
  - `daily_total` 不等於各別欄位 delta 總和
  - 各別欄位的日均異常（> 500）
  - 同一日多筆 `printer_metrics` 記錄的重複問題
- 新增「Ground Truth」比對模式：直接從 `printer_metrics` 重新計算增量並與 `daily_stats` 視圖比對
- 新增匯出功能：將異常清單匯出為 CSV
- 將日均異常門檻從固定 500 改為可設定（預設維持 500）

## Capabilities

### New Capabilities
- `daily-print-audit`: 每日列印量增量計算正確性檢查的規格與驗證邏輯

### Modified Capabilities
- `dashboard-data-refresh`: 修改 daily_stats 視圖的增量計算需求（若需修正 supabase.sql 的 daily_stats 定義）

## Impact

- `daily_audit.html` — 前端檢查邏輯與 UI
- `supabase.sql` — `daily_stats` 視圖定義（若需修正）
- `scripts/collect_metrics.go` — 若需新增重複記錄檢查
