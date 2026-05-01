## Why

The left-side unit menu is currently the only way to enter a unit-scoped dashboard view, even though the all-units chart already presents the same unit names. Making the chart legend interactive gives users a direct path from the unit overview into the corresponding unit detail state.

## What Changes

- Make the `依單位列印狀況` chart legend entries clickable when the dashboard is in the `All` scope.
- Clicking a unit legend entry changes the active unit exactly as if the user clicked that unit in the left-side menu.
- The active left-side unit menu state and all scoped dashboard content update after legend navigation.
- Keep existing legend behavior for selected-unit printer monthly charts unless a legend entry maps to a unit.
- Do not add any new database queries for this navigation.

## Capabilities

### New Capabilities

### Modified Capabilities
- `dashboard-unit-navigation`: Adds chart legend navigation as an equivalent way to select a unit from the all-units chart.

## Impact

- Affected code: `index.html`
- Affected UI: `依單位列印狀況` chart legend interaction
- Affected behavior: chart legend click calls the same scoped navigation path as the left-side unit menu
- No data model, Supabase query, Go script, authentication, or deployment workflow changes expected
