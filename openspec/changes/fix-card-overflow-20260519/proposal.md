## Why

首頁 dashboard 的卡片中數值文字過大，導致超出卡片範圍；日期選擇器的結束日期與各單位印表機列印量卡片也過大，會將右側卡片擠出畫面。

## What Changes

- 縮小首頁 SUMMARY 卡片中數值文字大小（列印數量、傳真數量、影印數量、總數量）
- 縮小日期選擇器卡片中結束日期的文字與容器尺寸
- 縮小各單位印表機列印量 CARD（包含圖表容器）的尺寸，避免擠出右側卡片

## Capabilities

### New Capabilities
- `homepage-card-sizing`: 首頁卡片數值文字與容器尺寸的規範

### Modified Capabilities
<!-- none -->

## Impact

- `index.html` 全部樣式與渲染邏輯
  - `.summary-value` 字體大小與卡片網格
  - `.date-panel` 與 `.date-field` 尺寸
  - `.dashboard-main-grid` 的 `pcnt` 區域寬度
  - `.chart-container` 高度
