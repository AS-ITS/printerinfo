## Context

The dashboard already supports an `All` scope and selected-unit scopes, with charts and the bottom printer table rendered from cached loaded data. The requested change is a presentation-level adjustment: the all-units homepage should emphasize aggregate charts, while printer row details should appear only after a unit is selected.

## Goals / Non-Goals

**Goals:**
- Hide the bottom printer table while `All` is selected.
- Keep the bottom printer table visible and filtered when a specific unit is selected.
- Move `年度總印量統計` before `依單位列印狀況` in the all-units homepage layout.
- Preserve existing cached-data rendering and manual refresh behavior.

**Non-Goals:**
- No database schema changes.
- No new Supabase queries.
- No changes to printer table columns or chart calculations beyond visibility and ordering.

## Decisions

- Use the existing active unit state as the source of truth for table visibility.
  - Rationale: The dashboard already re-renders all scoped sections from the selected unit, so table visibility can follow the same scope branch.
  - Alternative considered: Add a separate table visibility state. This would duplicate the scope state and increase the chance of inconsistent UI.
- Reorder the existing chart sections in the DOM or render flow instead of creating duplicate chart containers.
  - Rationale: Reusing the existing chart instances and containers keeps refresh and empty-state behavior unchanged.
  - Alternative considered: Render a second yearly chart for the homepage. This would add unnecessary chart lifecycle complexity.

## Risks / Trade-offs

- Hidden all-units table could remove a broad overview users previously used. Mitigation: the table remains available through selected-unit views, while aggregate charts remain on the homepage.
- Reordering chart containers could affect responsive spacing. Mitigation: preserve existing section styles and only change order.
