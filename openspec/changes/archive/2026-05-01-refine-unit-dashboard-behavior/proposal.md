## Why

The unit-scoped dashboard currently keeps the same unit aggregation chart regardless of whether the user is viewing all units or one unit, which limits usefulness after narrowing the scope. The dashboard also refreshes data too eagerly and the printer list displays the wrong label in the printer-name field, creating unnecessary database reads and confusing table content.

## What Changes

- Show the existing `依單位列印狀況` chart only when `All` is selected.
- When a specific unit is selected, replace that chart area with a per-printer monthly chart for the selected year.
- The per-printer monthly chart shows each printer in the selected unit and its month-by-month print volume for the current year.
- Prevent automatic database reloads when the page merely regains focus or auth state is re-observed without a meaningful session change.
- Add or expose an explicit refresh action that reloads dashboard data on demand.
- Keep the first authenticated dashboard entry loading data from Supabase.
- Fix the bottom printer list so the printer-name display content shows the unit name instead of the printer/location label.
- Preserve existing Supabase schema, table/view names, and authentication provider.

## Capabilities

### New Capabilities
- `dashboard-data-refresh`: Defines when dashboard data is loaded from Supabase and when cached data is reused.

### Modified Capabilities
- `dashboard-unit-navigation`: Refines scoped rendering so the all-units view uses unit aggregation, selected-unit view uses per-printer monthly data, and the printer list label matches the requested unit-name display.

## Impact

- Affected code: `index.html`
- Affected UI: unit chart section, selected-unit chart section, refresh control, printer table display labels
- Affected data flow: dashboard data loading must distinguish initial load, cached render, auth state changes, and manual refresh
- No expected database schema, Go data-generation, Supabase policy, or deployment workflow changes
