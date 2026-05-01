## 1. Yearly Chart UI

- [x] 1.1 Add a homepage yearly statistics chart section in `index.html`.
- [x] 1.2 Add a canvas and empty-state element for yearly chart rendering.
- [x] 1.3 Ensure the yearly chart appears in the all-units homepage chart area without replacing existing daily or unit charts.

## 2. Yearly Chart Data And Rendering

- [x] 2.1 Add a yearly chart renderer that aggregates loaded printer `trend` rows by year.
- [x] 2.2 Render annual total print volume from cached loaded data without new Supabase queries.
- [x] 2.3 Re-render the yearly chart after initial dashboard load and manual refresh.
- [x] 2.4 Render an empty state when no yearly trend data exists.

## 3. Verification

- [x] 3.1 Verify all-units homepage shows the yearly statistics chart when yearly data exists.
- [x] 3.2 Verify yearly chart totals are calculated from cached trend rows.
- [x] 3.3 Verify manual refresh recalculates the yearly chart without adding a separate query.
- [x] 3.4 Run JavaScript syntax validation and existing Go tests.
