## RENAMED Requirements

### Requirement: 頁面標題
- **FROM**: 「今日未列印印表機 | No Print Today」
- **TO**: 「未收到資料列表 | No Data Received」

#### Scenario: 頁面標題顯示正確
- **WHEN** 使用者開啟 no_print_today.html
- **THEN** 瀏覽器標籤顯示「未收到資料列表 | No Data Received」

### Requirement: 登入畫面副標題
- **FROM**: 「今日未列印印表機列表」
- **TO**: 「未收到資料列表」

#### Scenario: 登入畫面顯示正確
- **WHEN** 使用者尚未登入，看到登入畫面
- **THEN** 登入按鈕上方顯示「未收到資料列表」

### Requirement: Hero 區塊標題
- **FROM**: 「今日未列印印表機」
- **TO**: 「未收到資料列表」

#### Scenario: Hero 區塊標題顯示正確
- **WHEN** 使用者登入後看到儀表板
- **THEN** 頁面頂部標題顯示「未收到資料列表」

### Requirement: Hero 區塊副標說明
- **FROM**: 「total_count 與 printer_count 皆為 0 的印表機」
- **TO**: 「未收到有效資料的印表機」

#### Scenario: 副標說明正確
- **WHEN** 使用者看到儀表板頂部
- **THEN** 標題下方說明文字為「未收到有效資料的印表機」

### Requirement: 摘要面板篩選標籤
- **FROM**: 「今日未列印」
- **TO**: 「未收到資料」

#### Scenario: 摘要面板顯示正確
- **WHEN** 使用者看到儀表板摘要
- **THEN** 紅色摘要欄顯示「未收到資料：X 台」

### Requirement: 表格欄位標題 - 時間欄位
- **FROM**: 「最近統計時間」
- **TO**: 「上次抓取時間」

#### Scenario: 表格欄位標題正確
- **WHEN** 使用者查看印表機表格
- **THEN** 第四欄欄位標題顯示「上次抓取時間」

### Requirement: 最後更新標籤
- **FROM**: 「最後更新」
- **TO**: 「上次抓取」

#### Scenario: 最後更新標籤正確
- **WHEN** 使用者看到頂部操作區
- **THEN** 時間標籤顯示「上次抓取：2026/05/15 ...」

### Requirement: 變數名稱
- **FROM**: `last_recorded`
- **TO**: `last_fetched`

#### Scenario: 變數名稱更新
- **WHEN** 檢視 no_print_today.html 原始碼
- **THEN** 資料結構中使用 `last_fetched` 而非 `last_recorded`

### Requirement: 時間顯示變數
- **FROM**: `recordedAt`
- **TO**: `fetchedAt`

#### Scenario: 時間顯示變數更新
- **WHEN** 檢視 no_print_today.html 的 renderTable 函式
- **THEN** 時間格式化變數名為 `fetchedAt`
