# PrinterInfo

PrinterInfo 是一個印表機營運監控儀表板，前端以靜態 HTML 呈現，資料來源使用 Supabase PostgreSQL。專案包含即時監控頁、年度累計報告頁、資料庫 migration，以及一支 Go 腳本用來從 Supabase 匯出部署用的 `public/data.json`。

## 功能

- Google OAuth 登入後顯示全院印表機監控儀表板。
- 顯示今日總印量、印表機總數、耗材警告數。
- 依單位篩選印表機，並顯示機型、IP、保固天數、累計印量與今日印量。
- 顯示每日趨勢圖與單位別日、月、年印量統計。
- 顯示耗材低水位警告與近 30 天事件紀錄。
- 提供年度累計報告頁，可依年份查看各單位、各印表機的月度列印、影印、傳真合計。

## 專案結構

```text
.
├── index.html              # 主要監控儀表板
├── index_list.html         # 年度累計報告頁
├── supabase.sql            # Supabase/PostgreSQL schema、view 與 migration
├── scripts/main.go         # 從 Supabase 讀取資料並產生 public/data.json
├── .github/workflows/      # GitHub Actions 自動部署流程
├── images/                 # PWA 與 favicon 圖檔
├── site.webmanifest        # PWA manifest
├── go.mod / go.sum         # Go 腳本依賴
└── makefile                # 簡易 push 指令
```

## 資料模型

`supabase.sql` 會建立或更新下列資料表：

- `printers`：印表機基本資料，包含 IP、位置、型號、單位、購買日、保固日與狀態。
- `printer_metrics`：每日累積計數，包含總量、列印、影印、掃描、傳真與記錄日期。
- `supplies`：耗材狀態，支援 `toner`、`ink`、`paper`。
- `incidents`：故障或維護事件，包含錯誤碼、描述、處理狀態、停機時間與成本。

同時會建立兩個 view：

- `daily_stats`：依印表機與日期用 `LAG` 計算每日增量，產出 `daily_print`、`daily_copy`、`daily_scan`、`daily_fax`、`daily_total`。
- `dashboard_stats`：彙整印表機、耗材、保固天數、近 30 天事件數與未結事件數，供儀表板讀取。

## 前端頁面

### `index.html`

主要營運監控頁。頁面使用：

- Tailwind CSS CDN
- Chart.js CDN
- Supabase JS v2 CDN

此頁會透過 Supabase client 直接讀取：

- `dashboard_stats`
- `supplies`
- `daily_stats`
- `incidents`

頁面內目前已設定 Supabase URL 與 anon key，並使用 Supabase Auth 的 Google OAuth 登入流程。

### `index_list.html`

年度累計報告頁。此頁同樣透過 Supabase JS 直接讀取 `printers` 與 `printer_metrics`，在瀏覽器端計算每日差值、月統計與年度統計。

## 本機開發

這個專案的前端是靜態 HTML，可以直接用瀏覽器開啟：

```powershell
.\index.html
```

如果要用本機靜態伺服器預覽，可在專案根目錄執行：

```powershell
python -m http.server 8000
```

然後開啟：

```text
http://localhost:8000/
```

年度報告頁位於：

```text
http://localhost:8000/index_list.html
```

## 初始化或更新 Supabase schema

先確認已安裝 `psql`，並準備 Supabase PostgreSQL connection string。然後執行：

```powershell
psql "<SUPABASE_DB_CONNECTION>" -f supabase.sql
```

`supabase.sql` 使用 `CREATE TABLE IF NOT EXISTS`、`ALTER TABLE ... ADD COLUMN IF NOT EXISTS` 和重新建立 view 的方式設計，可用於初始化或更新既有資料庫結構。

## 產生 `public/data.json`

Go 腳本會從 Supabase PostgreSQL 讀取 dashboard 相關資料，輸出成 `public/data.json`：

```powershell
$env:SUPABASE_DB_CONNECTION="<SUPABASE_DB_CONNECTION>"
go run .\scripts\main.go
```

輸出的 JSON 結構包含：

- `summary.total_printers`
- `summary.today_total`
- `summary.supply_warning_count`
- `printers[]`
- 每台印表機的耗材、趨勢、事件與累計量

目前前端主要是直接查 Supabase；`public/data.json` 主要用於部署流程或需要靜態資料快照的情境。

## 部署流程

`.github/workflows/deploy.yml` 會在下列情況執行：

- push 到 `main`
- 每日排程執行一次

流程內容：

1. Checkout repository。
2. 安裝 Go。
3. 安裝 PostgreSQL client。
4. 執行 `supabase.sql` migration。
5. 執行 `go run scripts/main.go` 產生 `public/data.json`。
6. 安裝 Vercel CLI。
7. `vercel pull`、`vercel build --prod`。
8. 使用 `vercel deploy --prebuilt --prod` 部署。

GitHub Actions 需要設定下列 secrets：

- `SUPABASE_DB_CONNECTION`
- `SUPABASE_URL`
- `SUPABASE_KEY`
- `VERCEL_ORG_ID`
- `VERCEL_PROJECT_ID`
- `VERCEL_TOKEN`

## Vercel 設定

`.vercelignore` 會排除 Go 腳本、Go module 檔與其他不需要部署到靜態站台的檔案。部署時主要保留 HTML、manifest、favicon、images，以及 workflow 產生的 `public/data.json`。

## 常用指令

```powershell
# 更新 Go 依賴
go mod tidy

# 產生靜態 JSON 資料
$env:SUPABASE_DB_CONNECTION="<SUPABASE_DB_CONNECTION>"
go run .\scripts\main.go

# 執行資料庫 migration
psql "<SUPABASE_DB_CONNECTION>" -f supabase.sql

# 本機預覽
python -m http.server 8000
```

## 注意事項

- `scripts/main.go` 需要 `SUPABASE_DB_CONNECTION`，缺少時會直接結束。
- `daily_stats.daily_total` 是由 SQL view 計算出的每日增量，前端不應改用累積值 `total_count` 當 fallback。
- 耗材警告門檻目前是低於 15%。
- `supabase.sql` 會 drop 並重建 `daily_stats` 與 `dashboard_stats` view，但不會 drop 主要資料表。
- 若變更資料表欄位，需同步檢查 `supabase.sql`、`scripts/main.go`、`index.html` 與 `index_list.html`。
