## 1. Layout And Navigation

- [x] 1.1 Replace the inline `unit-filter-section` with a left-side unit navigation area in `index.html`.
- [x] 1.2 Add responsive layout classes so the unit menu appears on the left on desktop and remains usable above or beside content on mobile.
- [x] 1.3 Render `All` and each distinct unit into the unit navigation with a clear active state.

## 2. Scoped Rendering

- [x] 2.1 Refactor dashboard rendering so initial data load stores complete data in caches and then calls a shared scoped render function.
- [x] 2.2 Update unit selection handling to change `currentFilterUnit` and re-render all scoped right-side dashboard sections.
- [x] 2.3 Ensure `All` renders summary cards, trend chart, unit chart, supply warnings, incidents, and printer table from all loaded printers.
- [x] 2.4 Ensure a specific unit renders summary cards, trend chart, unit chart, supply warnings, incidents, and printer table only from printers in that unit.
- [x] 2.5 Filter incidents by the currently scoped printer IDs when a specific unit is selected.

## 3. Interaction Safety

- [x] 3.1 Avoid unsafe inline JavaScript interpolation for unit names by using event listeners or safe data attributes.
- [x] 3.2 Ensure printer table action buttons still open the correct history and supply monitor for the currently rendered filtered rows.
- [x] 3.3 Preserve existing Supabase auth, data loading, and query behavior without adding new fetches on unit selection.

## 4. Verification

- [x] 4.1 Verify default load selects `All` and displays whole-organization totals.
- [x] 4.2 Verify selecting a unit updates every right-side dashboard section without page reload.
- [x] 4.3 Verify switching back to `All` restores whole-organization statistics.
- [x] 4.4 Verify desktop and mobile layouts do not overlap or hide the unit navigation or dashboard content.
