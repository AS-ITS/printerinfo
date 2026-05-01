## Why

The homepage currently emphasizes daily trends and unit-level daily/monthly summaries, but users also need an annual view to compare long-term output across the whole dashboard. Adding a yearly statistics chart makes the all-items homepage more complete without requiring users to open individual history modals.

## What Changes

- Add a yearly statistics chart to the homepage all-items view.
- The yearly chart summarizes loaded printer history by year across all units when `All` is selected.
- The yearly chart uses existing cached `trend` / `daily_stats` data and does not add new Supabase queries.
- The chart updates when data is initially loaded or manually refreshed.
- Preserve existing daily trend chart, unit chart, selected-unit chart behavior, and table behavior.

## Capabilities

### New Capabilities

### Modified Capabilities
- `dashboard-unit-navigation`: Refines all-units scoped rendering to include a yearly statistics chart on the homepage.

## Impact

- Affected code: `index.html`
- Affected UI: homepage dashboard chart area
- Affected data flow: aggregate yearly totals from already loaded per-printer trend rows
- No expected database schema, Go script, Supabase query, authentication, or deployment workflow changes
