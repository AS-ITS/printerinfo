## Context

The bottom printer table in `index.html` is rendered from static table headers and row template markup inside `renderPrinters(printers)`. The current header order includes `彩色`, `IP`, and `保固`, but the row cells place the IP address where the color indicator belongs and place the color indicator where the IP address belongs. The warranty column is no longer needed.

## Goals / Non-Goals

**Goals:**
- Remove the `保固` table header and corresponding warranty row cell.
- Ensure the `彩色` column displays `✔` or `✘`.
- Ensure the `IP` column displays `p.ip_address`.
- Rename the table header `印表機` to `地點`.
- Preserve existing responsive visibility classes for the color and IP columns unless the implementation must minimally adjust them to keep alignment.

**Non-Goals:**
- No Supabase query changes.
- No change to how color capability is detected.
- No changes to modals, charts, summaries, or refresh behavior.
- No database, Go, or deployment changes.

## Decisions

1. Update header and row markup together.

   The table uses plain HTML headers plus template string cells. The implementation should remove or reorder cells in the same pass to keep header-to-cell alignment exact.

2. Keep `isColorPrinter(model)` as the source of color capability.

   The existing helper already maps model names to a boolean. The change is only to place its output under the correct column.

3. Remove warranty display only from the bottom table.

   The underlying `warranty_days` field can remain in built printer data because other future UI or code may still use it. This change removes only the visible table column.

## Risks / Trade-offs

- [Risk] Removing a table cell without removing the matching header can shift all following columns. -> Mitigation: update header and row template in one small edit and verify generated row order.
- [Risk] Responsive hidden classes can make visual checking harder across breakpoints. -> Mitigation: preserve the existing intended visibility classes while validating desktop markup alignment.
