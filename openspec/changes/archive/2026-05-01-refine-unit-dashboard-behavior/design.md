## Context

`index.html` currently uses one `renderUnitChart(printers)` path for both all-unit and selected-unit scopes. After unit navigation was introduced, selecting a single unit still shows an aggregation chart by unit, which is redundant because only one unit remains. The dashboard also calls `renderDashboard()` from auth initialization and auth state callbacks, so repeated auth notifications or focus-related session refresh behavior can trigger avoidable Supabase reads. The bottom table currently uses `p.location` as the prominent printer-name display, but the requested display should show the unit name in that field.

## Goals / Non-Goals

**Goals:**
- Keep the existing unit aggregation chart visible only for the `All` scope.
- Replace the unit aggregation chart with a selected-unit printer/month chart when one unit is active.
- Add a clear refresh control that explicitly reloads Supabase data.
- Avoid reloading Supabase data merely because the page regains focus or receives an auth callback for the same authenticated user.
- Fix the bottom printer list so the prominent name field displays `unit`.

**Non-Goals:**
- No Supabase schema or RLS change.
- No change to Google OAuth provider setup.
- No redesign of `index_list.html`.
- No Go script change.
- No new frontend dependency.

## Decisions

1. Split chart rendering by selected scope.

   `renderScopedDashboard()` should decide whether to call the all-unit aggregation renderer or a selected-unit monthly printer renderer. When `currentFilterUnit === 'all'`, the existing `依單位列印狀況` section remains visible and uses all loaded printers. When a unit is selected, the same chart section can be retitled and repopulated with a per-printer monthly chart for the current year.

   Alternative considered: show both charts in selected-unit mode. That would make the dashboard longer and keep a low-value one-unit aggregation visible.

2. Compute selected-unit monthly data from cached printer trends.

   Each selected printer already carries `trend` rows built from `daily_stats`. The selected-unit chart should aggregate `daily_total` by printer and month for the current year, using the printer identity for datasets or grouped labels.

   Alternative considered: add another Supabase query for monthly grouping. Cached trends avoid extra reads and match the no-reload requirement.

3. Gate data loading behind explicit load modes.

   Introduce a lightweight loaded/session guard around `renderDashboard()` or split it into `loadDashboardData({ force })` and `renderScopedDashboard()`. The dashboard loads data on first authenticated entry. Later auth callbacks for the same session render from cache unless the user clicks refresh or a force load is explicitly requested.

   Alternative considered: debounce all reloads. Debouncing reduces duplicate bursts but still reloads on focus/auth churn, so it does not fully satisfy the requirement.

4. Add a manual refresh action near the existing last-update/sign-out controls.

   A visible refresh button should call the forced data reload path and then update the scoped dashboard. The existing `last-update` text should reflect the last successful data fetch, not every cached render.

   Alternative considered: rely on browser reload. A local refresh control gives users the requested data update without leaving the application state model ambiguous.

5. Change only the requested table label, not the underlying data model.

   The prominent printer-name field in the bottom list should display `p.unit || '-'`. Supporting details such as location, IP, and model can remain available in secondary text or existing columns so operators still have enough context.

## Risks / Trade-offs

- [Risk] The per-printer monthly chart can become visually dense for units with many printers. -> Mitigation: use responsive chart options and readable labels; if necessary, group by printer rows or stacked monthly bars without adding dependencies.
- [Risk] Auth callbacks after login can still be necessary for the first render. -> Mitigation: key the guard by loaded state and current user ID/email, allowing the first authenticated entry while skipping duplicate callbacks for the same user.
- [Risk] Last-update semantics may confuse cached renders. -> Mitigation: update `last-update` only after a successful Supabase fetch and leave scoped switching as a client-side operation.
- [Risk] Showing unit name in both the unit column and name column may feel repetitive. -> Mitigation: keep location/model/IP in secondary text or adjacent columns to preserve operational context while honoring the requested field content.
