## 1. Data Loading And Refresh

- [x] 1.1 Refactor dashboard loading so the first authenticated entry fetches Supabase data and later scoped renders reuse cached data.
- [x] 1.2 Prevent repeated auth callbacks or page focus returns for the same user from triggering another Supabase data fetch.
- [x] 1.3 Add a visible refresh action near the dashboard header controls.
- [x] 1.4 Wire the refresh action to force a fresh Supabase fetch and then re-render the currently selected unit scope.
- [x] 1.5 Update the last-update display only after a successful data fetch.

## 2. Scoped Chart Behavior

- [x] 2.1 Keep the existing all-units `依單位列印狀況` aggregation chart visible only when `All` is selected.
- [x] 2.2 Add a selected-unit chart mode that shows current-year monthly print volume by printer.
- [x] 2.3 Compute selected-unit monthly chart data from cached printer `trend` rows without issuing new Supabase queries.
- [x] 2.4 Render a stable empty or zero-data state when the selected unit has no current-year trend data.
- [x] 2.5 Ensure switching between `All` and a unit updates the chart title, datasets, and rendered content without reloading the page.

## 3. Printer List Display

- [x] 3.1 Change the bottom printer list's prominent printer-name field to display the unit name.
- [x] 3.2 Preserve enough secondary printer context, such as location, IP, or model, so rows remain distinguishable.
- [x] 3.3 Verify history and supply actions still open the correct printer after the display-label change.

## 4. Verification

- [x] 4.1 Verify first authenticated dashboard entry fetches Supabase data exactly once for initial render.
- [x] 4.2 Verify repeated focus/auth callbacks for the same user do not refetch data.
- [x] 4.3 Verify manual refresh refetches data and preserves the current selected unit scope.
- [x] 4.4 Verify `All` shows the unit aggregation chart.
- [x] 4.5 Verify selecting a unit shows per-printer monthly current-year print volume.
- [x] 4.6 Verify the bottom printer list shows unit name in the prominent name field.
- [x] 4.7 Run JavaScript syntax validation and existing Go tests.
