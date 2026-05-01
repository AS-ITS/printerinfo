## Context

The homepage currently renders a daily trend chart and an all-units unit chart that contains today, month, and year values by unit. Users have asked for an additional yearly statistics chart for all homepage items, giving a clear annual comparison instead of only seeing annual values mixed into the unit chart.

## Goals / Non-Goals

**Goals:**
- Add a homepage yearly statistics chart for the `All` scope.
- Aggregate annual totals from existing loaded printer `trend` rows.
- Update the yearly chart after initial load and manual refresh.
- Keep existing daily, monthly, all-units, selected-unit, table, and modal behavior intact.

**Non-Goals:**
- No new Supabase query.
- No database schema change.
- No Go script change.
- No change to `index_list.html`.
- No replacement of the existing daily trend or unit chart.

## Decisions

1. Add a separate yearly chart section on the homepage.

   A dedicated chart makes annual totals visible without overloading the existing unit chart. The chart should sit near the current chart area so users see daily, unit/month, and yearly context together.

2. Use cached `trend` rows for aggregation.

   Each printer already contains `trend` entries built from `daily_stats`. The yearly chart can sum `daily_total` by year across the currently loaded all-units data without a new query.

3. Render for the all-units scope.

   The request targets the homepage all-items view. Selected-unit chart behavior should remain focused on per-printer monthly current-year output.

## Risks / Trade-offs

- [Risk] Printers with missing trend rows may make yearly totals look lower than expected. -> Mitigation: aggregate only available loaded trend data and render an empty state when no yearly data exists.
- [Risk] Too many historical years can crowd the x-axis. -> Mitigation: sort years ascending and rely on Chart.js tick handling; implementation can cap or format ticks later if needed.
