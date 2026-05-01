const fs = require('fs');

const html = fs.readFileSync('index.html', 'utf8');

function assert(condition, message) {
  if (!condition) {
    throw new Error(message);
  }
}

const yearlyIndex = html.indexOf('id="yearly-chart-section"');
const unitIndex = html.indexOf('id="unit-chart-section"');
const printerIndex = html.indexOf('id="printer-grid"');

assert(yearlyIndex !== -1, 'yearly chart section is missing');
assert(unitIndex !== -1, 'unit chart section is missing');
assert(printerIndex !== -1, 'printer grid section is missing');
assert(yearlyIndex < unitIndex, 'yearly chart must appear before unit chart');
assert(unitIndex < printerIndex, 'printer grid should remain below chart sections');

const renderPrintersStart = html.indexOf('function renderPrinters(printers)');
const isColorPrinterStart = html.indexOf('function isColorPrinter', renderPrintersStart);
assert(renderPrintersStart !== -1 && isColorPrinterStart !== -1, 'renderPrinters function bounds are missing');

const renderPrinters = html.slice(renderPrintersStart, isColorPrinterStart);
assert(renderPrinters.includes("document.getElementById('printer-grid')"), 'renderPrinters must control printer grid visibility');
assert(renderPrinters.includes("currentFilterUnit === 'all'"), 'renderPrinters must branch on all-units scope');
assert(renderPrinters.includes("grid.classList.add('hidden')"), 'all-units scope must hide printer grid');
assert(renderPrinters.includes("tbody.innerHTML = ''"), 'all-units scope must clear printer rows');
assert(renderPrinters.includes("grid.classList.remove('hidden')"), 'selected-unit scope must show printer grid');
assert(renderPrinters.includes("${isColorPrinter(p.model) ? '✔' : '✘'}"), 'color column behavior must remain intact');
assert(renderPrinters.includes('${p.ip_address}</td>'), 'IP column behavior must remain intact');

const renderScopedStart = html.indexOf('function renderScopedDashboard()');
const computeSummaryStart = html.indexOf('function computeSummary', renderScopedStart);
assert(renderScopedStart !== -1 && computeSummaryStart !== -1, 'renderScopedDashboard function bounds are missing');

const renderScoped = html.slice(renderScopedStart, computeSummaryStart);
assert(renderScoped.includes('const scopedPrinters = getScopedPrinters();'), 'scoped printers must be computed once');
assert(renderScoped.includes('renderTrendChart(scopedPrinters);'), 'daily trend chart must use scoped printers');
assert(renderScoped.includes('renderYearlyStatsChart(scopedPrinters);'), 'yearly chart must use scoped printers');
assert(renderScoped.includes('renderScopedUnitChart(scopedPrinters);'), 'unit chart must use scoped printers');

const yearlyStart = html.indexOf('function renderYearlyStatsChart(printers)');
const unitChartComment = html.indexOf('// == 依單位列印狀況圖表 ==', yearlyStart);
assert(yearlyStart !== -1 && unitChartComment !== -1, 'renderYearlyStatsChart function bounds are missing');

const yearlyRenderer = html.slice(yearlyStart, unitChartComment);
assert(!yearlyRenderer.includes("currentFilterUnit !== 'all'"), 'yearly chart must not be hidden for selected units');
assert(yearlyRenderer.includes("document.getElementById('yearly-chart-title')"), 'yearly chart title must be scoped');
assert(yearlyRenderer.includes("currentFilterUnit === 'all'"), 'yearly chart must distinguish all and selected-unit labels');
assert(yearlyRenderer.includes('printers.forEach'), 'yearly chart must aggregate from the provided scoped printers');
assert(!yearlyRenderer.includes('cachedPrinters.forEach'), 'yearly chart must not aggregate directly from all cached printers');

const scopedUnitStart = html.indexOf('function renderScopedUnitChart(printers)');
const unitAggregationStart = html.indexOf('function renderUnitAggregationChart', scopedUnitStart);
assert(scopedUnitStart !== -1 && unitAggregationStart !== -1, 'renderScopedUnitChart function bounds are missing');

const scopedUnitRenderer = html.slice(scopedUnitStart, unitAggregationStart);
assert(scopedUnitRenderer.includes("currentFilterUnit === 'all'"), 'unit chart must branch on all-units scope');
assert(scopedUnitRenderer.includes('renderUnitAggregationChart(printers);'), 'all-units scope must use unit aggregation chart');
assert(scopedUnitRenderer.includes('renderPrinterMonthlyChart(printers);'), 'selected-unit scope must use per-printer chart');

console.log('dashboard layout and scoped chart checks passed');
