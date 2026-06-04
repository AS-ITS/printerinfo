const SUPABASE_URL = 'https://aopzclrptvpfovklhcip.supabase.co';
const SUPABASE_ANON_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImFvcHpjbHJwdHZwZm92a2xoY2lwIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzY4MDI1MzgsImV4cCI6MjA5MjM3ODUzOH0.ZXpDc6ADxFXwSj6LYjcZn9viytSHEPA_69E7VY2bVT0';

// == 共用 counter delta 計算 ==
function counterDelta(current, previous, key) {
    return Math.max(0, Number(current?.[key] || 0) - Number(previous?.[key] || 0));
}

function hasCounterReading(row) {
    return ['total_count', 'print_count', 'copy_count', 'scan_total', 'fax_count']
        .some(key => Number(row?.[key] || 0) > 0);
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
