## 1. 固定 Date Card 寬度

- [x] 1.1 將 `.hero-side` 的 date-panel 對應欄位從 `minmax(140px, 160px)` 改為 `200px`
- [x] 1.2 在 `.date-panel` 加入 `overflow: hidden` 防止裝飾箭頭溢出

## 2. 增加 Date Card 與 Summary Cards 間距

- [x] 2.1 將 `.hero-side` 的 gap 從 12px 改為 8px
- [x] 2.2 將 `.summary-strip` 各欄最小寬度縮小以釋放空間（130→110, 100→90, 100→90, 148→128, 82→72）

## 3. 驗證

- [ ] 3.1 在 1920px 寬度下檢查結束日期不被遮蓋
- [ ] 3.2 在 1366px 寬度下檢查結束日期不被遮蓋
- [ ] 3.3 在 980px 響應式斷點下檢查 date-panel 不縮小
