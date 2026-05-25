# Vercel 部署指南

本專案可以部署到 Vercel 作為靜態網站，無需後端服務。

## 部署步驟

### 1. 在 Vercel 上建立專案

1. 前往 https://vercel.com
2. 點擊 "Add New Project"
3. 從 GitHub 選擇 `AS-ITS/printerinfo` 專案
4. 點擊 "Deploy"

### 2. 設定環境變數

在 Vercel 專案設定頁面中，新增以下環境變數：

| 變數名稱 | 值 | 說明 |
|---------|---|------|
| `SUPABASE_URL` | `https://xxx.supabase.co` | Supabase 專案 URL |
| `SUPABASE_ANON_KEY` | `eyJhbGci...` | Supabase Anonymous Key |

**取得方式：**
1. 登入 Supabase Dashboard
2. 前往 Settings → API
3. 複製 Project URL 和 anon/public key

### 3. 部署選項

#### 自動部署（推送到 GitHub）

專案已設定自動部署，當推送到 `main` 分支時會自動部署到 Vercel。

#### 手動部署

```bash
# 安裝 Vercel CLI
npm install --global vercel@latest

# 登入 Vercel
vercel login

# 部署到生產環境
vercel --prod
```

## 功能說明

部署後可使用以下頁面：

- **儀表板** - `index.html`
  - 顯示今日總印量、印表機總數、耗材警告數
  - 依單位篩選印表機
  - 每日趨勢圖與統計
  - 耗材低水位警告與事件紀錄

- **年度累計報告** - `index_list.html`
  - 依年份查看各單位、印表機的月度列印統計

- **每日掃描佇列** - `daily-scan-queue.html` (新增)
  - 查看所有印表機的耗材狀態
  - 今日印量統計
  - 故障紀錄
  - 篩選功能
  - 匯出 JSON/CSV

## 安全性注意

1. **不要**將 `SUPABASE_ANON_KEY` 用於伺服器端
2. **不要**將資料庫連線字串（service_role key）暴露在前端
3. Supabase 的 anon key 已授權公開訪問，適合前端使用
4. 在 Supabase Settings → API 中確保「RLS (Row Level Security)」已啟用

## 無障礙設計

所有頁面都支援：
- 響應式設計（手機、平板、桌機）
- 鍵盤導航
- 清晰的對比度和字體大小

## 故障排除

### 資料無法載入

1. 檢查 Supabase 環境變數是否正確設定
2. 確認資料庫 table 存在：`printers`, `supplies`, `printer_metrics`, `incidents`
3. 檢查瀏覽器 Console 是否有錯誤訊息

### 登入失敗

1. 確認 Google OAuth 在 Supabase Settings → Authentication → Providers 中已啟用
2. 設定正確的 `AUTH_ORIGIN` 環境變數（如果需要）

## CI/CD

專案使用 GitHub Actions 自動部署：

- **觸發條件**: 推送到 `main` 分支
- **部署內容**: 靜態網頁
- **不包含**: Go 腳本、資料庫 migration

資料庫 migration 需要手動在 Supabase 執行（如需要）：
```bash
psql "你的連線字串" -f supabase.sql
```

## 授權

此專案使用 Google OAuth 登入，需要：
- Google Cloud Console 中的 OAuth 2.0 Client ID
- Supabase Settings → Authentication → Providers → Google