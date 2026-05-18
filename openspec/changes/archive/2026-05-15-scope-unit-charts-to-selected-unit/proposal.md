## Why

When a user selects an individual unit, every visible chart should reflect that unit rather than continuing to imply all-unit totals. The selected-unit view needs daily, yearly, and per-printer charting to line up with the active unit filter.

## What Changes

- Keep `每日總印量趨勢` scoped to the selected unit when a unit is active.
- Show `年度總印量統計` for the selected unit instead of hiding it or showing all-unit totals.
- Keep `年度總印量統計` as all-unit totals only when `All` is selected.
- Change `依單位列印狀況` in a selected-unit view to show per-printer print status/statistics inside that unit.
- Continue using cached loaded printer trend data without additional Supabase queries when switching units.

## Capabilities

### New Capabilities

### Modified Capabilities
- `dashboard-unit-navigation`: Refine selected-unit chart scoping so daily trends, yearly totals, and the `依單位列印狀況` chart all reflect the active unit.

## Impact

- `index.html` chart visibility, titles, and render logic.
- OpenSpec dashboard unit navigation requirements.
- No database schema, Supabase API, or dependency changes.
