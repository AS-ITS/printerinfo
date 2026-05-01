## Context

The dashboard currently has two unit-selection surfaces: the left-side menu and the all-units `依單位列印狀況` chart. The chart displays unit names in its legend, but those legend entries only control chart dataset visibility through Chart.js defaults. Users now expect clicking a unit legend entry to enter the same unit-scoped view as clicking the left-side menu.

## Goals / Non-Goals

**Goals:**
- Make all-units chart legend entries navigate to the matching unit scope.
- Reuse the existing `filterByUnit(unit)` path so menu state, scoped content, incidents, tables, and selected-unit chart behavior stay consistent.
- Avoid extra Supabase reads when navigating through the legend.
- Keep selected-unit monthly printer chart behavior intact.

**Non-Goals:**
- No new chart library or dependency.
- No database or API changes.
- No change to left-side menu semantics.
- No change to selected-unit printer monthly chart legend behavior unless needed to avoid accidental unit navigation.

## Decisions

1. Use Chart.js legend click customization for the all-units chart.

   The `renderUnitAggregationChart(printers)` function already has access to the sorted `units` array. The legend click handler can resolve the clicked legend item to a unit name and call `filterByUnit(unit)`.

2. Prefer shared navigation over duplicate rendering logic.

   The legend handler should call `filterByUnit(unit)` instead of manually setting state and rendering sections. This keeps behavior equivalent to the left-side menu.

3. Scope the behavior to the all-units chart.

   The selected-unit monthly chart legend represents printer datasets, not units. It should continue to serve chart readability rather than navigating between units.

## Risks / Trade-offs

- [Risk] Chart.js legend item indexing can differ between dataset legends and label legends. -> Mitigation: use the legend item's text when available and validate it against the known `units` list before calling `filterByUnit`.
- [Risk] Users may still expect dataset hide/show behavior. -> Mitigation: this change intentionally repurposes all-units legend clicks for navigation because the legend entries represent unit destinations.
