## ADDED Requirements

### Requirement: 結束日期欄位不得超出 Date card 範圍
結束日期輸入框必須完全顯示於 date-panel 容器內，不得溢出邊界。

#### Scenario: 結束日期顯示完整
- **WHEN** 用戶查看首頁日期選擇器
- **THEN** 結束日期輸入框右緣與 date-panel 右緣對齊且不溢出

### Requirement: Date card 與 Summary cards 不得重疊
date-panel 與 summary-strip 之間必須有足夠間距，summary cards 不得覆蓋 date-panel 內容。

#### Scenario: 日期卡片不被覆蓋
- **WHEN** 用戶查看首頁 header 區域
- **THEN** date-panel 的右緣與 summary-strip 的左緣之間有至少 15px 間距

### Requirement: Date card 固定寬度
date-panel 寬度應固定為 160px，不隨螢幕尺寸縮小。

#### Scenario: 窄螢幕下 date card 不縮小
- **WHEN** 螢幕寬度小於 1200px
- **THEN** date-panel 寬度保持 160px 不變
