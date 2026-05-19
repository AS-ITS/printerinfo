## ADDED Requirements

### Requirement: Summary card numeric values must fit within card boundaries
首頁 SUMMARY 卡片中的數值文字（列印數量、傳真數量、影印數量、總數量）必須在卡片範圍內完整顯示，不得溢出。

#### Scenario: Numeric value fits in summary card
- **WHEN** 用戶查看首頁 summary 卡片
- **THEN** 數值文字大小不超過 28px 且完整顯示於卡片內

#### Scenario: Date picker end date fits in card
- **WHEN** 用戶查看日期選擇器
- **THEN** 結束日期輸入框完整顯示於卡片內，不溢出邊界

### Requirement: Unit printer chart card must not挤 out right-side cards
各單位印表機列印量卡片（包含圖表容器）不得將右側卡片擠出畫面。

#### Scenario: Unit chart fits within grid area
- **WHEN** 首頁載入完成
- **THEN** unit-chart-section 與 yearly-chart-section 共享右側 300px 欄且不溢出
