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

console.log('dashboard layout checks passed');
