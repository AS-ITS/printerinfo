## Context

`.hero-side` 是巢狀 grid，左側放 `.date-panel`，右側放 `.summary-strip`（5 個 summary cards）。當螢幕縮小或 summary cards 最小寬度合計過大時，summary-strip 的左側卡片（列印數量）會覆蓋 date-panel 的右緣，同時 date-panel 內的 `.date-rule` 裝飾箭頭會溢出到 panel 外。

## Goals / Non-Goals

**Goals:**
- date-panel 固定 160px 寬，不再縮小
- 結束日期欄位完全顯示於 date-panel 內
- date-panel 與 summary-strip 間有足夠間距，不產生重疊

**Non-Goals:**
- 改變整體 dashboard 佈局
- 修改 summary-cards 的數值或邏輯

## Decisions

1. **date-panel 固定 200px（而非 minmax）**
   - 原因：minmax(140px, 160px) 在窄螢幕仍會縮小導致內容擠壓
   - 替代方案：使用 flex-basis 或 width: 200px；grid 使用 200px 最簡潔

2. **date-panel 加入 overflow: hidden**
   - 原因：`.date-rule` 裝飾箭頭以 `right: -13px` 定位，會溢出 panel 外
   - 替代方案：修正 `.date-rule` 定位；但 overflow 更全面，可防其他內容溢出

3. **hero-side gap 12px → 8px**
   - 原因：縮小間距讓 date-panel 有更大寬度，同時保留足夠不重疊
   - 替代方案：使用 margin；gap 更簡潔

4. **summary-strip 各欄最小寬度各縮減 20px**
   - 原因：釋放空間給 date-panel，避免整體超出 hero-side 容器
   - 新值：130→110, 100→90, 100→90, 148→128, 82→72

## Risks / Trade-offs

- [summary-cards 文字可能因空間縮小而截斷] → 數字可換行（white-space: normal），label 文字固定不縮短
- [hero-side gap 增大後總寬增加] → 同步縮小 summary-strip 欄最小寬度平衡
