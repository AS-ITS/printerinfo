## 1. Chart Scoping

- [x] 1.1 Review the current `renderScopedDashboard`, `renderTrendChart`, `renderYearlyStatsChart`, and `renderScopedUnitChart` data flow in `index.html`.
- [x] 1.2 Ensure `每日總印量趨勢` renders from selected-unit printers when a unit is active and from all printers when `All` is active.
- [x] 1.3 Update `年度總印量統計` so it renders for both `All` and selected-unit scopes using the active scoped printer set.
- [x] 1.4 Ensure selected-unit yearly totals exclude printers from other units.

## 2. Per-Printer Unit Chart

- [x] 2.1 Ensure `依單位列印狀況` uses the all-unit aggregation chart only when `All` is active.
- [x] 2.2 Ensure `依單位列印狀況` uses per-printer monthly print statistics when a specific unit is active.
- [x] 2.3 Update chart titles or labels so selected-unit chart scope is clear to the user.
- [x] 2.4 Preserve existing empty states when the active scope has no trend data.

## 3. Data Loading Behavior

- [x] 3.1 Confirm unit switching recalculates charts from cached loaded data without issuing new Supabase queries.
- [x] 3.2 Confirm manual refresh recalculates daily, yearly, and per-printer charts for the active scope.

## 4. Verification

- [x] 4.1 Add or update local scripted checks for selected-unit yearly chart visibility and scoped aggregation.
- [x] 4.2 Add or update local scripted checks for selected-unit per-printer `依單位列印狀況` behavior.
- [x] 4.3 Run JavaScript syntax validation for `index.html`.
- [x] 4.4 Run `openspec validate scope-unit-charts-to-selected-unit --strict`.
