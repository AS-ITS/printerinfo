## Context

The dashboard already computes `scopedPrinters` from the active unit selection and passes that filtered set into summary, daily trend, yearly chart, unit chart, supply warning, printer table, and incident renderers. Daily trend rendering is already scope-aware through that data path, while yearly chart rendering is currently hidden outside the all-units scope. The selected-unit `依單位列印狀況` chart already has a per-printer monthly rendering path, but this change makes that behavior the explicit contract.

## Goals / Non-Goals

**Goals:**
- Ensure the daily trend chart reflects the selected unit when a unit is active.
- Show the yearly statistics chart for selected units using only printers in that unit.
- Keep yearly statistics as all-unit totals when `All` is active.
- Ensure the `依單位列印狀況` chart for a selected unit shows per-printer print statistics within that unit.
- Preserve cached-data rendering and avoid new Supabase queries on unit switching.

**Non-Goals:**
- No changes to Supabase schema or query shape.
- No new chart library or dependency.
- No change to the bottom printer table behavior from the previous layout change.

## Decisions

- Treat `scopedPrinters` as the single data source for all visible dashboard charts.
  - Rationale: The render flow already computes this once per active unit and keeps unit switching local to cached data.
  - Alternative considered: Let each chart filter from `cachedPrinters` independently. That would duplicate filtering rules and make future scope changes more error-prone.
- Remove the all-units-only guard from yearly chart rendering and update copy/title as needed for selected-unit context.
  - Rationale: The same yearly aggregation logic works for all scopes if it receives scoped printer data.
  - Alternative considered: Create a separate selected-unit yearly chart function. That would duplicate chart lifecycle and empty-state logic.
- Keep the selected-unit unit-chart path as a per-printer monthly chart.
  - Rationale: A selected unit no longer needs a unit aggregation chart with only one unit; per-printer data is more useful at that level.

## Risks / Trade-offs

- Selected units with sparse trend rows may show an empty yearly chart. Mitigation: preserve the existing empty-state behavior with scoped data.
- Chart titles could be ambiguous if they do not mention the selected scope. Mitigation: update titles or labels where needed so users can distinguish all-unit totals from selected-unit totals.
