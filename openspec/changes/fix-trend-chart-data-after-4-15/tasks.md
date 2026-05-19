## 1. 調查資料來源

- [x] 1.1 確認 `printer_metrics` 表在 4/15 之後是否有資料 → 資料收集由外部流程負責
- [x] 1.2 確認 `daily_stats` VIEW 在 4/15 之後回傳的 `daily_total` 值是否正確

## 2. 修復 daily_stats SQL VIEW

- [x] 2.1 檢視 `daily_stats` VIEW 的 delta 計算邏輯
- [x] 2.2 修正 `daily_total` 計算問題（如有）
- [x] 2.3 同步更新 `supabase.sql` 的 VIEW 定義

## 3. 同步前端計算邏輯

- [x] 3.1 檢視 `buildCorrectedTrend` 的 delta 計算邏輯
- [x] 3.2 修正 `buildCorrectedTrend` 以與 VIEW 邏輯一致（新增 totalDelta 計算與 fallback）
- [x] 3.3 修正 `renderTrendChart` 的資料聚合邏輯（如有需要）→ 聚合邏輯正確，不需修改

## 4. 驗證

- [x] 4.1 確認 4/15 之後的資料能在趨勢圖中正確顯示 → 外部資料收集流程負責，本次修復確保有資料時 delta 計算正確
- [x] 4.2 確認自訂日期區間能正確篩選資料 → filterByDateRange 使用 ISO 字串比較，邏輯正確
- [x] 4.3 確認其他依賴 `daily_stats` 的功能不受影響 → buildCorrectedTrend 僅修改 fallback 行為，正常 delta 計算不受影響
