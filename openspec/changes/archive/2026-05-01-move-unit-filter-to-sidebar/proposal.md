## Why

The current unit filter appears inline above the printer table and only refreshes the table content, which makes the dashboard feel disconnected from the selected unit. Moving units into a persistent left-side menu will make unit selection a primary navigation action and keep the right-side statistics aligned with the selected scope.

## What Changes

- Replace the inline unit filter button group with a left-side unit navigation menu.
- Keep `All` as the default selection.
- When a specific unit is selected, the right-side dashboard content displays only that unit's filtered statistics and printer content.
- When `All` is selected, the right-side dashboard content displays whole-organization statistics across all units.
- Update all dashboard sections that summarize or list scoped data, including summary cards, trend charts, unit charts, supply warnings, incident list, and printer table.
- Preserve existing Supabase data loading and authentication behavior.

## Capabilities

### New Capabilities
- `dashboard-unit-navigation`: Defines the left-side unit navigation and the scoped dashboard content behavior for all-unit and per-unit views.

### Modified Capabilities

## Impact

- Affected code: `index.html`
- Affected UI: dashboard layout, unit filter controls, summary cards, trend chart, unit chart, supply warnings, incidents, printer table
- Affected data flow: existing in-browser cached printer, metric, supply, and incident data must be filtered consistently before rendering
- No database schema, Go script, Supabase API, or deployment workflow changes are expected
