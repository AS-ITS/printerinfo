-- 1. 建立印表機清單表
CREATE TABLE printers (
  id SERIAL PRIMARY KEY,
  ip_address TEXT UNIQUE NOT NULL,
  location TEXT,
  model TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. 建立每日數據紀錄表
CREATE TABLE printer_metrics (
  id SERIAL PRIMARY KEY,
  printer_id INTEGER REFERENCES printers(id),
  total_count INTEGER,
  print_count INTEGER,
  copy_count INTEGER,
  scan_total INTEGER,
  fax_count INTEGER,
  recorded_at DATE DEFAULT CURRENT_DATE,
  UNIQUE(printer_id, recorded_at) -- 防止同一天重複記錄
);

-- 3. 建立視圖計算每日增量 (統計當天真正印了多少)
--    使用 COALESCE 處理第一天無前日資料的情況（LAG 回傳 NULL → 改為 0）
CREATE OR REPLACE VIEW daily_stats AS
SELECT 
  p.ip_address,
  p.location,
  p.model,
  m.recorded_at,
  -- 絕對數值
  m.total_count,
  m.print_count,
  m.copy_count,
  m.scan_total,
  m.fax_count,
  -- 每日增量（第一天顯示 0，非 NULL）
  COALESCE(m.print_count - LAG(m.print_count) OVER w, 0) AS daily_print,
  COALESCE(m.copy_count  - LAG(m.copy_count)  OVER w, 0) AS daily_copy,
  COALESCE(m.scan_total  - LAG(m.scan_total)  OVER w, 0) AS daily_scan,
  COALESCE(m.fax_count   - LAG(m.fax_count)   OVER w, 0) AS daily_fax
FROM printer_metrics m
JOIN printers p ON m.printer_id = p.id
WINDOW w AS (PARTITION BY m.printer_id ORDER BY m.recorded_at);