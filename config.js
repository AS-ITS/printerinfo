const SUPABASE_URL = 'https://aopzclrptvpfovklhcip.supabase.co';
const SUPABASE_ANON_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFvcHpjbHJwdHZwZm92a2xoY2lwIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzY4MDI1MzgsImV4cCI6MjA5MjM3ODUzOH0.ZXpDc6ADxFXwSj6LYjcZn9viytSHEPA_69E7VY2bVT0';

// == 資料保留上限：預設只查詢最近 MAX_RETENTION_YEARS 年的 daily_stats ==
const MAX_RETENTION_YEARS = 5;
function getRetentionCutoff() {
    const cutoff = new Date();
    cutoff.setFullYear(cutoff.getFullYear() - MAX_RETENTION_YEARS);
    return cutoff.toISOString().substring(0, 10);
}

// == 共用 counter delta 計算 ==
function counterDelta(current, previous, key) {
    return Math.max(0, Number(current?.[key] || 0) - Number(previous?.[key] || 0));
}

function hasCounterReading(row) {
    return ['total_count', 'print_count', 'copy_count', 'scan_total', 'fax_count']
        .some(key => Number(row?.[key] || 0) > 0);
}

// == 共用 Chart.js options 工廠：套用預設樣式 + 個別覆寫 ==
function createChartOptions(overrides = {}) {
    const defaults = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                labels: { color: '#334155', font: { family: 'Plus Jakarta Sans', size: 12, weight: '700' } }
            },
            tooltip: {
                backgroundColor: 'rgba(23, 32, 51, 0.94)',
                titleFont: { family: 'Plus Jakarta Sans' },
                bodyFont: { family: 'Plus Jakarta Sans' }
            }
        },
        scales: {
            x: {
                grid: { color: '#eef1f5' },
                ticks: { color: '#526071', font: { family: 'Plus Jakarta Sans', weight: '700' } }
            },
            y: {
                grid: { color: '#e5e9f0' },
                ticks: { color: '#526071', font: { family: 'Plus Jakarta Sans', weight: '700' } }
            }
        }
    };
    // Deep merge: overrides 會逐層覆蓋 defaults
    return deepMerge(defaults, overrides);
}

function deepMerge(target, source) {
    const result = { ...target };
    for (const key in source) {
        if (Object.prototype.hasOwnProperty.call(source, key)) {
            if (isPlainObject(result[key]) && isPlainObject(source[key])) {
                result[key] = deepMerge(result[key], source[key]);
            } else {
                result[key] = source[key];
            }
        }
    }
    return result;
}

function isPlainObject(val) {
    return val !== null && typeof val === 'object' && !Array.isArray(val);
}

// == 共用 corrected trend 建構 ==
function buildCorrectedTrend(rawMetrics) {
    if (!rawMetrics || rawMetrics.length === 0) return [];
    let previousComparable = null;
    return rawMetrics.map(row => {
        const hasReading = hasCounterReading(row);
        let dailyPrint = 0, dailyCopy = 0, dailyScan = 0, dailyFax = 0, dailyTotal = 0;

        if (hasReading && previousComparable) {
            dailyPrint = counterDelta(row, previousComparable, 'print_count');
            dailyCopy = counterDelta(row, previousComparable, 'copy_count');
            dailyScan = counterDelta(row, previousComparable, 'scan_total');
            dailyFax = counterDelta(row, previousComparable, 'fax_count');
            const totalDelta = counterDelta(row, previousComparable, 'total_count');
            dailyTotal = dailyPrint + dailyCopy + dailyFax;

            // counter reset 回退：所有 daily 為 0 但 total 有變化時使用 totalDelta
            if (dailyPrint === 0 && dailyCopy === 0 && dailyFax === 0 && totalDelta > 0) {
                dailyTotal = totalDelta;
            }
        }

        const corrected = { ...row, daily_print: dailyPrint, daily_copy: dailyCopy, daily_scan: dailyScan, daily_fax: dailyFax, daily_total: dailyTotal };

        if (hasReading) previousComparable = row;
        else if (!previousComparable) previousComparable = row;

        return corrected;
    });
}
