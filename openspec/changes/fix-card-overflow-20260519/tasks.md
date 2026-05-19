## 1. 縮小 Summary 卡片數值文字

- [x] 1.1 將 `.summary-value` 的 `font-size` 從 34px 改為 28px
- [x] 1.2 將 `.summary-value` 的 `white-space: nowrap` 改為 `normal`（移除或保留，視需要）
- [x] 1.3 同步調整 `.summary-strip` 的 `grid-template-columns` 最小寬度，使總寬不超過 1560px 容器

## 2. 縮小日期選擇器卡片

- [x] 2.1 將 `.hero-side` 的 `.date-panel` 對應欄位從 `minmax(160px, 190px)` 改為 `minmax(140px, 160px)`
- [x] 2.2 將 `.date-field` 的 `font-size` 從 13px 改為 11px
- [x] 2.3 將 `.date-panel` 的 padding 從 `10px 12px` 改為 `8px 10px`

## 3. 縮小各單位印表機列印量卡片

- [x] 3.1 將 `.dashboard-main-grid` 的 `pcnt` 區域與 `yearly` 區域共享右側欄：改為 `grid-template-columns: minmax(0, 1fr) minmax(0, 1fr)` 或固定比例
- [x] 3.2 將 `.chart-container` 的預設高度從 252px 改為 200px（若需要更多空間，可保留 220px）

## 4. 驗證

- [x] 4.1 在 1920px 寬度下檢查所有卡片無溢出
- [x] 4.2 在 1366px 寬度下檢查所有卡片無溢出
- [x] 4.3 在 980px 響應式斷點下檢查佈局正確
