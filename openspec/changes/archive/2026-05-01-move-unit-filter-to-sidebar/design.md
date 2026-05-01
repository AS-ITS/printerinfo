## Context

`index.html` currently renders the full dashboard after loading Supabase data, then provides an inline unit filter above the printer table. The existing `filterByUnit(unit)` path only refreshes the printer table and the unit filter buttons, so summary cards, charts, warnings, and incident content remain based on the full data set.

The requested change makes unit selection a primary dashboard navigation control. The left side will list `All` and every available unit. The right side will render the dashboard for the selected scope.

## Goals / Non-Goals

**Goals:**
- Move unit selection from the inline filter section into a persistent left-side menu.
- Make `All` the default selection and display organization-wide statistics.
- Re-render every scoped dashboard section when the selected unit changes.
- Keep the current Supabase queries, auth flow, and source data model unchanged.
- Preserve responsive usability on smaller screens.

**Non-Goals:**
- No database schema change.
- No Go data generation change.
- No new frontend dependency.
- No redesign of `index_list.html`.
- No change to Google OAuth or Supabase project configuration.

## Decisions

1. Use a single scoped render function for the right-side dashboard.

   `renderDashboard()` should continue to load and normalize all data into global caches. After data loading, it should call a new or refactored scoped render function, for example `renderScopedDashboard()`, that derives `filteredPrinters` from `currentFilterUnit` and refreshes all right-side sections.

   Alternative considered: update each click handler to manually call every render function. A central scoped render path is easier to keep correct when future dashboard sections are added.

2. Keep cached source data unfiltered.

   `cachedPrinters`, `cachedMetrics`, `cachedSupplies`, and `cachedIncidents` should continue to represent the complete loaded data set. Filtering should happen at render time so switching from one unit to another does not require another Supabase request.

   Alternative considered: mutate caches per selected unit. That would make `All` restoration error-prone and could break modal detail views that rely on stable indexes.

3. Filter incidents through the selected printers.

   The incident list should show all incidents when `All` is selected. For a specific unit, it should show incidents whose `printer_id` belongs to the filtered printer set.

   Alternative considered: leave incidents global. That would violate the expectation that the right side reflects the selected unit.

4. Replace inline filter UI with a left navigation container.

   The existing `unit-filter-section` should be removed or repurposed into a sidebar menu rendered from the same unit list. The selected item should have a visually distinct active state, and the dashboard content area should sit to the right on desktop.

   Alternative considered: keep both inline buttons and sidebar. Duplicate controls add state synchronization risk and visual clutter.

5. Preserve mobile behavior with a stacked or horizontally scrollable menu.

   On narrow viewports, the left menu can collapse into a top section above the dashboard content or become a horizontal list. It must remain usable without overlapping dashboard content.

## Risks / Trade-offs

- [Risk] Existing modal functions use table row indexes from the currently rendered printer list, while cached data remains global. -> Mitigation: ensure action buttons pass an index or identifier that resolves against the same filtered printer array used to render the table, or refactor actions to use stable printer IDs.
- [Risk] Charts may appear sparse for units with limited history. -> Mitigation: render the existing empty or low-data states gracefully and keep `All` available.
- [Risk] The left menu can crowd smaller screens. -> Mitigation: use responsive layout rules so the menu stacks above content on mobile and remains fixed-width only on larger screens.
- [Risk] Unit names with quotes or special characters can break inline `onclick` handlers. -> Mitigation: prefer event listeners or safe dataset attributes over interpolated inline JavaScript.
