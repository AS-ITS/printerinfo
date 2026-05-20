## Why

印表機首頁的結束日期欄位超出 Date card 範圍，且會被左側的列印數量 card 遮蓋。原因是 date-panel 寬度會隨螢幕縮小、與 summary-strip 間距不足，以及 date-panel 缺少 overflow 保護。

## What Changes

- 將 `.date-panel` 寬度固定為 200px（不再縮小）
- 將 `.date-panel` 加入 `overflow: hidden` 以防止內容溢出
- 將 `.hero-side` 的 gap 從 12px 擴大為 20px，增加 date-panel 與 summary-strip 間距
- 縮小 `.summary-strip` 各欄最小寬度，使其在有限空間內不致溢出

## Capabilities

### New Capabilities
- `date-end-overlap-fix`: 首頁日期選擇器卡片與 summary 卡片的重疊修復

### Modified Capabilities
<!-- none -->

## Impact

- `index.html` — `.date-panel`, `.hero-side`, `.summary-strip` CSS
