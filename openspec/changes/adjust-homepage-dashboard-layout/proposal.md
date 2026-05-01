## Why

首頁的全單位總覽應該保持在統計圖表與整體趨勢，避免最下方明細表造成資訊過重。單位明細表則應在使用者進入個別單位後才出現，並調整首頁圖表順序讓年度總印量統計先於依單位列印狀況。

## What Changes

- 在 `All` 首頁範圍隱藏最下方印表機表格。
- 在選取個別單位時顯示該單位的最下方印表機表格。
- 將首頁 `年度總印量統計` 區塊移到 `依單位列印狀況` 之前。
- 保留個別單位畫面既有的篩選、統計與表格欄位行為。

## Capabilities

### New Capabilities

### Modified Capabilities
- `dashboard-unit-navigation`: Change all-units dashboard layout so the printer table is hidden on the homepage, remains visible for selected units, and the yearly chart appears before the unit aggregation chart.

## Impact

- `index.html` dashboard layout and render logic.
- OpenSpec dashboard unit navigation requirements.
- No database schema, Supabase API, or dependency changes.
