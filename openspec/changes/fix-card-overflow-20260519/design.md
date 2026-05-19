## Context

首頁 dashboard 位於 `index.html`，包含三個溢出問題的區域：

1. **Summary 卡片數值文字**：`.summary-value` 目前 34px，數字過大溢出卡片
2. **日期選擇器**：`.date-field` 的結束日期卡片寬度不足
3. **各單位印表機列印量卡片**：`.dashboard-main-grid` 中 `pcnt` 區域與圖表容器過大

當前 `.dashboard-layout` 為 `minmax(0, 1fr) 300px` 雙欄，`pcnt` 與 `yearly` 共享右側欄。
`.summary-strip` 的 5 欄網格最小寬度合計超過 650px。

## Goals / Non-Goals

**Goals:**
- 數值文字縮小到不溢出卡片的尺寸
- 日期選擇器結束日期不再溢出
- 各單位印表機列印量卡片縮小，不擠出右側卡片

**Non-Goals:**
- 不改變整體佈局結構（仍是雙欄）
- 不變更響應式斷點策略
- 不修改資料邏輯或 API

## Decisions

1. **數值文字：34px → 28px**
   - 縮小 18% 即可容納 6 位數+千分位於現有卡片
   - 同時將 `white-space: nowrap` 移除，改為允許換行

2. **日期面板：寬度從 190px 縮至 160px，`.date-field` 字體 13px → 11px**
   - 結束日期 `<input type="date">` 的瀏覽器預設寬度較大
   - 縮小字體與 padding 即可容納

3. **各單位印表機列印量：`pcnt` 區域使用 `minmax(0, 0.4fr)` 替代預設流式分配**
   - 與 `yearly` 共享右側 300px 欄
   - 圖表容器高度從 252px → 200px 減少垂直空間

## Risks / Trade-offs

- [數值文字縮小後可讀性降低] → 保持在 28px 仍有足夠大小
- [日期面板縮小後 label 可能溢出] → 同步縮小 label 字體與 padding
- [pcnt 區域縮小後圖表標籤可能被截斷] → 配合縮小圖表容器高度
