-- ============================================
-- PrinterInfo Dashboard - Full Migration
-- 可在空資料庫或已有資料的資料庫上執行（安全、可重複執行）
-- 使用方式：psql -d <dbname> -f supabase.sql
-- ============================================

-- 第一步：為 printers 表擴增欄位（已有資料不受影響）
ALTER TABLE printers ADD COLUMN IF NOT EXISTS purchase_date DATE;
ALTER TABLE printers ADD COLUMN IF NOT EXISTS warranty_end DATE;
ALTER TABLE printers ADD COLUMN IF NOT EXISTS status TEXT DEFAULT 'normal';
ALTER TABLE printers ADD COLUMN IF NOT EXISTS unit TEXT;

-- 第二步：建立 printer_metrics 每日數據紀錄表（若不存在）
CREATE TABLE IF NOT EXISTS printer_metrics (
  id SERIAL PRIMARY KEY,
  printer_id INTEGER REFERENCES printers(id),
  total_count INTEGER,
  print_count INTEGER,
  copy_count INTEGER,
  scan_total INTEGER,
  fax_count INTEGER,
  recorded_at DATE DEFAULT CURRENT_DATE,
  UNIQUE(printer_id, recorded_at)
);

-- 第三步：建立 supplies 耗材狀態表（全新表，空表）
CREATE TABLE IF NOT EXISTS supplies (
  id SERIAL PRIMARY KEY,
  printer_id INTEGER REFERENCES printers(id),
  supply_type TEXT NOT NULL CHECK (supply_type IN ('toner', 'ink', 'paper')),
  capacity INTEGER,
  remaining_percent INTEGER DEFAULT 100,
  last_refill_date DATE,
  UNIQUE(printer_id, supply_type)
);

-- 第四步：建立 incidents 故障紀錄表（全新表，空表）
CREATE TABLE IF NOT EXISTS incidents (
  id SERIAL PRIMARY KEY,
  printer_id INTEGER REFERENCES printers(id),
  error_code TEXT,
  description TEXT,
  resolution TEXT,
  cost NUMERIC(10, 2),
  downtime_minutes INTEGER,
  incident_date DATE DEFAULT CURRENT_DATE
);

-- 第五步：為現有印表機建立預設耗材紀錄（安全：用 WHERE NOT EXISTS 避免重複）
INSERT INTO supplies (printer_id, supply_type, remaining_percent)
SELECT p.id, s.type, 100
FROM printers p
CROSS JOIN (VALUES ('toner'), ('ink'), ('paper')) AS s(type)
WHERE NOT EXISTS (
  SELECT 1 FROM supplies WHERE supplies.printer_id = p.id AND supplies.supply_type = s.type
);

-- 第六步：建立/更新 daily_stats 視圖（增加 daily_total 欄位）
CREATE OR REPLACE VIEW daily_stats AS
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
  COALESCE(m.print_count - LAG(m.print_count) OVER w, 0) AS daily_print,
  COALESCE(m.copy_count  - LAG(m.copy_count)  OVER w, 0) AS daily_copy,
  COALESCE(m.scan_total  - LAG(m.scan_total)  OVER w, 0) AS daily_scan,
  COALESCE(m.fax_count   - LAG(m.fax_count)   OVER w, 0) AS daily_fax,
  COALESCE(m.print_count - LAG(m.print_count) OVER w, 0)
    + COALESCE(m.copy_count  - LAG(m.copy_count)  OVER w, 0)
    + COALESCE(m.scan_total  - LAG(m.scan_total)  OVER w, 0)
    + COALESCE(m.fax_count   - LAG(m.fax_count)   OVER w, 0) AS daily_total
FROM printer_metrics m
JOIN printers p ON m.printer_id = p.id
WINDOW w AS (PARTITION BY m.printer_id ORDER BY m.recorded_at);

-- 第七步：建立 dashboard_stats 儀表板整合視圖
-- 查詢所有印表機當前狀態、耗材、保固、故障數
CREATE OR REPLACE VIEW dashboard_stats AS
SELECT
  p.id,
  p.ip_address,
  p.location,
  p.model,
  p.unit,
  COALESCE(p.status, 'normal') AS printer_status,
  COALESCE(supplies_agg.toner_percent, 100) AS toner_percent,
  COALESCE(supplies_agg.ink_percent, 100) AS ink_percent,
  COALESCE(supplies_agg.paper_percent, 100) AS paper_percent,
  p.warranty_end,
  EXTRACT(DAY FROM COALESCE(p.warranty_end, CURRENT_DATE + INTERVAL '365 days') - CURRENT_DATE)::INTEGER AS warranty_days,
  COALESCE(inc_agg.recent_incidents, 0) AS recent_incidents_30d
FROM printers p
LEFT JOIN (
  SELECT printer_id,
    MAX(CASE WHEN supply_type = 'toner' THEN remaining_percent END) AS toner_percent,
    MAX(CASE WHEN supply_type = 'ink'  THEN remaining_percent END) AS ink_percent,
    MAX(CASE WHEN supply_type = 'paper' THEN remaining_percent END) AS paper_percent
  FROM supplies
  GROUP BY printer_id
) supplies_agg ON p.id = supplies_agg.printer_id
LEFT JOIN (
  SELECT printer_id, COUNT(*) AS recent_incidents_30d
  FROM incidents
  WHERE incident_date >= CURRENT_DATE - INTERVAL '30 days'
  GROUP BY printer_id
) inc_agg ON p.id = inc_agg.printer_id;

-- ============================================
-- 執行完成後可跑以下 SQL 確認
-- ============================================
-- 印表機總數：SELECT COUNT(*) FROM printers;
-- 耗材紀錄數：SELECT COUNT(*) FROM supplies;
-- 每日資料數：SELECT COUNT(*) FROM printer_metrics;
-- 印表機狀態：SELECT ip_address, status, toner_percent, ink_percent, warranty_days FROM dashboard_stats;
