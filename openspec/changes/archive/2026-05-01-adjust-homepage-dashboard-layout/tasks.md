## 1. Layout Rendering

- [x] 1.1 Locate the dashboard sections for the yearly statistics chart, unit aggregation chart, and bottom printer table in `index.html`.
- [x] 1.2 Reorder the all-units homepage layout so `年度總印量統計` renders before `依單位列印狀況`.
- [x] 1.3 Hide the bottom printer table when the active unit scope is `All`.
- [x] 1.4 Show the bottom printer table when a specific unit is selected, using the existing selected-unit filtered printer rows.

## 2. Behavior Integration

- [x] 2.1 Ensure switching between `All` and specific units updates table visibility without a data reload.
- [x] 2.2 Ensure manual refresh preserves the active scope and recalculates visible sections in the new order.
- [x] 2.3 Preserve existing chart empty states and printer table column behavior.

## 3. Verification

- [x] 3.1 Add or update local scripted checks for all-units table hiding, selected-unit table visibility, and chart section order.
- [x] 3.2 Run JavaScript syntax validation for `index.html`.
- [x] 3.3 Run `openspec validate adjust-homepage-dashboard-layout --strict`.
