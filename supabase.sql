-- ============================================
-- PrinterInfo Dashboard - Full Migration (修正版)
-- 可在空資料庫或已有資料的資料庫上執行（安全、可重複執行）
-- 使用方式：psql -d <dbname> -f supabase_fixed.sql
-- ============================================

-- 第一步：確保 printers 基礎表存在（修正：空資料庫也能安全執行）
CREATE TABLE IF NOT EXISTS printers (
  id         SERIAL PRIMARY KEY,
  ip_address TEXT,
  location   TEXT,
  model      TEXT
);

-- 第二步：為 printers 表擴增欄位（已有資料不受影響）
ALTER TABLE printers ADD COLUMN IF NOT EXISTS purchase_date  DATE;
ALTER TABLE printers ADD COLUMN IF NOT EXISTS warranty_end   DATE;
ALTER TABLE printers ADD COLUMN IF NOT EXISTS status         TEXT DEFAULT 'normal';
ALTER TABLE printers ADD COLUMN IF NOT EXISTS unit           TEXT;

-- 第三步：建立 printer_metrics 每日數據紀錄表（若不存在）
CREATE TABLE IF NOT EXISTS printer_metrics (
  id          SERIAL PRIMARY KEY,
  printer_id  INTEGER REFERENCES printers(id),
  total_count INTEGER,
  print_count INTEGER,
  copy_count  INTEGER,
  scan_total  INTEGER,
  fax_count   INTEGER,
  recorded_at DATE DEFAULT CURRENT_DATE,
  UNIQUE(printer_id, recorded_at)
);

-- 第四步：建立 supplies 耗材狀態表
-- 修正：supply_type 的合法值與第六步 INSERT 保持一致，並以 ENUM-like CHECK 集中管理
CREATE TABLE IF NOT EXISTS supplies (
  id                SERIAL PRIMARY KEY,
  printer_id        INTEGER REFERENCES printers(id),
  supply_type       TEXT NOT NULL CHECK (supply_type IN ('toner', 'ink', 'paper')),
  capacity          INTEGER,
  remaining_percent INTEGER DEFAULT 100,
  last_refill_date  DATE,
  UNIQUE(printer_id, supply_type)
);

-- 第五步：建立 incidents 故障紀錄表
-- 修正：新增結構化 incident_status 欄位，支援未解決故障統計與警示
CREATE TABLE IF NOT EXISTS incidents (
  id               SERIAL PRIMARY KEY,
  printer_id       INTEGER REFERENCES printers(id),
  error_code       TEXT,
  description      TEXT,
  resolution       TEXT,
  cost             NUMERIC(10, 2),
  downtime_minutes INTEGER,
  incident_date    DATE    DEFAULT CURRENT_DATE,
  -- 'open' = 未處理, 'in_progress' = 處理中, 'resolved' = 已解決
  incident_status  TEXT    NOT NULL DEFAULT 'open'
                   CHECK (incident_status IN ('open', 'in_progress', 'resolved'))
);

-- 第六步：為已存在的 incidents 表補上 incident_status 欄位
-- 修正：CREATE TABLE IF NOT EXISTS 在表已存在時會跳過，需用 ALTER TABLE 確保欄位存在
ALTER TABLE incidents ADD COLUMN IF NOT EXISTS
  incident_status TEXT NOT NULL DEFAULT 'open'
  CHECK (incident_status IN ('open', 'in_progress', 'resolved'));

-- 第七步：為現有印表機建立預設耗材紀錄（安全：用 WHERE NOT EXISTS 避免重複）
-- 注意：若未來 supply_type CHECK 新增類型，需同步更新此處的 VALUES 清單
INSERT INTO supplies (printer_id, supply_type, remaining_percent)
SELECT p.id, s.type, 100
FROM printers p
CROSS JOIN (VALUES ('toner'), ('ink'), ('paper')) AS s(type)
WHERE NOT EXISTS (
  SELECT 1 FROM supplies
  WHERE supplies.printer_id = p.id
    AND supplies.supply_type = s.type
);

-- 第八步：建立/更新 daily_stats 視圖
-- 修正：LAG 差值改用 GREATEST(..., 0) 防止計數器重置時出現負數
-- 修正：先 DROP 再 CREATE，避免欄位順序或名稱異動時報錯
DROP VIEW IF EXISTS daily_stats CASCADE;
CREATE VIEW daily_stats AS
SELECT
  p.ip_address,
  p.location,
  p.model,
  m.recorded_at,
  m.total_count,
  m.print_count,
  m.copy_count,
  m.scan_total,
  m.fax_count,
  GREATEST(m.print_count - LAG(m.print_count) OVER w, 0) AS daily_print,
  GREATEST(m.copy_count  - LAG(m.copy_count)  OVER w, 0) AS daily_copy,
  GREATEST(m.scan_total  - LAG(m.scan_total)  OVER w, 0) AS daily_scan,
  GREATEST(m.fax_count   - LAG(m.fax_count)   OVER w, 0) AS daily_fax,
  GREATEST(m.print_count - LAG(m.print_count) OVER w, 0)
    + GREATEST(m.copy_count  - LAG(m.copy_count)  OVER w, 0)
    + GREATEST(m.scan_total  - LAG(m.scan_total)  OVER w, 0)
    + GREATEST(m.fax_count   - LAG(m.fax_count)   OVER w, 0) AS daily_total
FROM printer_metrics m
JOIN printers p ON m.printer_id = p.id
WINDOW w AS (PARTITION BY m.printer_id ORDER BY m.recorded_at);

-- 第九步：建立 dashboard_stats 儀表板整合視圖
-- 修正：warranty_end 為 NULL 時 warranty_days 回傳 NULL，避免誤顯示為 365
-- 修正：新增 open_incidents_count 統計未解決故障數
-- 修正：先 DROP 再 CREATE，避免欄位順序或名稱異動時報錯
DROP VIEW IF EXISTS dashboard_stats CASCADE;
CREATE VIEW dashboard_stats AS
SELECT
  p.id,
  p.ip_address,
  p.location,
  p.model,
  p.unit,
  COALESCE(p.status, 'normal') AS printer_status,
  COALESCE(supplies_agg.toner_percent, 100) AS toner_percent,
  COALESCE(supplies_agg.ink_percent,   100) AS ink_percent,
  COALESCE(supplies_agg.paper_percent, 100) AS paper_percent,
  p.warranty_end,
  COALESCE(EXTRACT(DAY FROM (p.warranty_end - CURRENT_DATE))::INTEGER, 0) AS warranty_days,
  COALESCE(inc_agg.recent_incidents_30d, 0) AS recent_incidents_30d,
  -- 新增：未解決故障數（open + in_progress）
  COALESCE(inc_agg.open_incidents_count, 0) AS open_incidents_count
FROM printers p
LEFT JOIN (
  SELECT
    printer_id,
    MAX(CASE WHEN supply_type = 'toner' THEN remaining_percent END) AS toner_percent,
    MAX(CASE WHEN supply_type = 'ink'   THEN remaining_percent END) AS ink_percent,
    MAX(CASE WHEN supply_type = 'paper' THEN remaining_percent END) AS paper_percent
  FROM supplies
  GROUP BY printer_id
) supplies_agg ON p.id = supplies_agg.printer_id
LEFT JOIN (
  SELECT
    printer_id,
    COUNT(*) FILTER (WHERE incident_date >= CURRENT_DATE - INTERVAL '30 days')
      AS recent_incidents_30d,
    COUNT(*) FILTER (WHERE incident_status IN ('open', 'in_progress'))
      AS open_incidents_count
  FROM incidents
  GROUP BY printer_id
) inc_agg ON p.id = inc_agg.printer_id;

-- ============================================
-- 執行完成後可跑以下 SQL 確認
-- ============================================
-- 印表機總數：         SELECT COUNT(*) FROM printers;
-- 耗材紀錄數：         SELECT COUNT(*) FROM supplies;
-- 每日資料數：         SELECT COUNT(*) FROM printer_metrics;
-- 未解決故障數：       SELECT COUNT(*) FROM incidents WHERE incident_status IN ('open', 'in_progress');
-- 印表機完整狀態：     SELECT ip_address, status, toner_percent, ink_percent,
--                             warranty_days, open_incidents_count
--                      FROM dashboard_stats;